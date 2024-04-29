package models

type UserRepo struct {
	repo SqlExecutor
}

func InitializeUserRepo(db SqlExecutor) *UserRepo {
	return &UserRepo{repo: db}
}

type UserModel interface {
	IsNotExistUser(userId string) (bool, error)
	RegisterUser(userId string) error
}

func (ur *UserRepo) IsNotExistUser(userId string) (bool, error) {
	stmt, err := ur.repo.Prepare(
		`
			SELECT
				COUNT(*)
			FROM
				users
			WHERE
				users.line_user_id = $1;
			`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(userId).Scan(&count)
	if err != nil {
		return false, err
	}

	if count >= 1 {
		return false, err
	}

	return true, nil
}

func (ur *UserRepo) RegisterUser(userId string) error {
	stmt, err := ur.repo.Prepare(
		`
			INSERT INTO
				users(line_user_id)
			VALUES
				($1)
			`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId)
	if err != nil {
		return err
	}

	return nil
}
