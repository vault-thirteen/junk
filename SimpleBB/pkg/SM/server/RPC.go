package server

// RPC handlers.

import (
	"encoding/json"
	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	sm "github.com/vault-thirteen/SimpleBB/pkg/SM/rpc"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	cs "github.com/vault-thirteen/SimpleBB/pkg/common/models/settings"
)

func (srv *Server) initRpc() (err error) {
	rpcDurationFieldName := cs.RpcDurationFieldName
	rpcRequestIdFieldName := cs.RpcRequestIdFieldName

	ps := &jrm1.ProcessorSettings{
		CatchExceptions:    true,
		LogExceptions:      true,
		CountRequests:      true,
		DurationFieldName:  &rpcDurationFieldName,
		RequestIdFieldName: &rpcRequestIdFieldName,
	}

	srv.js, err = jrm1.NewProcessor(ps)
	if err != nil {
		return err
	}

	fns := []jrm1.RpcFunction{
		srv.Ping,
		srv.AddSubscription,
		srv.IsSelfSubscribed,
		srv.IsUserSubscribed,
		srv.IsUserSubscribedS,
		srv.CountSelfSubscriptions,
		srv.GetSelfSubscriptions,
		srv.GetSelfSubscriptionsOnPage,
		srv.GetUserSubscriptions,
		srv.GetUserSubscriptionsOnPage,
		srv.GetThreadSubscribersS,
		srv.DeleteSelfSubscription,
		srv.DeleteSubscription,
		srv.DeleteSubscriptionS,
		srv.ClearThreadSubscriptionsS,
		srv.GetDKey,
		srv.ShowDiagnosticData,
		srv.Test,
	}

	for _, fn := range fns {
		err = srv.js.AddFunc(fn)
		if err != nil {
			return err
		}
	}

	return nil
}

// Ping.

func (srv *Server) Ping(_ *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	result = sm.PingResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// Subscription.

func (srv *Server) AddSubscription(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *sm.AddSubscriptionParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *sm.AddSubscriptionResult
	r, re = srv.addSubscription(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) IsSelfSubscribed(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *sm.IsSelfSubscribedParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *sm.IsSelfSubscribedResult
	r, re = srv.isSelfSubscribed(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) IsUserSubscribed(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *sm.IsUserSubscribedParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *sm.IsUserSubscribedResult
	r, re = srv.isUserSubscribed(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) IsUserSubscribedS(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *sm.IsUserSubscribedSParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *sm.IsUserSubscribedSResult
	r, re = srv.isUserSubscribedS(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) CountSelfSubscriptions(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *sm.CountSelfSubscriptionsParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *sm.CountSelfSubscriptionsResult
	r, re = srv.countSelfSubscriptions(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetSelfSubscriptions(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *sm.GetSelfSubscriptionsParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *sm.GetSelfSubscriptionsResult
	r, re = srv.getSelfSubscriptions(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetSelfSubscriptionsOnPage(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *sm.GetSelfSubscriptionsOnPageParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *sm.GetSelfSubscriptionsOnPageResult
	r, re = srv.getSelfSubscriptionsOnPage(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetUserSubscriptions(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *sm.GetUserSubscriptionsParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *sm.GetUserSubscriptionsResult
	r, re = srv.getUserSubscriptions(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetUserSubscriptionsOnPage(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *sm.GetUserSubscriptionsOnPageParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *sm.GetUserSubscriptionsOnPageResult
	r, re = srv.getUserSubscriptionsOnPage(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetThreadSubscribersS(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *sm.GetThreadSubscribersSParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *sm.GetThreadSubscribersSResult
	r, re = srv.getThreadSubscribersS(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) DeleteSelfSubscription(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *sm.DeleteSelfSubscriptionParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *sm.DeleteSelfSubscriptionResult
	r, re = srv.deleteSelfSubscription(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) DeleteSubscription(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *sm.DeleteSubscriptionParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *sm.DeleteSubscriptionResult
	r, re = srv.deleteSubscription(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) DeleteSubscriptionS(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *sm.DeleteSubscriptionSParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *sm.DeleteSubscriptionSResult
	r, re = srv.deleteSubscriptionS(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) ClearThreadSubscriptionsS(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *sm.ClearThreadSubscriptionsSParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *sm.ClearThreadSubscriptionsSResult
	r, re = srv.clearThreadSubscriptionsS(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

// Other.

func (srv *Server) GetDKey(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *sm.GetDKeyParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *sm.GetDKeyResult
	r, re = srv.getDKey(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) ShowDiagnosticData(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *sm.ShowDiagnosticDataParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *sm.ShowDiagnosticDataResult
	r, re = srv.showDiagnosticData()
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) Test(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *sm.TestParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *sm.TestResult
	r, re = srv.test(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}
