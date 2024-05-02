package models

import (
	"errors"
)

type DosageRepo struct {
	repo SqlExecutor
}

func InitializeDosageRepo(db SqlExecutor) *DosageRepo {
	return &DosageRepo{repo: db}
}

type DosageModel interface {
	RegisterMedications(ml []MedicationListForRegisterMedications, userId string) error
	GetMedications(userId string) ([]MedicationListForGetMedications, error)
	DeleteMedications(userId string, dosageId int) error
}

type MedicationListForRegisterMedications struct {
	Name        string `json:"name"`
	IsMorning   bool   `json:"isMorning"`
	IsAfternoon bool   `json:"isAfternoon"`
	IsEvening   bool   `json:"isEvening"`
	Duration    int    `json:"duration"`
}

type MedicationListForGetMedications struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Amount      int    `json:"amount"`
	Duration    int    `json:"duration"`
	Time        string `json:"time"`
	IsMorning   bool   `json:"isMorning"`
	IsAfternoon bool   `json:"isAfternoon"`
	IsEvening   bool   `json:"isEvening"`
}

func (dr *DosageRepo) RegisterMedications(ml []MedicationListForRegisterMedications, userId string) error {
	stmt, err := dr.repo.Prepare(
		`
			INSERT INTO
				dosage(name, amount, duration, morning_flg, afternoon_flg, evening_flg, user_id)
			VALUES
				($1, $2, $3, $4, $5, $6, (SELECT id FROM users WHERE line_user_id = $7));
			`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, m := range ml {
		_, err := stmt.Exec(m.Name, m.Duration, m.IsMorning, m.IsAfternoon, m.IsEvening, userId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dr *DosageRepo) GetMedications(userId string) ([]MedicationListForGetMedications, error) {
	stmt, err := dr.repo.Prepare(
		`
				SELECT
					dosage.id ,dosage.name,dosage.amount,dosage.duration,
					dosage.morning_flg,dosage.afternoon_flg,dosage.evening_flg
				FROM
					dosage
				INNER JOIN
					users
				ON
					dosage.user_id = users.id
				WHERE
					users.line_user_id = $1;
				`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//クエリ実行時にエラーがあるのならばこの時点で処理しておく
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var medications []MedicationListForGetMedications

	for rows.Next() {
		//レスポンス用の構造体をappendして構造体配列にして返す
		var medication MedicationListForGetMedications

		err := rows.Scan(&medication.Id, &medication.Name, &medication.Amount, &medication.Duration, &medication.IsMorning, &medication.IsAfternoon, &medication.IsEvening)
		if err != nil {
			return nil, err
		}

		medications = append(medications, medication)
	}

	//配列が取得で木なかった場合のエラーをここで返す
	if len(medications) == 0 {
		return nil, errors.New("Record not found.")
	}

	return medications, nil
}

func (dr *DosageRepo) DeleteMedications(userId string, dosageId int) error {
	stmt, err := dr.repo.Prepare(
		`
			DELETE
			FROM
				dosage
			WHERE
				id = $1
			AND
				user_id = (SELECT id FROM users WHERE line_user_id = $2);
			`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(dosageId, userId)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("Record not found.")
	}

	return nil
}
