package router

import (
	"errors"
	"fmt"
	_ "github.com/SakuraBurst/miniature-octo-happiness/cmd/gophermart/docs"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/controller"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/types"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

type Router struct {
	*echo.Echo
	Endpoint       string
	userController *controller.GopherMartUserController
}

//swag init --parseDependency --parseInternal -d ./,../../internal/gophermart/router

func CreateRouter(endpoint string) *Router {
	router := &Router{Echo: echo.New(), Endpoint: endpoint, userController: controller.InitUserController()}
	router.GET("/swagger/*", echoSwagger.WrapHandler)
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	userApi := router.Group("/api/user")
	userApi.POST("/register", router.Register)
	userApi.POST("/login", router.Login)
	authGroup := userApi.Group("/")
	authGroup.Use(echojwt.WithConfig(controller.UserTokenConfig()))
	authGroup.POST("orders", router.CreateOrder)
	authGroup.GET("orders", router.GetOrders)
	authGroup.GET("balance", router.GetBalance)
	authGroup.POST("withdraw", router.Withdraw)
	authGroup.GET("withdrawals", router.GetWithdrawals)
	return router
}

// Register godoc
//
//	@Summary		Регистрация
//	@Description	Регистрация
//	@Tags			user
//	@Accept			json
//	@Param			login		body	string	true	"Логин нового пользователя"
//	@Param			password	body	string	true	"Пароль нового пользователя"
//	@Success		200
//	@Failure		400	{object}	echo.HTTPError
//	@Failure		409	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/user/register [post]
func (r Router) Register(c echo.Context) error {
	userRequest := new(types.UserRequest)
	c.Bind(userRequest)
	if !userRequest.IsValid() {
		return echo.ErrBadRequest
	}
	t, err := r.userController.Register(userRequest.Login, userRequest.Password, c.Request().Context())
	if errors.Is(err, controller.ErrExistingUser) {
		return echo.NewHTTPError(http.StatusConflict)
	} else if err != nil {
		return echo.ErrInternalServerError
	}
	c.Response().Header().Set("Authorization", "Bearer "+t)
	c.Response().WriteHeader(http.StatusOK)
	return nil
}

// Login godoc
//
//	@Summary		Логин
//	@Description	Логин
//	@Tags			user
//	@Accept			json
//	@Param			login		body	string	true	"Логин пользователя"
//	@Param			password	body	string	true	"Пароль пользователя"
//	@Success		200
//	@Failure		400	{object}	echo.HTTPError
//	@Failure		401	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/user/login [post]
func (r Router) Login(c echo.Context) error {
	userRequest := new(types.UserRequest)
	c.Bind(userRequest)
	if !userRequest.IsValid() {
		return echo.ErrBadRequest
	}
	t, err := r.userController.Login(userRequest.Login, userRequest.Password, c.Request().Context())
	if errors.Is(err, controller.ErrNoExist) {
		return echo.NewHTTPError(http.StatusUnauthorized)
	} else if err != nil {
		fmt.Println(err)
		return echo.ErrInternalServerError
	}
	c.Response().Header().Set("Authorization", "Bearer "+t)
	c.Response().WriteHeader(http.StatusOK)
	return nil
}
func (r Router) CreateOrder(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome "+controller.UserLoginFromToken(c.Get("token"))+"!")
}
func (r Router) GetOrders(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
func (r Router) GetBalance(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
func (r Router) Withdraw(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
func (r Router) GetWithdrawals(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
