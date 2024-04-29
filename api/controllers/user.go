package controllers

import (
	"net/http"
	"touyakun/models"
)

type UserController struct {
	userModel models.UserModel
	w         http.ResponseWriter
}

func InitializeUserController(u models.UserModel, w http.ResponseWriter) *UserController {
	return &UserController{
		userModel: u,
		w:         w,
	}
}

func (uc *UserController) RegisterUser(userId string) {
	isNotExist, err := uc.userModel.IsNotExistUser(userId)
	if err != nil {
		uc.w.WriteHeader(500)
		uc.w.Write([]byte(err.Error()))
		return
	}

	if !isNotExist {
		uc.w.WriteHeader(409)
		uc.w.Write([]byte("user already exists"))
		return
	}

	err = uc.userModel.RegisterUser(userId)
	if err != nil {
		uc.w.WriteHeader(500)
		uc.w.Write([]byte(err.Error()))
		return
	}
	uc.w.WriteHeader(200)
	return
}

func (uc *UserController) DeleteUser(userId string) {
	isNotExist, err := uc.userModel.IsNotExistUser(userId)
	if err != nil {
		uc.w.WriteHeader(500)
		uc.w.Write([]byte(err.Error()))
		return
	}

	if isNotExist {
		uc.w.WriteHeader(404)
		uc.w.Write([]byte("user does not exist"))
		return
	}

	err = uc.userModel.DeleteUser(userId)
	if err != nil {
		uc.w.WriteHeader(500)
		uc.w.Write([]byte(err.Error()))
		return
	}

	uc.w.WriteHeader(200)
	return
}
