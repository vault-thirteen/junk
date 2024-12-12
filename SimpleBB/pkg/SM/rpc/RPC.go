package rpc

import (
	"github.com/vault-thirteen/SimpleBB/pkg/SM/models"
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	rpc2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

// Ping.

type PingParams = rpc2.PingParams
type PingResult = rpc2.PingResult

// Subscription.

type AddSubscriptionParams struct {
	rpc2.CommonParams
	ThreadId base2.Id `json:"threadId"`
	UserId   base2.Id `json:"userId"`
}
type AddSubscriptionResult = rpc2.CommonResultWithSuccess

type IsSelfSubscribedParams struct {
	rpc2.CommonParams
	ThreadId base2.Id `json:"threadId"`
}
type IsSelfSubscribedResult struct {
	rpc2.CommonResult
	UserId       base2.Id   `json:"userId"`
	ThreadId     base2.Id   `json:"threadId"`
	IsSubscribed base2.Flag `json:"isSubscribed"`
}

type IsUserSubscribedParams struct {
	rpc2.CommonParams
	UserId   base2.Id `json:"userId"`
	ThreadId base2.Id `json:"threadId"`
}
type IsUserSubscribedResult struct {
	rpc2.CommonResult
	UserId       base2.Id   `json:"userId"`
	ThreadId     base2.Id   `json:"threadId"`
	IsSubscribed base2.Flag `json:"isSubscribed"`
}

type IsUserSubscribedSParams struct {
	rpc2.CommonParams
	rpc2.DKeyParams
	UserId   base2.Id `json:"userId"`
	ThreadId base2.Id `json:"threadId"`
}
type IsUserSubscribedSResult struct {
	rpc2.CommonResult
	UserId       base2.Id   `json:"userId"`
	ThreadId     base2.Id   `json:"threadId"`
	IsSubscribed base2.Flag `json:"isSubscribed"`
}

type CountSelfSubscriptionsParams struct {
	rpc2.CommonParams
}
type CountSelfSubscriptionsResult struct {
	rpc2.CommonResult
	UserSubscriptionsCount base2.Count `json:"userSubscriptionsCount"`
}

type GetSelfSubscriptionsParams struct {
	rpc2.CommonParams
}
type GetSelfSubscriptionsResult struct {
	rpc2.CommonResult
	UserSubscriptions *models.UserSubscriptions `json:"userSubscriptions"`
}

type GetSelfSubscriptionsOnPageParams struct {
	rpc2.CommonParams
	Page base2.Count `json:"page"`
}
type GetSelfSubscriptionsOnPageResult struct {
	rpc2.CommonResult
	UserSubscriptions *models.UserSubscriptions `json:"userSubscriptions"`
}

type GetUserSubscriptionsParams struct {
	rpc2.CommonParams
	UserId base2.Id `json:"userId"`
}
type GetUserSubscriptionsResult struct {
	rpc2.CommonResult
	UserSubscriptions *models.UserSubscriptions `json:"userSubscriptions"`
}

type GetUserSubscriptionsOnPageParams struct {
	rpc2.CommonParams
	UserId base2.Id    `json:"userId"`
	Page   base2.Count `json:"page"`
}
type GetUserSubscriptionsOnPageResult struct {
	rpc2.CommonResult
	UserSubscriptions *models.UserSubscriptions `json:"userSubscriptions"`
}

type GetThreadSubscribersSParams struct {
	rpc2.CommonParams
	rpc2.DKeyParams
	ThreadId base2.Id `json:"threadId"`
}
type GetThreadSubscribersSResult struct {
	rpc2.CommonResult
	ThreadSubscriptions *models.ThreadSubscriptionsRecord `json:"threadSubscriptions"`
}

type DeleteSelfSubscriptionParams struct {
	rpc2.CommonParams
	ThreadId base2.Id `json:"threadId"`
}
type DeleteSelfSubscriptionResult = rpc2.CommonResultWithSuccess

type DeleteSubscriptionParams struct {
	rpc2.CommonParams
	ThreadId base2.Id `json:"threadId"`
	UserId   base2.Id `json:"userId"`
}
type DeleteSubscriptionResult = rpc2.CommonResultWithSuccess

type DeleteSubscriptionSParams struct {
	rpc2.CommonParams
	rpc2.DKeyParams
	ThreadId base2.Id `json:"threadId"`
	UserId   base2.Id `json:"userId"`
}
type DeleteSubscriptionSResult = rpc2.CommonResultWithSuccess

type ClearThreadSubscriptionsSParams struct {
	rpc2.CommonParams
	rpc2.DKeyParams
	ThreadId base2.Id `json:"threadId"`
}
type ClearThreadSubscriptionsSResult = rpc2.CommonResultWithSuccess

// Other.

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
