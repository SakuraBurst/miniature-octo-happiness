package controller

import (
	"context"
	"errors"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/repoitory"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/types"
	"github.com/jackc/pgx/v5"
)

type GopherMartWithdrawController struct {
	repository repoitory.WithdrawTable
}

func InitWithdrawController(table repoitory.WithdrawTable) *GopherMartWithdrawController {
	return &GopherMartWithdrawController{repository: table}
}

func (c *GopherMartWithdrawController) CreateWithdraw(login, orderId string, sum float64, userController *GopherMartUserController, context context.Context) error {
	if !Luhn(orderId) {
		return ErrInvalidOrderId
	}
	err := userController.WithdrawUserBalance(login, sum, context)
	if err != nil {
		return err
	}
	return c.repository.CreateWithdraw(login, orderId, sum, context)
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
