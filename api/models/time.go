package models

import (
	"time"
)

type TimeRepo struct {
	repo SqlExecutor
}

func InitializeTimeRepo(db SqlExecutor) *TimeRepo {
	return &TimeRepo{repo: db}
}

type TimeModel interface {
	RegisterMorningTime(lineUserId string, insertTime time.Time) error
	RegisterAfternonnTime(lineUserId string, insertTime time.Time) error
	RegisterEveningTime(lineUserId string, insertTime time.Time) error
}

type MedicationNotificationTimeInfoForInsertData struct {
	time        time.Time
	isMorning   bool
	isAfternoon bool
	isEvening   bool
}

func (tr *TimeRepo) RegisterMorningTime(lineUserId string, insertTime time.Time) error {
	formattedRFC3339Time := insertTime.Format(time.TimeOnly)

	stmt, err := tr.repo.Prepare(
		`
			UPDATE
				users
			SET
				morning_medication_time = $1
			WHERE
				line_user_id = $2
			`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(formattedRFC3339Time, lineUserId)
	if err != nil {
		return err
	}

	return nil
}

func (tr *TimeRepo) RegisterAfternonnTime(lineUserId string, insertTime time.Time) error {
	formattedRFC3339Time := insertTime.Format(time.TimeOnly)

	stmt, err := tr.repo.Prepare(
		`
			UPDATE
			    users
			SET
				afternoon_medication_time = $1
			WHERE
				line_user_id = $2
			`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(formattedRFC3339Time, lineUserId)
	if err != nil {
		return err
	}

	return nil
}

func (tr *TimeRepo) RegisterEveningTime(lineUserId string, insertTime time.Time) error {
	formattedRFC3339Time := insertTime.Format(time.TimeOnly)

	stmt, err := tr.repo.Prepare(
		`
			UPDATE
			    users
			SET
				evening_medication_time = $1
			WHERE
				line_user_id = $2
			`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(formattedRFC3339Time, lineUserId)
	if err != nil {
		return err
	}

	return nil
}
