// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/dao"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/db"
	"github.com/yyh-gl/hobigon-golang-api-server/app/presentation/cli"
	"github.com/yyh-gl/hobigon-golang-api-server/app/usecase"
	"github.com/yyh-gl/hobigon-golang-api-server/cmd/rest/di"
)

// Injectors from wire.go:

func initApp() *di.ContainerCLI {
	task := dao.NewTask()
	slack := dao.NewSlack()
	notification := usecase.NewNotification(task, slack)
	cliNotification := cli.NewNotification(notification)
	gormDB := db.NewDB()
	containerCLI := &di.ContainerCLI{
		HandlerNotification: cliNotification,
		DB:                  gormDB,
	}
	return containerCLI
}

// wire.go:

var appSet = wire.NewSet(infra.CLISet, usecase.CLISet, cli.WireSet)
