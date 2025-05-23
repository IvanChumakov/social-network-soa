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
	grpcClient  pb.PostsServiceClient
	statsClient pb.StatisticsServiceClient
}

func NewApp(client pb.PostsServiceClient, statsClient pb.StatisticsServiceClient) *App {
	return &App{
		grpcClient:  client,
		statsClient: statsClient,
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

func (a *App) LikeEvent(w http.ResponseWriter, r *http.Request) {
	logger.Info("POST /like")
	if r.Method != http.MethodPost {
		logger.Error("POST /like: method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := JWTTokenVerify(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}

	userID := r.Header.Get("user_id")
	postId := r.URL.Query().Get("id")
	postIdCasted, _ := strconv.Atoi(postId)

	like := pb.Like{
		PostId: int32(postIdCasted),
	}

	md := metadata.Pairs("user_id", userID)
	ctx := metadata.NewOutgoingContext(r.Context(), md)
	_, err = a.grpcClient.LikeEvent(ctx, &like)
	if err != nil {
		logger.Error(fmt.Sprintf("Like event failed: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *App) ViewEvent(w http.ResponseWriter, r *http.Request) {
	logger.Info("POST /view")
	if r.Method != http.MethodPost {
		logger.Error("POST /view: method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := JWTTokenVerify(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}

	userID := r.Header.Get("user_id")
	postId := r.URL.Query().Get("id")
	postIdCasted, _ := strconv.Atoi(postId)

	view := pb.View{
		PostId: int32(postIdCasted),
	}

	md := metadata.Pairs("user_id", userID)
	ctx := metadata.NewOutgoingContext(r.Context(), md)
	_, err = a.grpcClient.ViewEvent(ctx, &view)
	if err != nil {
		logger.Error(fmt.Sprintf("View event failed: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *App) AddComment(w http.ResponseWriter, r *http.Request) {
	logger.Info("POST /comment")
	if r.Method != http.MethodPost {
		logger.Error("POST /comment: method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := JWTTokenVerify(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}

	var comment pb.Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		logger.Error(fmt.Sprintf("comment json decode failed: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Header.Get("user_id")
	md := metadata.Pairs("user_id", userID)
	ctx := metadata.NewOutgoingContext(r.Context(), md)

	_, err = a.grpcClient.AddComment(ctx, &comment)
	if err != nil {
		logger.Error(fmt.Sprintf("add comment failed: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *App) GetAllComments(w http.ResponseWriter, r *http.Request) {
	logger.Info("GET /comments")
	if r.Method != http.MethodGet {
		logger.Error("GET /comments: method not allowed")
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

	postId := r.URL.Query().Get("id")
	md := metadata.Pairs("post_id", postId)
	ctx := metadata.NewOutgoingContext(r.Context(), md)

	comments, err := a.grpcClient.GetAllCommentsPaginated(ctx, &pb.Pagination{
		PageSize:  int32(pageSize),
		PageIndex: int32(index),
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Get all comments failed: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(comments)
}

func (a *App) GetLikesViewsComments(w http.ResponseWriter, r *http.Request) {
	logger.Info("GET /stats")
	if r.Method != http.MethodGet {
		logger.Error("GET /stats: method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := JWTTokenVerify(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}

	postId, _ := strconv.Atoi(r.URL.Query().Get("post_id"))
	postIdProto := pb.PostId{
		PostId: int32(postId),
	}

	stats, err := a.statsClient.GetLikesViewsComments(r.Context(), &postIdProto)
	if err != nil {
		logger.Error(fmt.Sprintf("Get likes views comments failed: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(stats)
}

func (a *App) GetViewsDynamic(w http.ResponseWriter, r *http.Request) {
	logger.Info("GET /views-dynamic")
	if r.Method != http.MethodGet {
		logger.Error("GET /views-dynamic: method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := JWTTokenVerify(r); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}

	postId, _ := strconv.Atoi(r.URL.Query().Get("post_id"))
	postIdProto := pb.PostId{
		PostId: int32(postId),
	}

	dynamic, err := a.statsClient.GetViewsDynamic(r.Context(), &postIdProto)
	if err != nil {
		logger.Error(fmt.Sprintf("Get views dynamic failed: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(dynamic)
}

func (a *App) GetLikesDynamic(w http.ResponseWriter, r *http.Request) {
	logger.Info("GET /likes-dynamic")
	if r.Method != http.MethodGet {
		logger.Error("GET /likes-dynamic: method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := JWTTokenVerify(r); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}

	postId, _ := strconv.Atoi(r.URL.Query().Get("post_id"))
	postIdProto := pb.PostId{
		PostId: int32(postId),
	}

	dynamic, err := a.statsClient.GetLikesDynamic(r.Context(), &postIdProto)
	if err != nil {
		logger.Error(fmt.Sprintf("Get views dynamic failed: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(dynamic)
}

func (a *App) GetCommentsDynamic(w http.ResponseWriter, r *http.Request) {
	logger.Info("GET /comments-dynamic")
	if r.Method != http.MethodGet {
		logger.Error("GET /comments-dynamic: method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := JWTTokenVerify(r); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}

	postId, _ := strconv.Atoi(r.URL.Query().Get("post_id"))
	postIdProto := pb.PostId{
		PostId: int32(postId),
	}

	dynamic, err := a.statsClient.GetCommentsDynamic(r.Context(), &postIdProto)
	if err != nil {
		logger.Error(fmt.Sprintf("Get comments dynamic failed: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(dynamic)
}

func (a *App) GetPostsTop(w http.ResponseWriter, r *http.Request) {
	logger.Info("GET /posts-top")
	if r.Method != http.MethodGet {
		logger.Error("GET /top-post: method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := JWTTokenVerify(r); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}

	typeEvent := r.URL.Query().Get("type")
	typePb := pb.Type{
		Type: typeEvent,
	}

	top, err := a.statsClient.GetPostsTop(r.Context(), &typePb)
	if err != nil {
		logger.Error(fmt.Sprintf("Get top posts failed: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(top)
}

func (a *App) GetUsersTop(w http.ResponseWriter, r *http.Request) {
	logger.Info("GET /users-top")
	if r.Method != http.MethodGet {
		logger.Error("GET /top-user: method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := JWTTokenVerify(r); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}

	typeEvent := r.URL.Query().Get("type")
	typePb := pb.Type{
		Type: typeEvent,
	}

	top, err := a.statsClient.GetUsersTop(r.Context(), &typePb)
	if err != nil {
		logger.Error(fmt.Sprintf("Get top users failed: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(top)
}
