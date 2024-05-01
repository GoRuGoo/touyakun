package models

import (
	"database/sql"
	"testing"
	"time"
)

func TestRegisterMorningTime(t *testing.T) {
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

	testLineUserId := "test_id"
	testTime := time.Now()

	err = repo.RegisterMorningTime(testLineUserId, testTime)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to register morning time: %v", err)
	}

	tx.Rollback()
}

func TestRegisterAfternoonTime(t *testing.T) {
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

	testLineUserId := "test_id"
	testTime := time.Now()

	err = repo.RegisterAfternonnTime(testLineUserId, testTime)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to register morning time: %v", err)
	}

	tx.Rollback()
}

func TestRegisterEveningTime(t *testing.T) {
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

	testLineUserId := "test_id"
	testTime := time.Now()

	err = repo.RegisterEveningTime(testLineUserId, testTime)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to register morning time: %v", err)
	}

	tx.Rollback()
}
