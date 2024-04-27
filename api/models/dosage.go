package models

import (
	"database/sql"
)

type DosageRepo struct {
	repo *sql.DB
}

func InitializeDosageRepo(db *sql.DB) *DosageRepo {
	return &DosageRepo{repo: db}
}

type DosageModel interface {
	GetMedications(authKey string) ([]MedicationListForGetMedications, error)
}

type MedicationListForGetMedications struct {
	Name       string `json:"name"`
	Amount     int    `json:"amount"`
	Duration   int    `json:"duration"`
	MorningFlg bool   `json:"morning_flg"`
	NoonFlg    bool   `json:"noon_flg"`
	EveningFlg bool   `json:"evening_flg"`
	Time       string `json:"time"`
}

func (dr *DosageRepo) GetMedications(authKey string) ([]MedicationListForGetMedications, error) {
	stmt, err := dr.repo.Prepare(
		`
			SELECT
				dosage.name,dosage.amount,dosage.duration,
				time.morning_flg,time.noon_flg,time.evening_flg,time.night_flg
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
				users.auth_key = ?
			`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(authKey)
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

	return medications, nil
}
