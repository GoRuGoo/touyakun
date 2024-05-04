package models

import (
	"database/sql"
	. "github.com/go-playground/assert"
	"testing"
	"time"
)

func TestGetNotificationList(t *testing.T) {
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
	traitSql, err := ReadSQLFile("./testdata/notification_test.sql")
	if err != nil {
		t.Fatalf("Failed to read SQL file: %v", err)
	}
	_, err = tx.Exec(traitSql)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to load test data: %v", err)
	}

	// Initialize the NotificationRepo
	repo := InitializeNotificationRepo(tx)

	testTime, _ := time.Parse("15:04", "08:00")
	// Call the GetNotificationList method
	notifications, err := repo.GetNotificationList(testTime)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to get notifications: %v", err)
	}

	// Check the results...
	// (Your test assertions here)
	expected := []NotificationList{
		{
			LineUserId:  "test_user_1",
			DosageName:  "Test Drug 1",
			DosageAmout: "2",
		},
	}

	Equal(t, expected, notifications)

	tx.Rollback()
}
