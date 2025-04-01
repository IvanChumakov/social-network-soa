package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"social-network/api-gateway/internal/logger"
	_ "social-network/api-gateway/internal/models"
	pb "social-network/protos"
	"strconv"
	"strings"

	"google.golang.org/grpc/metadata"
)

type App struct {
	grpcClient pb.PostsServiceClient
}

func NewApp(client pb.PostsServiceClient) *App {
	return &App{
		grpcClient: client,
	}
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

func (a *App) CreatePost(w http.ResponseWriter, r *http.Request) {
	logger.Info("POST /post")
	if r.Method != http.MethodPost {
		logger.Error("POST /post: method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := JWTTokenVerify(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}

	var post pb.PostEssential
	data, _ := io.ReadAll(r.Body)
	err = json.Unmarshal(data, &post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}

	userID := r.Header.Get("user_id")
	md := metadata.Pairs("user_id", userID)

	ctx := metadata.NewOutgoingContext(r.Context(), md)
	_, err = a.grpcClient.AddPost(ctx, &post)
	if err != nil {
		logger.Error(fmt.Sprintf("Add post failed: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (a *App) DeletePost(w http.ResponseWriter, r *http.Request) {
	logger.Info("DELETE /post")
	if r.Method != http.MethodDelete {
		logger.Error("DELETE /post: method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := JWTTokenVerify(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	message := pb.PostId{
		PostId: int32(id),
	}
	_, err = a.grpcClient.DeletePost(context.Background(), &message)
	if err != nil {
		logger.Error(fmt.Sprintf("Delete post failed: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}
}

func (a *App) UpdatePost(w http.ResponseWriter, r *http.Request) {
	logger.Info("PUT /post")
	if r.Method != http.MethodPut {
		logger.Error("PUT /post: method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := JWTTokenVerify(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}

	var post pb.PostWithNoUser
	data, _ := io.ReadAll(r.Body)
	err = json.Unmarshal(data, &post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = a.grpcClient.UpdatePost(context.Background(), &post)
	if err != nil {
		logger.Error(fmt.Sprintf("Update post failed: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, err.Error())
	}
}

func (a *App) GetPostById(w http.ResponseWriter, r *http.Request) {
	logger.Info("GET /get-post-by-id")
	if r.Method != http.MethodGet {
		logger.Error("GET /get-post-by-id: method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := JWTTokenVerify(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}

	split := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(split[len(split)-1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	message := pb.PostId{
		PostId: int32(id),
	}
	post, err := a.grpcClient.GetPostById(context.Background(), &message)
	if err != nil {
		logger.Error(fmt.Sprintf("Get post failed: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}

	_ = json.NewEncoder(w).Encode(post)
}

func (a *App) GetPosts(w http.ResponseWriter, r *http.Request) {
	logger.Info("GET /get-posts")
	if r.Method != http.MethodGet {
		logger.Error("GET /get-posts: method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := JWTTokenVerify(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	index, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Header.Get("user_id")
	md := metadata.Pairs("user_id", userID)
	ctx := metadata.NewOutgoingContext(r.Context(), md)

	pagination := pb.Pagination{
		PageSize:  int32(pageSize),
		PageIndex: int32(index),
	}

	posts, err := a.grpcClient.GetAllPostsPaginated(ctx, &pagination)
	if err != nil {
		logger.Error(fmt.Sprintf("Get all posts failed: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(posts)
}
