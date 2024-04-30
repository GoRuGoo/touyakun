package models

import (
	"errors"
	"time"
)

type TimeRepo struct {
	repo SqlExecutor
}

func InitializeTimeRepo(db SqlExecutor) *TimeRepo {
	return &TimeRepo{repo: db}
}

type TimeModel interface {
	DeleteTime(lineUserId string, id int) error
	RegisterTime(lineUserId string, time time.Time, isMorning, isAfternoon, isEvening bool) error
}

func (tr *TimeRepo) DeleteTime(lineUserId string, id int) error {
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

type MedicationNotificationTimeInfoForInsertData struct {
	lineUserId  string
	time        time.Time
	isMorning   bool
	isAfternoon bool
	isEvening   bool
}

func (tr *TimeRepo) RegisterTime(lineUserId string, mn MedicationNotificationTimeInfoForInsertData) error {
	formattedRFC3339Time := mn.time.Format(time.RFC3339)

	stmt, err := tr.repo.Prepare(
		`
			INSERT INTO
				time(user_id, time, is_morning, is_afternoon, is_evening)
			VALUES
				((SELECT id FROM users WHERE line_user_id = $1), $2, $3, $4, $5)
			`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(lineUserId, formattedRFC3339Time, mn.isMorning, mn.isAfternoon, mn.isEvening)
	if err != nil {
		return err
	}

	return nil
}