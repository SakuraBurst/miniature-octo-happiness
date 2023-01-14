package main

import (
	"encoding/base64"
	"fmt"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/router"
)

//	@title			Gophermart API
//	@version		1.0
//	@description	No description.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api

// @securityDefinitions.basic	BasicAuth
func main() {
	credentials := []byte("govno@gmail.com:pass")
	fmt.Println(base64.StdEncoding.EncodeToString(credentials))
	r := router.CreateRouter(":8080")
	r.Logger.Fatal(r.Start(r.Endpoint))
}
