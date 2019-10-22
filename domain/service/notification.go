package service

import (
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
)

// NotificationService は通知サービスのインターフェース
type NotificationService interface {
	SendTodayBirthdayToSlack(model.Birthday) error
}
