package server

import (
	"context"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/fx"
	"net/http"
	"social-network/api-gateway/internal/app"
	"social-network/api-gateway/internal/config"
	"social-network/api-gateway/internal/logger"
)

func NewServer(cfg *config.Config, app *app.App) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/register", http.HandlerFunc(app.Register))
	mux.Handle("/login", http.HandlerFunc(app.Login))
	mux.HandleFunc("/user-profile", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			app.GetUserProfile(w, r)
		case http.MethodPut:
			app.UpdateUserProfile(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			app.CreatePost(w, r)
		case http.MethodDelete:
			app.DeletePost(w, r)
		case http.MethodPut:
			app.UpdatePost(w, r)
		case http.MethodGet:
			app.GetPosts(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.Handle("/post/", http.HandlerFunc(app.GetPostById))
	mux.Handle("/like", http.HandlerFunc(app.LikeEvent))
	mux.Handle("/view", http.HandlerFunc(app.ViewEvent))
	mux.HandleFunc("/comments", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			app.GetAllComments(w, r)
		case http.MethodPost:
			app.AddComment(w, r)
		}
	})

	mux.Handle("/stats", http.HandlerFunc(app.GetLikesViewsComments))
	mux.Handle("/views-dynamic", http.HandlerFunc(app.GetViewsDynamic))
	mux.Handle("/likes-dynamic", http.HandlerFunc(app.GetLikesDynamic))
	mux.Handle("/comments-dynamic", http.HandlerFunc(app.GetCommentsDynamic))
	mux.Handle("/posts-top", http.HandlerFunc(app.GetPostsTop))
	mux.Handle("/users-top", http.HandlerFunc(app.GetUsersTop))

	mux.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL("swagger/swagger/doc.json")))

	return &http.Server{
		Addr:    cfg.Port,
		Handler: mux,
	}
}

func InvokeServer(lc fx.Lifecycle, srv *http.Server) error {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				logger.Info("starting server on " + srv.Addr)
				if err := srv.ListenAndServe(); err != nil {
					logger.Error("error starting server: " + err.Error())
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return nil
}
