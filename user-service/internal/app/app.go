package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	customError "social-network/user-service/internal/errors"
	"social-network/user-service/internal/logger"
	"social-network/user-service/internal/repository"
	"social-network/user-service/internal/service"
)

type App struct {
	userService service.UserServiceInterface
}

func NewApp(userService service.UserServiceInterface) *App {
	return &App{
		userService: userService,
	}
}

func (app *App) Register(w http.ResponseWriter, r *http.Request) {
	logger.Info("POST /register")

	user := repository.User{}
	data, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(data, &user)
	if err != nil {
		logger.Error(fmt.Sprintf("bad request"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := app.userService.Register(&user)
	if err != nil {
		var alreadyExists *customError.LoginAlreadyTakenError
		w.WriteHeader(http.StatusBadRequest)

		if errors.As(err, &alreadyExists) {
			_, _ = fmt.Fprint(w, err.Error())
		}
		return
	}
	_ = json.NewEncoder(w).Encode(token)
	w.WriteHeader(http.StatusOK)
}

func (app *App) Login(w http.ResponseWriter, r *http.Request) {
	logger.Info("POST /login")

	user := repository.User{}
	data, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(data, &user)
	if err != nil {
		logger.Error(fmt.Sprintf("bad request"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := app.userService.Login(&user)
	if err != nil {
		var notFoundError *customError.NotFoundUserError
		w.WriteHeader(http.StatusBadRequest)

		if errors.As(err, &notFoundError) {
			_, _ = fmt.Fprint(w, err.Error())
		}
		return
	}
	_ = json.NewEncoder(w).Encode(token)
	w.WriteHeader(http.StatusOK)
}

func (app *App) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	logger.Info("GET /user-profile")

	user, err := app.userService.GetUserProfile(r.Header.Get("login"))
	if err != nil {
		var notFoundError *customError.NotFoundUserError
		w.WriteHeader(http.StatusBadRequest)

		if errors.As(err, &notFoundError) {
			_, _ = fmt.Fprint(w, err.Error())
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(*user)
	w.WriteHeader(http.StatusOK)
}

func (app *App) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	logger.Info("PUT /user-profile")

	user := repository.User{}
	data, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(data, &user)
	if err != nil {
		logger.Error(fmt.Sprintf("bad request"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = app.userService.UpdateUserProfile(r.Header.Get("login"), &user)
	if err != nil {
		var updateCredsErr *customError.UpdateCredentialsError
		w.WriteHeader(http.StatusBadRequest)

		if errors.As(err, &updateCredsErr) {
			_, _ = fmt.Fprint(w, err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}
