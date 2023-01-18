package router

import (
	"bytes"
	"errors"
	_ "github.com/SakuraBurst/miniature-octo-happiness/cmd/gophermart/docs"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/controller"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/repoitory"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/types"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"io"
	"net/http"
)

type Router struct {
	*echo.Echo
	Endpoint        string
	userController  *controller.GopherMartUserController
	orderController *controller.GopherMartOrderController
}

//swag init --parseDependency --parseInternal -d ./,../../internal/gophermart/router
//swag fmt -d ./,../../internal/gophermart/router

func CreateRouter(endpoint string) *Router {
	userRep, orderRep := repoitory.InitDataBase()
	router := &Router{Echo: echo.New(), Endpoint: endpoint, userController: controller.InitUserController(userRep), orderController: controller.InitOrderController(orderRep)}
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
//	@Summary		Регистрация пользователя
//	@Description	Регистрация производится по паре логин/пароль. Каждый логин должен быть уникальным. После успешной регистрации должна происходить автоматическая аутентификация пользователя. Для передачи аутентификационных данных используйте механизм cookies или HTTP-заголовок Authorization.
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
//	@Summary		Аутентификация пользователя
//	@Description	Аутентификация производится по паре логин/пароль. Для передачи аутентификационных данных используйте механизм cookies или HTTP-заголовок Authorization.
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
		return echo.ErrUnauthorized
	} else if err != nil {
		return echo.ErrInternalServerError
	}
	c.Response().Header().Set("Authorization", "Bearer "+t)
	c.Response().WriteHeader(http.StatusOK)
	return nil
}

// CreateOrder godoc
//
//	@Summary		Загрузка номера заказа
//	@Description	Хендлер доступен только аутентифицированным пользователям. Номером заказа является последовательность цифр произвольной длины.
//	@Tags			orders
//	@Accept			plain
//	@Param			order	body	string	true	"Номер заказа"
//	@Success		200
//	@Success		202
//	@Failure		400	{object}	echo.HTTPError
//	@Failure		401	{object}	echo.HTTPError
//	@Failure		409	{object}	echo.HTTPError
//	@Failure		422	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/user/orders [post]
func (r Router) CreateOrder(c echo.Context) error {
	buf := &bytes.Buffer{}
	_, err := io.Copy(buf, c.Request().Body)
	if err != nil {
		return echo.ErrInternalServerError
	}
	err = c.Request().Body.Close()
	if err != nil {
		return echo.ErrInternalServerError
	}
	if buf.Len() == 0 {
		return echo.ErrBadRequest
	}
	err = r.orderController.CreateOrder(buf.String(), controller.UserLoginFromToken(c.Get("token")), c.Request().Context())
	if errors.Is(err, controller.ErrInvalidOrderId) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity)
	}
	if errors.Is(err, controller.ErrExistingOrderForCurrentUser) {
		c.Response().WriteHeader(http.StatusOK)
		return nil
	}
	if errors.Is(err, controller.ErrExistingOrderForAnotherUser) {
		return echo.NewHTTPError(http.StatusConflict)
	}
	if err != nil {
		return echo.ErrInternalServerError
	}
	c.Response().WriteHeader(http.StatusAccepted)
	return nil
}

// GetOrders godoc
//
//	@Summary		Загрузка номера заказа
//	@Description	Хендлер доступен только авторизованному пользователю. Номера заказа в выдаче должны быть отсортированы по времени загрузки от самых старых к самым новым. Формат даты — RFC3339.
//	@Tags			orders
//	@Success		200	{array}		types.Order
//	@Failure		401	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/user/orders [get]
func (r Router) GetOrders(c echo.Context) error {
	orders, err := r.orderController.GetUserOrders(controller.UserLoginFromToken(c.Get("token")), c.Request().Context())
	if err != nil {
		return echo.ErrInternalServerError
	}
	return c.JSONPretty(http.StatusOK, orders, "  ")
}

// GetBalance godoc
//
//	@Summary		Получение текущего баланса пользователя
//	@Description	Хендлер доступен только авторизованному пользователю. В ответе должны содержаться данные о текущей сумме баллов лояльности, а также сумме использованных за весь период регистрации баллов.
//	@Tags			withdraws
//	@Success		200	{object}	types.UserBalance
//	@Failure		401	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/user/balance [get]
func (r Router) GetBalance(c echo.Context) error {
	balance, err := r.userController.GetUserBalance(controller.UserLoginFromToken(c.Get("token")), c.Request().Context())
	if err != nil {
		return echo.ErrInternalServerError
	}
	return c.JSONPretty(http.StatusOK, balance, "  ")
}

// Withdraw godoc
//
//	@Summary		Запрос на списание средств
//	@Description	Хендлер доступен только авторизованному пользователю. Номер заказа представляет собой гипотетический номер нового заказа пользователя, в счёт оплаты которого списываются баллы.
//	@Tags			withdraws
//	@Accept			json
//	@Param			order	body	string	true	"Номер заказа"
//	@Param			sum		body	number	true	"Сумма баллов к списанию в счёт оплаты"
//	@Success		200
//	@Failure		401	{object}	echo.HTTPError
//	@Failure		402	{object}	echo.HTTPError
//	@Failure		422	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/user/balance/withdraw [post]
func (r Router) Withdraw(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// GetWithdrawals godoc
//
//	@Summary		Получение информации о выводе средств
//	@Description	Хендлер доступен только авторизованному пользователю. Факты выводов в выдаче должны быть отсортированы по времени вывода от самых старых к самым новым. Формат даты — RFC3339.
//	@Tags			withdraws
//	@Success		200	{array}	types.Withdraw
//	@Success		204
//	@Failure		401	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/user/balance [get]
func (r Router) GetWithdrawals(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
