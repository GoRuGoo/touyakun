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
	GetMedications(userId string) ([]MedicationListForGetMedications, error)
	DeleteMedications(userId string, dosageId int) error
}

type MedicationListForGetMedications struct {
	Name     string `json:"name"`
	Amount   int    `json:"amount"`
	Duration int    `json:"duration"`
	Time     string `json:"time"`
}

func (dr *DosageRepo) GetMedications(userId string) ([]MedicationListForGetMedications, error) {
	stmt, err := dr.repo.Prepare(
		`
				SELECT
					dosage.name,dosage.amount,dosage.duration,
					time.morning_flg,time.afternoon_flg,time.evening_flg
				FROM
					dosage
				INNER JOIN
					users
				ON
					dosage.user_id = users.id
				INNER JOIN
					time
				ON
					dosage.time_id = time.id
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

		//朝昼夜のフラグを返すのは責務に反しているので一旦変数に落としこんで処理を行ってから構造体に入れる
		var morningFlg, afternoonFlg, eveningFlg bool

		err := rows.Scan(&medication.Name, &medication.Amount, &medication.Duration, &morningFlg, &afternoonFlg, &eveningFlg)
		if err != nil {
			return nil, err
		}

		// 朝昼夜のフラグを見て、どの時間帯に飲むかの設定はモデル層で行うのが適切なのでここで処理しておく
		if morningFlg {
			medication.Time = "morning"
		} else if afternoonFlg {
			medication.Time = "afternoon"
		} else if eveningFlg {
			medication.Time = "evening"
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
				dosage.id = $1
			AND
				dosage.user_id = (SELECT id FROM users WHERE line_user_id = $2);
			`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(dosageId, userId)
	if err != nil {
		return err
	}

	return nil
}
