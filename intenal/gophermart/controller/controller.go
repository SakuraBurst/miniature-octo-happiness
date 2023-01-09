package controller

import (
	"context"
	"fmt"
)

type GopherMartController struct {
}

func (c GopherMartController) IsUserLoggedIn(email, password string, context context.Context) (bool, error) {
	fmt.Printf("%s, %s, %s, %v", "auth check triggered", email, password, context)
	return true, nil
}
