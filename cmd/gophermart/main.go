package main

import (
	"encoding/base64"
	"fmt"
	"github.com/SakuraBurst/miniature-octo-happiness/intenal/gophermart/router"
)

func main() {
	credentials := []byte("govno@gmail.com:pass")
	fmt.Println(base64.StdEncoding.EncodeToString(credentials))
	r := router.CreateRouter(":8080")
	r.Logger.Fatal(r.Start(r.Endpoint))
}
