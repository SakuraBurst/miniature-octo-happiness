package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/repoitory"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type GopherMartUserController struct {
	repository repoitory.UserTable
}

var ErrExistingUser = errors.New("user already exist")
var ErrNoExist = errors.New("no exist")

func (uc *GopherMartUserController) IsUserLoggedIn(login, password string, context context.Context) (bool, error) {
	fmt.Printf("%s, %s, %s, %v", "auth check triggered", login, password, context)
	return true, nil
}

func (uc *GopherMartUserController) Register(login, password string, c context.Context) error {
	_, err := uc.repository.GetUser(login, c)
	fmt.Println("3434")
	fmt.Println(err)
	if err == nil {
		return ErrExistingUser
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return err
	}
	return uc.repository.CreateUser(login, string(hashedPassword), c)
}

func (uc *GopherMartUserController) Login(login, password string, c context.Context) error {
	user, err := uc.repository.GetUser(login, c)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNoExist
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
