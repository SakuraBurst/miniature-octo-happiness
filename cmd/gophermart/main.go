package main

import (
	"flag"
	"fmt"
	"strconv"
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
	fmt.Println(strconv.FormatFloat(310.2262268066406, 'f', 2, 32))
}

func checkFlags(cfg *config) {
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "Адрес сервера, где будет работать приложение")
	flag.StringVar(&cfg.LoyaltyServiceAddress, "r", cfg.LoyaltyServiceAddress, "Адрес сервиса лояльности")
	flag.StringVar(&cfg.DataBaseURI, "d", cfg.DataBaseURI, "Ссылка для подключения к базе данных")
	flag.Parse()
}
