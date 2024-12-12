package rpc

import (
	"github.com/vault-thirteen/SimpleBB/pkg/NM/models"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived2"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/UidList"
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	rpc2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

// Ping.

type PingParams = rpc2.PingParams
type PingResult = rpc2.PingResult

// Notification.

type AddNotificationParams struct {
	rpc2.CommonParams
	UserId base2.Id   `json:"userId"`
	Text   base2.Text `json:"text"`
}
type AddNotificationResult struct {
	rpc2.CommonResult

	// ID of the created notification.
	NotificationId base2.Id `json:"notificationId"`
}

type AddNotificationSParams struct {
	rpc2.CommonParams
	rpc2.DKeyParams
	UserId base2.Id   `json:"userId"`
	Text   base2.Text `json:"text"`
}
type AddNotificationSResult struct {
	rpc2.CommonResult

	// ID of the created notification.
	NotificationId base2.Id `json:"notificationId"`
}

type SendNotificationIfPossibleSParams struct {
	rpc2.CommonParams
	rpc2.DKeyParams
	UserId base2.Id   `json:"userId"`
	Text   base2.Text `json:"text"`
}
type SendNotificationIfPossibleSResult struct {
	rpc2.CommonResult

	// ID and status of the created notification when it is available.
	IsSent         base2.Flag `json:"isSent"`
	NotificationId base2.Id   `json:"notificationId"`
}

type GetNotificationParams struct {
	rpc2.CommonParams
	NotificationId base2.Id `json:"notificationId"`
}
type GetNotificationResult struct {
	rpc2.CommonResult
	Notification *models.Notification `json:"notification"`
}

type GetNotificationsParams struct {
	rpc2.CommonParams
}
type GetNotificationsResult struct {
	rpc2.CommonResult
	NotificationIds *ul.UidList           `json:"notificationIds"`
	Notifications   []models.Notification `json:"notifications"`
}

type GetNotificationsOnPageParams struct {
	rpc2.CommonParams
	Page base2.Count `json:"page"`
}
type GetNotificationsOnPageResult struct {
	rpc2.CommonResult
	NotificationsOnPage *models.NotificationsOnPage `json:"nop"`
}

type GetUnreadNotificationsParams struct {
	rpc2.CommonParams
}
type GetUnreadNotificationsResult struct {
	rpc2.CommonResult
	NotificationIds *ul.UidList           `json:"notificationIds"`
	Notifications   []models.Notification `json:"notifications"`
}

type CountUnreadNotificationsParams struct {
	rpc2.CommonParams
}
type CountUnreadNotificationsResult struct {
	rpc2.CommonResult
	UNC base2.Count `json:"unc"`
}

type MarkNotificationAsReadParams struct {
	rpc2.CommonParams

	// Identifier of a notification.
	NotificationId base2.Id `json:"notificationId"`
}
type MarkNotificationAsReadResult = rpc2.CommonResultWithSuccess

type DeleteNotificationParams struct {
	rpc2.CommonParams
	NotificationId base2.Id `json:"notificationId"`
}
type DeleteNotificationResult = rpc2.CommonResultWithSuccess

// Resource.

type AddResourceParams struct {
	rpc2.CommonParams
	Resource any `json:"resource"`
}
type AddResourceResult struct {
	rpc2.CommonResult

	// ID of the created resource.
	ResourceId base2.Id `json:"resourceId"`
}

type GetResourceParams struct {
	rpc2.CommonParams
	ResourceId base2.Id `json:"resourceId"`
}
type GetResourceResult struct {
	rpc2.CommonResult
	Resource derived2.IResource `json:"resource"`
}

type GetResourceValueParams struct {
	rpc2.CommonParams
	ResourceId base2.Id `json:"resourceId"`
}
type GetResourceValueResult struct {
	rpc2.CommonResult
	Resource models.ResourceWithValue `json:"resource"`
}

type GetListOfAllResourcesOnPageParams struct {
	rpc2.CommonParams
	Page base2.Count `json:"page"`
}
type GetListOfAllResourcesOnPageResult struct {
	rpc2.CommonResult
	ResourcesOnPage *models.ResourcesOnPage `json:"rop"`
}

type DeleteResourceParams struct {
	rpc2.CommonParams
	ResourceId base2.Id `json:"resourceId"`
}
type DeleteResourceResult = rpc2.CommonResultWithSuccess

// Other.

type ProcessSystemEventSParams struct {
	rpc2.CommonParams
	rpc2.DKeyParams
	SystemEventData derived2.ISystemEventData `json:"systemEventData"`
}
type ProcessSystemEventSResult = rpc2.CommonResultWithSuccess

type GetDKeyParams struct {
	rpc2.CommonParams
}
type GetDKeyResult struct {
	rpc2.CommonResult
	DKey base2.Text `json:"dKey"`
}

type ShowDiagnosticDataParams struct{}
type ShowDiagnosticDataResult struct {
	rpc2.CommonResult
	rpc2.RequestsCount
}

type TestParams struct {
	N uint `json:"n"`
}
type TestResult struct {
	rpc2.CommonResult
}
