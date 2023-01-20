package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/repoitory"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/types"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type GopherMartUserController struct {
	repository repoitory.UserTable
}
type jwtTokenClaims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}

var ErrExistingUser = errors.New("user already exist")
var ErrNoExist = errors.New("no exist")
var secret = []byte("secret")

func InitUserController(table repoitory.UserTable) *GopherMartUserController {
	return &GopherMartUserController{repository: table}
}

func UserTokenConfig() echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtTokenClaims)
		},
		ContextKey: "token",
		SigningKey: secret,
	}
}

func UserLoginFromToken(token any) string {
	user := token.(*jwt.Token)
	claims := user.Claims.(*jwtTokenClaims)
	return claims.Login
}

func (uc *GopherMartUserController) Register(login, password string, c context.Context) (string, error) {
	_, err := uc.repository.GetUser(login, c)
	if err == nil {
		return "", ErrExistingUser
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		fmt.Println(err)
		return "", err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return "", err
	}
	err = uc.repository.CreateUser(login, string(hashedPassword), c)
	if err != nil {
		return "", err
	}
	return uc.createUserToken(login)
}

func (uc *GopherMartUserController) Login(login, password string, c context.Context) (string, error) {
	user, err := uc.repository.GetUser(login, c)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNoExist
	}
	if err != nil {
		return "", err
	}
	fmt.Println(user, login, password)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return "", ErrNoExist
	}
	if err != nil {
		return "", err
	}
	return uc.createUserToken(login)
}

func (uc *GopherMartUserController) GetUserBalance(login string, c context.Context) (*types.UserBalance, error) {
	user, err := uc.repository.GetUser(login, c)
	if err != nil {
		return nil, err
	}
	balance := new(types.UserBalance)
	balance.CurrentBalance = user.Balance
	balance.Withdraw = user.Withdraw
	return balance, nil
}

func (uc *GopherMartUserController) UpdateUserBalance(login string, balance float64, c context.Context) error {
	user, err := uc.repository.GetUser(login, c)
	if err != nil {
		return err
	}
	return uc.repository.UpdateBalanceAndWithdraw(login, user.Balance+balance, user.Withdraw, c)
}

func (uc *GopherMartUserController) createUserToken(login string) (string, error) {
	claims := &jwtTokenClaims{
		login,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
