package server

// RPC handlers.

import (
	"encoding/json"
	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	nm "github.com/vault-thirteen/SimpleBB/pkg/NM/rpc"
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
		srv.AddNotification,
		srv.AddNotificationS,
		srv.SendNotificationIfPossibleS,
		srv.GetNotification,
		srv.GetNotifications,
		srv.GetNotificationsOnPage,
		srv.GetUnreadNotifications,
		srv.CountUnreadNotifications,
		srv.MarkNotificationAsRead,
		srv.DeleteNotification,
		srv.AddResource,
		srv.GetResource,
		srv.GetResourceValue,
		srv.GetListOfAllResourcesOnPage,
		srv.DeleteResource,
		srv.ProcessSystemEventS,
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
	result = nm.PingResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// Notification.

func (srv *Server) AddNotification(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.AddNotificationParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.AddNotificationResult
	r, re = srv.addNotification(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) AddNotificationS(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.AddNotificationSParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.AddNotificationSResult
	r, re = srv.addNotificationS(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) SendNotificationIfPossibleS(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.SendNotificationIfPossibleSParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.SendNotificationIfPossibleSResult
	r, re = srv.sendNotificationIfPossibleS(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetNotification(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.GetNotificationParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.GetNotificationResult
	r, re = srv.getNotification(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetNotifications(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.GetNotificationsParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.GetNotificationsResult
	r, re = srv.getNotifications(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetNotificationsOnPage(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.GetNotificationsOnPageParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.GetNotificationsOnPageResult
	r, re = srv.getNotificationsOnPage(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetUnreadNotifications(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.GetUnreadNotificationsParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.GetUnreadNotificationsResult
	r, re = srv.getUnreadNotifications(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) CountUnreadNotifications(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.CountUnreadNotificationsParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.CountUnreadNotificationsResult
	r, re = srv.countUnreadNotifications(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) MarkNotificationAsRead(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.MarkNotificationAsReadParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.MarkNotificationAsReadResult
	r, re = srv.markNotificationAsRead(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) DeleteNotification(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.DeleteNotificationParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.DeleteNotificationResult
	r, re = srv.deleteNotification(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

// Resource.

func (srv *Server) AddResource(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.AddResourceParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.AddResourceResult
	r, re = srv.addResource(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetResource(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.GetResourceParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.GetResourceResult
	r, re = srv.getResource(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetResourceValue(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.GetResourceValueParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.GetResourceValueResult
	r, re = srv.getResourceValue(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetListOfAllResourcesOnPage(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.GetListOfAllResourcesOnPageParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.GetListOfAllResourcesOnPageResult
	r, re = srv.getListOfAllResourcesOnPage(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) DeleteResource(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.DeleteResourceParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.DeleteResourceResult
	r, re = srv.deleteResource(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

// Other.

func (srv *Server) ProcessSystemEventS(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.ProcessSystemEventSParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.ProcessSystemEventSResult
	r, re = srv.processSystemEventS(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetDKey(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.GetDKeyParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.GetDKeyResult
	r, re = srv.getDKey(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) ShowDiagnosticData(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.ShowDiagnosticDataParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.ShowDiagnosticDataResult
	r, re = srv.showDiagnosticData()
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) Test(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.TestParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.TestResult
	r, re = srv.test(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}
