package controllers

type APIError struct {
	StatusCode             int
	ErrorMessages          string
	ProductionErrorMessage string
}
