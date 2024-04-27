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
	GetMedications() int
}

func (dr *DosageRepo) GetMedications() int {
	return 0
}
