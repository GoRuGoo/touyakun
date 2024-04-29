package controllers

import (
	"errors"
	"touyakun/models"
)

type UserController struct {
	userModel models.UserModel
}

func InitializeUserController(u models.UserModel) *UserController {
	return &UserController{
		userModel: u,
	}
}

func (uc *UserController) RegisterUser(userId string) error {
	isNotExist, err := uc.userModel.IsNotExistUser(userId)
	if err != nil {
		return err
	}

	if !isNotExist {
		return errors.New("user already exists")
	}

	err = uc.userModel.RegisterUser(userId)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UserController) DeleteUser(userId string) error {
	isNotExist, err := uc.userModel.IsNotExistUser(userId)
	if err != nil {
		return err
	}

	if isNotExist {
		return errors.New("user does not exist")
	}

	err = uc.userModel.DeleteUser(userId)
	if err != nil {
		return err
	}

	return nil
}
