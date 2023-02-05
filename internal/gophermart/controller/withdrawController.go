package controller

import (
	"context"
	"errors"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/repository"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/types"
	"github.com/jackc/pgx/v5"
)

type GopherMartWithdrawController struct {
	repository repository.WithdrawTable
}

func InitWithdrawController(table repository.WithdrawTable) *GopherMartWithdrawController {
	return &GopherMartWithdrawController{repository: table}
}

func (c *GopherMartWithdrawController) CreateWithdraw(login, orderID string, sum float64, userController *GopherMartUserController, context context.Context) error {
	if !Luhn(orderID) {
		return ErrInvalidOrderID
	}
	err := userController.WithdrawUserBalance(login, sum, context)
	if err != nil {
		return err
	}
	return c.repository.CreateWithdraw(login, orderID, sum, context)
}

func (c *GopherMartWithdrawController) GetUserWithdrawals(login string, context context.Context) ([]types.Withdraw, error) {
	w, err := c.repository.GetAllWithdrawalsByLogin(login, context)
	if err == nil {
		return w, err
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return nil, err
}
