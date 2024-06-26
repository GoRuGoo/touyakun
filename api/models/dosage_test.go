package models

import (
	"bytes"
	"database/sql"
	"io/ioutil"
	"reflect"
	"testing"

	_ "github.com/lib/pq"
)

func TestRegisterMedications(t *testing.T) {
	// Connect to the database
	db, err := sql.Open("postgres", "host=localhost port=5433 user=testcase password=password dbname=testcase sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// Begin a transaction
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

	// Initialize the DosageRepo
	repo := InitializeDosageRepo(tx)

	testMedicationList := []MedicationListForRegisterMedications{
		{Name: "トラネキサム", Amount: 1, Duration: 2, IsMorning: true, IsAfternoon: false, IsEvening: true},
		{Name: "トラネキサム2", Amount: 1, Duration: 2, IsMorning: false, IsAfternoon: false, IsEvening: true},
		{Name: "トラネキサム3", Amount: 1, Duration: 2, IsMorning: true, IsAfternoon: false, IsEvening: false},
	}
	testId := "test_id_for_register_medications"

	// Call the RegisterMedications method
	err = repo.RegisterMedications(testMedicationList, testId)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to register medication: %v", err)
	}

	// Check if the medication was registered successfully
	medications, err := repo.GetMedications(testId)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to get medications: %v", err)
	}

	expected := []MedicationListForGetMedications{
		{
			Id:          1,
			Name:        "トラネキサム",
			Amount:      1,
			Duration:    2,
			IsMorning:   true,
			IsAfternoon: false,
			IsEvening:   true,
		},
		{
			Id:          1,
			Name:        "トラネキサム2",
			Amount:      1,
			Duration:    2,
			IsMorning:   false,
			IsAfternoon: false,
			IsEvening:   true,
		},
		{
			Id:          1,
			Name:        "トラネキサム3",
			Amount:      1,
			Duration:    2,
			IsMorning:   true,
			IsAfternoon: false,
			IsEvening:   false,
		},
	}

	for i, medication := range medications {
		//Idは連番でつくので考慮しないようにあらかじめ期待する結果を入れておく
		medication.Id = expected[i].Id

		if !reflect.DeepEqual(medication, expected[i]) {
			tx.Rollback()
			t.Errorf("Unexpected result: got %+v, want %+v", medication, expected[i])
		}
	}

	tx.Rollback()
}
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
	medications, err := repo.GetMedications("test_id")
	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to get medications: %v", err)
	}

	// 期待する結果
	expected := []MedicationListForGetMedications{
		{
			Id:          1,
			Name:        "トラネキサム",
			Amount:      2,
			Duration:    7,
			IsMorning:   true,
			IsAfternoon: false,
			IsEvening:   true,
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

func TestDeleteMedications(t *testing.T) {
	// Connect to the database
	db, err := sql.Open("postgres", "host=localhost port=5433 user=testcase password=password dbname=testcase sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	// Load the test data
	traitSql, err := ReadSQLFile("./testdata/dosage_test.sql")
	if err != nil {
		t.Fatalf("Failed to read SQL file: %v", err)
	}
	_, err = tx.Exec(traitSql)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to load test data: %v", err)
	}

	// Initialize the DosageRepo
	repo := InitializeDosageRepo(tx)

	// Call the DeleteMedications method
	err = repo.DeleteMedications("test_id", 1)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to delete medication: %v", err)
	}

	// Check if the medication was deleted successfully
	medications, err := repo.GetMedications("test_id")
	// If the medication was deleted successfully, the length of medications should be 0
	if err == nil {
		tx.Rollback()
		t.Fatalf("Expected no medications, but got %d", len(medications))
	}

	tx.Rollback()
}

func TestGettingAnErrorWhenTryingToDeleteARecordThatDoesNotExist(t *testing.T) {
	// Connect to the database
	db, err := sql.Open("postgres", "host=localhost port=5433 user=testcase password=password dbname=testcase sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	// Load the test data
	traitSql, err := ReadSQLFile("./testdata/dosage_test.sql")
	if err != nil {
		t.Fatalf("Failed to read SQL file: %v", err)
	}
	_, err = tx.Exec(traitSql)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to load test data: %v", err)
	}

	// Initialize the DosageRepo
	repo := InitializeDosageRepo(tx)

	// Call the DeleteMedications method
	err = repo.DeleteMedications("test_id", 3343214)
	if err == nil {
		tx.Rollback()
		t.Fatalf("Expected an error, but got none")
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
