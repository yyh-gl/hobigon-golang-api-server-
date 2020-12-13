// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package test

import (
	"github.com/google/wire"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/dao"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/db"
	"github.com/yyh-gl/hobigon-golang-api-server/app/presentation/http"
	"github.com/yyh-gl/hobigon-golang-api-server/app/usecase"
	"github.com/yyh-gl/hobigon-golang-api-server/cmd/http/di"
)

// Injectors from wire.go:

func InitTestApp() *di.ContainerAPI {
	gormDB := db.NewDB()
	blog := dao.NewBlog(gormDB)
	slack := dao.NewSlack()
	usecaseBlog := usecase.NewBlog(blog, slack)
	httpBlog := http.NewBlog(usecaseBlog)
	birthday := dao.NewBirthday(gormDB)
	usecaseBirthday := usecase.NewBirthday(birthday)
	httpBirthday := http.NewBirthday(usecaseBirthday)
	task := dao.NewTask()
	notification := usecase.NewNotification(task, slack, birthday)
	httpNotification := http.NewNotification(notification)
	logger := app.NewAPILogger()
	containerAPI := &di.ContainerAPI{
		HandlerBlog:         httpBlog,
		HandlerBirthday:     httpBirthday,
		HandlerNotification: httpNotification,
		DB:                  gormDB,
		Logger:              logger,
	}
	return containerAPI
}

// wire.go:

var testAppSet = wire.NewSet(app.APISet, infra.APISet, usecase.APISet, http.WireSet)
