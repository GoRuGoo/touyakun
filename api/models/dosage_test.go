package models

import (
	"bytes"
	"database/sql"
	"io/ioutil"
	"testing"

	_ "github.com/lib/pq"
)

func TestGetMedications(t *testing.T) {
	// データベースに接続
	db, err := sql.Open("postgres", "host=localhost port=5433 user=testcase password=password dbname=testcase sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// テスト用のSQLとはいえデータを破壊したくないのでトランザクションを開始
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	// テストデータをロード
	traitSql, err := ReadSQLFile("./testdata/dosage_test.sql")
	if err != nil {
		t.Fatalf("Failed to read SQL file: %v", err)
	}
	_, err = tx.Exec(traitSql)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to load test data: %v", err)
	}

	// テスト対象のリポジトリを初期化
	repo := InitializeDosageRepo(tx)

	// GetMedicationsメソッドをテスト
	medications, err := repo.GetMedications("test_auth")
	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to get medications: %v", err)
	}

	// 期待する結果
	expected := []MedicationListForGetMedications{
		{
			Name:     "トラネキサム",
			Amount:   1,
			Duration: 3,
			Time:     "morning",
		},
	}

	// 結果が期待通りであることを確認
	if len(medications) != len(expected) {
		tx.Rollback()
		t.Fatalf("Expected %d result(s), but got %d", len(expected), len(medications))
	}

	for i, medication := range medications {
		if medication.Name != expected[i].Name || medication.Amount != expected[i].Amount || medication.Duration != expected[i].Duration || medication.Time != expected[i].Time {
			tx.Rollback()
			t.Errorf("Unexpected result: got %+v, want %+v", medication, expected[i])
		}
	}
	tx.Rollback()
}

func ReadSQLFile(path string) (string, error) {
	// ファイルを読み込む
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	b := bytes.NewBuffer(data)

	return b.String(), nil
}
