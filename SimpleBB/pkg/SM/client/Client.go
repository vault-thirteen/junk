package c

import (
	cc "github.com/vault-thirteen/SimpleBB/pkg/common/models/Client"
)

// List of supported functions.
const (
	// Ping.
	FuncPing = cc.FuncPing

	// Subscription.
	FuncAddSubscription            = "AddSubscription"
	FuncIsSelfSubscribed           = "IsSelfSubscribed"
	FuncIsUserSubscribed           = "IsUserSubscribed"
	FuncIsUserSubscribedS          = "IsUserSubscribedS"
	FuncCountSelfSubscriptions     = "CountSelfSubscriptions"
	FuncGetSelfSubscriptions       = "GetSelfSubscriptions"
	FuncGetSelfSubscriptionsOnPage = "GetSelfSubscriptionsOnPage"
	FuncGetUserSubscriptions       = "GetUserSubscriptions"
	FuncGetUserSubscriptionsOnPage = "GetUserSubscriptionsOnPage"
	FuncGetThreadSubscribersS      = "GetThreadSubscribersS"
	FuncDeleteSelfSubscription     = "DeleteSelfSubscription"
	FuncDeleteSubscription         = "DeleteSubscription"
	FuncDeleteSubscriptionS        = "DeleteSubscriptionS"
	FuncClearThreadSubscriptionsS  = "ClearThreadSubscriptionsS"

	// Other.
	FuncGetDKey            = "GetDKey"
	FuncShowDiagnosticData = cc.FuncShowDiagnosticData
	FuncTest               = "Test"
)
