package c

import (
	cc "github.com/vault-thirteen/SimpleBB/pkg/common/models/Client"
)

// List of supported functions.
const (
	// Ping.
	FuncPing = cc.FuncPing

	// Notification.
	FuncAddNotification             = "AddNotification"
	FuncAddNotificationS            = "AddNotificationS"
	FuncSendNotificationIfPossibleS = "SendNotificationIfPossibleS"
	FuncGetNotification             = "GetNotification"
	FuncGetNotifications            = "GetNotifications"
	FuncGetNotificationsOnPage      = "GetNotificationsOnPage"
	FuncGetUnreadNotifications      = "GetUnreadNotifications"
	FuncCountUnreadNotifications    = "CountUnreadNotifications"
	FuncMarkNotificationAsRead      = "MarkNotificationAsRead"
	FuncDeleteNotification          = "DeleteNotification"

	// Resource.
	FuncAddResource                 = "AddResource"
	FuncGetResource                 = "GetResource"
	FuncGetResourceValue            = "GetResourceValue"
	FuncGetListOfAllResourcesOnPage = "GetListOfAllResourcesOnPage"
	FuncDeleteResource              = "DeleteResource"

	// Other.
	FuncProcessSystemEventS = "ProcessSystemEventS"
	FuncGetDKey             = "GetDKey"
	FuncShowDiagnosticData  = cc.FuncShowDiagnosticData
	FuncTest                = "Test"
)
