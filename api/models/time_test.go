package models

import (
	"database/sql"
	"testing"
)

func TestDeleteTime(t *testing.T) {
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
	traitSql, err := ReadSQLFile("./testdata/time_test.sql")
	if err != nil {
		t.Fatalf("Failed to read SQL file: %v", err)
	}
	_, err = tx.Exec(traitSql)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to load test data: %v", err)
	}

	repo := InitializeTimeRepo(tx)

	const TEST_AUTH = "test_auth"
	const TEST_ID = 1

	// DeleteTimeメソッドをテスト
	err = repo.DeleteTime(TEST_AUTH, TEST_ID)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to delete time: %v", err)
	}

	// Check if the record was deleted
	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM time WHERE id = $1", TEST_ID).Scan(&count)

	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to select time: %v", err)
	}

	if count != 0 {
		t.Fatalf("Failed to delete time: time record still exists")
	}
}
