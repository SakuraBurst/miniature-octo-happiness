package main

import (
	"context"
	"database/sql/driver"
	"errors"
	"flag"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/controller"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/repository"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/router"
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type appState int

const (
	appStateRunning  appState = iota
	appStateShutdown appState = iota
)

var currentState appState = -1

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
	log.SetReportCaller(true)
	if err := env.Parse(cfg); err != nil {
		log.Fatal(err)
	}
	checkFlags(cfg)
	userRep, orderRep, withdrawRep, db, err := repository.InitDataBase(cfg.DataBaseURI)
	if err != nil {
		log.Fatal(err)
	}
	userController := controller.InitUserController(userRep, cfg.SecretTokenKey)
	orderController, err := controller.InitOrderController(orderRep, cfg.LoyaltyServiceAddress)
	if err != nil {
		log.Info(err)
		time.Sleep(time.Second * 3)
		log.Fatal(err)
	}
	withdrawController := controller.InitWithdrawController(withdrawRep)
	r := router.CreateRouter(":8080", userController, orderController, withdrawController)
	go func() {
		currentState = appStateRunning
		if err := r.Start(r.Endpoint); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	watcherChan := make(chan struct{})
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		defer close(sigChan)
		defer close(watcherChan)
		dbWatch(db)
	}()
	<-sigChan
	log.Info("server is shutting down")
	currentState = appStateShutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := r.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Info("router is shut down")
	<-watcherChan
	shutdown(orderController, db)
}

func shutdown(orderRep *controller.GopherMartOrderController, db repository.DB) {
	orderRep.CloseQueue()
	for !orderRep.IsQueueEmpty() {
		time.Sleep(time.Millisecond * 50)
	}
	log.Info("queue is cleared")
	ctx, cl := context.WithTimeout(context.Background(), time.Second)
	defer cl()
	if err := db.Close(ctx); err != nil {
		log.Error(err)
	}
	log.Info("db is closed")
}

func dbWatch(pinger driver.Pinger) {
	for {
		if currentState != appStateRunning {
			return
		}
		ctx, cl := context.WithTimeout(context.Background(), time.Millisecond*500)
		err := pinger.Ping(ctx)
		if err != nil {
			cl()
			return
		}
		cl()
		time.Sleep(time.Millisecond * 500)
	}
}

func checkFlags(cfg *config) {
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "Адрес сервера, где будет работать приложение")
	flag.StringVar(&cfg.LoyaltyServiceAddress, "r", cfg.LoyaltyServiceAddress, "Адрес сервиса лояльности")
	flag.StringVar(&cfg.DataBaseURI, "d", cfg.DataBaseURI, "Ссылка для подключения к базе данных")
	flag.Parse()
}
