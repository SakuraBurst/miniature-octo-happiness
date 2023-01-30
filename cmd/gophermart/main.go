package main

import (
	"flag"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/controller"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/repoitory"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/router"
	"github.com/caarlos0/env/v6"
	"log"
)

type config struct {
	ServerAddress         string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
	LoyaltyServiceAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"http://localhost:8081/"`
	SecretTokenKey        string `env:"SECRET_KEY" envDefault:"secret"`
	DataBaseURI           string `env:"DATABASE_URI" envDefault:"postgres://postgres:pescola@localhost:5432/gophermart"`
}

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
	cfg := new(config)
	if err := env.Parse(cfg); err != nil {
		log.Fatal(err)
	}
	checkFlags(cfg)
	userRep, orderRep, withdrawRep := repoitory.InitDataBase(cfg.DataBaseURI)
	userController := controller.InitUserController(userRep, cfg.SecretTokenKey)
	orderController := controller.InitOrderController(orderRep, cfg.LoyaltyServiceAddress)
	withdrawController := controller.InitWithdrawController(withdrawRep)
	r := router.CreateRouter(":8080", userController, orderController, withdrawController)
	r.Logger.Fatal(r.Start(r.Endpoint))
}

func checkFlags(cfg *config) {
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "Адрес сервера, где будет работать приложение")
	flag.StringVar(&cfg.LoyaltyServiceAddress, "r", cfg.LoyaltyServiceAddress, "Адрес сервиса лояльности")
	flag.StringVar(&cfg.DataBaseURI, "d", cfg.DataBaseURI, "Ссылка для подключения к базе данных")
	flag.Parse()
}
