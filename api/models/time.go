package models

import "errors"

type TimeRepo struct {
	repo SqlExecutor
}

func InitializeTimeRepo(db SqlExecutor) *TimeRepo {
	return &TimeRepo{repo: db}
}

type TimeModel interface {
	DeleteTime(authKey string, id int) error
}

func (tr *TimeRepo) DeleteTime(authKey string, id int) error {
	stmt, err := tr.repo.Prepare(
		`
			DELETE
			FROM
				time
			USING
				users
			WHERE
				time.id = $1
			AND
				time.user_id = users.id
			AND 
				users.auth_key = $2
			`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id, authKey)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows were deleted")
	}

	return nil
}
