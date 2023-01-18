package controller

import (
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/repoitory"
)

type GopherMartOrderController struct {
	repository repoitory.OrderTable
}

func InitOrderController(table repoitory.OrderTable) *GopherMartOrderController {
	return &GopherMartOrderController{repository: table}
}
