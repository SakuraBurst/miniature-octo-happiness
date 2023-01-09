package router

import (
	"github.com/SakuraBurst/miniature-octo-happiness/intenal/gophermart/controller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type Router struct {
	*echo.Echo
	Endpoint   string
	controller controller.GopherMartController
}

func CreateRouter(endpoint string) *Router {
	router := &Router{Echo: echo.New(), Endpoint: endpoint, controller: controller.GopherMartController{}}
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	userApi := router.Group("/api/user")
	userApi.GET("/register", router.Register)
	userApi.GET("/login", router.Login)
	authGroup := userApi.Group("/")
	authGroup.Use(authMiddleware(router.controller))
	authGroup.POST("orders", router.CreateOrder)
	authGroup.GET("orders", router.GetOrders)
	authGroup.GET("balance", router.GetBalance)
	authGroup.POST("withdraw", router.Withdraw)
	authGroup.GET("withdrawals", router.GetWithdrawals)
	return router
}

func authMiddleware(c controller.GopherMartController) echo.MiddlewareFunc {
	return middleware.BasicAuth(func(email string, password string, context echo.Context) (bool, error) {
		return c.IsUserLoggedIn(email, password, context.Request().Context())
	})
}

func (r Router) Register(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
func (r Router) Login(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
func (r Router) CreateOrder(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
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
