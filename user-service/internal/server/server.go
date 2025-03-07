package server

import (
	"context"
	"go.uber.org/fx"
	"net/http"
	"social-network/user-service/internal/app"
	"social-network/user-service/internal/config"
	"social-network/user-service/internal/logger"
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

	return &http.Server{
		Addr:    cfg.ServerAddr,
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
