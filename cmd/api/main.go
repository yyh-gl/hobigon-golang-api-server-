package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
)

func main() {
	// システム共通で使用するものを用意
	//  -> logger, DB
	app.Init(app.APILogFilename)
	defer func() { _ = app.DB.Close() }()

	// 依存関係を定義
	blogHandler := initBlogHandler()
	birthdayHandler := initBirthdayHandler()
	notificationHandler := initNotificationHandler()

	// ルーティング設定
	r := httprouter.New()
	r.GlobalOPTIONS = wrapHandler(preflightHandler)

	// ブログ関連のAPI
	r.HandlerFunc(http.MethodPost, "/api/v1/blogs", wrapHandler(blogHandler.Create))
	r.HandlerFunc(http.MethodGet, "/api/v1/blogs/:title", wrapHandler(blogHandler.Show))
	r.HandlerFunc(http.MethodPost, "/api/v1/blogs/:title/like", wrapHandler(blogHandler.Like))

	// 誕生日関連のAPI
	r.HandlerFunc(http.MethodPost, "/api/v1/birthday", wrapHandler(birthdayHandler.Create))

	// 通知系API
	r.HandlerFunc(http.MethodPost, "/api/v1/notifications/slack/tasks/today", wrapHandler(notificationHandler.NotifyTodayTasksToSlack))
	//// TODO: 誕生日の人が複数いたときに対応
	r.HandlerFunc(http.MethodPost, "/api/v1/notifications/slack/birthdays/today", wrapHandler(notificationHandler.NotifyTodayBirthdayToSlack))
	r.HandlerFunc(http.MethodPost, "/api/v1/notifications/slack/rankings/access", wrapHandler(notificationHandler.NotifyAccessRankingToSlack))

	fmt.Println("========================")
	fmt.Println("Server Start >> http://localhost:3000")
	fmt.Println(" ↳  Log File -> " + os.Getenv("LOG_PATH") + "/" + app.APILogFilename)
	fmt.Println("========================")
	app.Logger.Print("Server Start")
	app.Logger.Fatal(http.ListenAndServe(":3000", r))
}

// wrapHandler : 全ハンドラー共通処理
func wrapHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエスト内容をログ出力
		// TODO: Body の内容を記録
		app.Logger.Print("[AccessLog] " + r.Method + " " + r.URL.String())

		// CORS用ヘッダーを付与
		switch {
		case app.IsPrd():
			w.Header().Add("Access-Control-Allow-Origin", "https://yyh-gl.github.io")
		case app.IsDev() || app.IsTest():
			w.Header().Add("Access-Control-Allow-Origin", "http://localhost:1313")
			w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3001")
		}
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json;charset=utf-8")

		h.ServeHTTP(w, r)
	}
}

// preflightHandler : preflight用のハンドラー
func preflightHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
