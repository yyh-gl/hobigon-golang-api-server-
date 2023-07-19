package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/app/presentation/rest/middleware"
	"github.com/yyh-gl/hobigon-golang-api-server/cmd/rest/di"
)

func main() {
	app.NewLogger()

	diContainer := initApp()
	defer func() { _ = diContainer.DB.Close() }()

	router := newRouter(diContainer)

	s := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	errCh := make(chan error, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		fmt.Println("========================")
		fmt.Println("Server Start >> http://localhost" + s.Addr)
		fmt.Println("========================")
		middleware.CountUpRunningVersion(app.GetVersion())
		errCh <- s.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		fmt.Println("Error happened:", err.Error())
	case sig := <-sigCh:
		fmt.Println("Signal received:", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		fmt.Println("Graceful shutdown failed:", err.Error())
	}
	middleware.CountDownRunningVersion(app.GetVersion())
	fmt.Println("Server shutdown")
}

func newRouter(diContainer *di.ContainerAPI) *mux.Router {
	r := mux.NewRouter()

	// Preflight handler
	r.PathPrefix("/").Handler(middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})).Methods(http.MethodOptions)

	// Health Check
	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	// Debug Handlers
	debugGetFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	r.HandleFunc(
		middleware.CreateHandlerFuncWithMiddleware(debugGetFunc, "/api/debug", "debug_get"),
	).Methods(http.MethodGet)

	debugPostFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}
	r.HandleFunc(
		middleware.CreateHandlerFuncWithMiddleware(debugPostFunc, "/api/debug", "debug_post"),
	).Methods(http.MethodPost)

	// Blog handlers
	r.HandleFunc(
		middleware.CreateHandlerFuncWithMiddleware(
			diContainer.HandlerBlog.Create,
			"/api/v1/blogs",
			"blog_create",
		),
	).Methods(http.MethodPost)

	r.HandleFunc(
		middleware.CreateHandlerFuncWithMiddleware(
			diContainer.HandlerBlog.Show,
			"/api/v1/blogs/{title}",
			"blog_show",
		),
	).Methods(http.MethodGet)

	r.HandleFunc(
		middleware.CreateHandlerFuncWithMiddleware(
			diContainer.HandlerBlog.Like,
			"/api/v1/blogs/{title}/like",
			"blog_like",
		),
	).Methods(http.MethodPost)

	// Calendar handlers
	r.HandleFunc(
		middleware.CreateHandlerFuncWithMiddleware(
			diContainer.HandlerCalendar.Create,
			"/api/v1/calendars",
			"calendar_create",
		),
	).Methods(http.MethodPost)

	// Notification handlers
	r.HandleFunc(
		middleware.CreateHandlerFuncWithMiddleware(
			diContainer.HandlerNotification.NotifyTodayTasksToSlack,
			"/api/v1/notifications/slack/tasks/today",
			"today_tasks_notification_send",
		),
	).Methods(http.MethodPost)

	r.HandleFunc(
		middleware.CreateHandlerFuncWithMiddleware(
			diContainer.HandlerNotification.NotifyAccessRankingToSlack,
			"/api/v1/notifications/slack/rankings/access",
			"access_ranking_notification_send",
		),
	).Methods(http.MethodPost)

	r.Handle("/metrics", promhttp.Handler())

	return r
}
