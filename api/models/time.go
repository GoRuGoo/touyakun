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
	GetMedicationRemindTimeList(lineUserId string) (MedicationRemindTimeList, error)
	RegisterMorningTime(lineUserId string, insertTime time.Time) error
	RegisterAfternonnTime(lineUserId string, insertTime time.Time) error
	RegisterEveningTime(lineUserId string, insertTime time.Time) error
}

type MedicationRemindTimeList struct {
	MorningTime   string `json:"morningTime"`
	AfternoonTime string `json:"afternoonTime"`
	EveningTime   string `json:"eveningTime"`
}

func (tr *TimeRepo) GetMedicationRemindTimeList(lineUserId string) (MedicationRemindTimeList, error) {
	stmt, err := tr.repo.Prepare(
		`
			SELECT
			    morning_medication_time, afternoon_medication_time, evening_medication_time
			FROM
				users
			WHERE
				line_user_id = $1
			`)
	if err != nil {
		return MedicationRemindTimeList{}, err
	}
	defer stmt.Close()

	var morningTime, afternoonTime, eveningTime time.Time
	err = stmt.QueryRow(lineUserId).Scan(&morningTime, &afternoonTime, &eveningTime)
	if err != nil {
		return MedicationRemindTimeList{}, err
	}

	varidatedMorningTime := morningTime.Format(time.TimeOnly)
	varidatedAfternoonTime := afternoonTime.Format(time.TimeOnly)
	varidatedEveningTime := eveningTime.Format(time.TimeOnly)

	return MedicationRemindTimeList{
		MorningTime:   varidatedMorningTime,
		AfternoonTime: varidatedAfternoonTime,
		EveningTime:   varidatedEveningTime,
	}, nil
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
