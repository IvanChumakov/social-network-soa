package app

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"social-network/api-gateway/internal/logger"
	_ "social-network/api-gateway/internal/models"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) createProxy(host string) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   host,
	})

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = host
		req.Host = host
	}
	return proxy
}

// Register godoc
// @Summary      Регистрация
// @Description  Зарегистрироваться в сервисе
// @Tags         Auth
// @Accept		 json
// @Produce      json
// @Param 		 user body models.RegisterModel true "Создать пользователя"
// @Success      200  {string} string
// @Router       /register [post]
func (a *App) Register(w http.ResponseWriter, r *http.Request) {
	logger.Info("POST /register")

	if r.Method != http.MethodPost {
		logger.Error("POST /register: method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	proxy := a.createProxy("user-service:8081")
	proxy.ServeHTTP(w, r)
}

// Login godoc
// @Summary      Войти
// @Description  Войти в систему
// @Tags         Auth
// @Accept		 json
// @Produce      json
// @Param 		 user body models.LoginModel true "Войти в систему"
// @Success      200  {string} string
// @Router       /login [post]
func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	logger.Info("POST /login")

	if r.Method != http.MethodPost {
		logger.Error("POST /login: method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	proxy := a.createProxy("user-service:8081")
	proxy.ServeHTTP(w, r)
}

// GetUserProfile godoc
// @Summary      Получить пользователя
// @Description  Получить пользователя
// @Tags         User
// @Accept		 application/x-www-form-urlencoded
// @Security BearerAuth
// @Produce      json
// @Success      200  {object} models.UserModel
// @Router       /user-profile [get]
func (a *App) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	logger.Info("GET /user-profile")

	if r.Method != http.MethodGet {
		logger.Error("GET /user-profile: method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := JWTTokenVerify(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}

	proxy := a.createProxy("user-service:8081")
	proxy.ServeHTTP(w, r)
}

// UpdateUserProfile godoc
// @Summary      Обновить пользователя
// @Description  Обновить данные о пользователе
// @Tags         User
// @Accept		 json
// @Security BearerAuth
// @Produce      json
// @Param 		 user body models.UserModel true "Обновить пользователя"
// @Success      200
// @Router       /user-profile [put]
func (a *App) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	logger.Info("PUT /user-profile")

	if r.Method != http.MethodPut {
		logger.Error("PUT /user-profile: method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := JWTTokenVerify(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}

	proxy := a.createProxy("user-service:8081")
	proxy.ServeHTTP(w, r)
}
