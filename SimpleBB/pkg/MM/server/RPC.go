package server

// RPC handlers.

import (
	"encoding/json"
	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	mm "github.com/vault-thirteen/SimpleBB/pkg/MM/rpc"
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
		srv.AddSection,
		srv.ChangeSectionName,
		srv.ChangeSectionParent,
		srv.GetSection,
		srv.MoveSectionUp,
		srv.MoveSectionDown,
		srv.DeleteSection,
		srv.AddForum,
		srv.ChangeForumName,
		srv.ChangeForumSection,
		srv.GetForum,
		srv.MoveForumUp,
		srv.MoveForumDown,
		srv.DeleteForum,
		srv.AddThread,
		srv.ChangeThreadName,
		srv.ChangeThreadForum,
		srv.GetThread,
		srv.GetThreadNamesByIds,
		srv.MoveThreadUp,
		srv.MoveThreadDown,
		srv.DeleteThread,
		srv.ThreadExistsS,
		srv.AddMessage,
		srv.ChangeMessageText,
		srv.ChangeMessageThread,
		srv.GetMessage,
		srv.GetLatestMessageOfThread,
		srv.DeleteMessage,
		srv.ListThreadAndMessages,
		srv.ListThreadAndMessagesOnPage,
		srv.ListForumAndThreads,
		srv.ListForumAndThreadsOnPage,
		srv.ListSectionsAndForums,
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
	result = mm.PingResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// Section.

func (srv *Server) AddSection(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.AddSectionParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.AddSectionResult
	r, re = srv.addSection(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) ChangeSectionName(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.ChangeSectionNameParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.ChangeSectionNameResult
	r, re = srv.changeSectionName(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) ChangeSectionParent(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.ChangeSectionParentParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.ChangeSectionParentResult
	r, re = srv.changeSectionParent(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetSection(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.GetSectionParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.GetSectionResult
	r, re = srv.getSection(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) MoveSectionUp(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.MoveSectionUpParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.MoveSectionUpResult
	r, re = srv.moveSectionUp(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) MoveSectionDown(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.MoveSectionDownParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.MoveSectionDownResult
	r, re = srv.moveSectionDown(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) DeleteSection(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.DeleteSectionParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.DeleteSectionResult
	r, re = srv.deleteSection(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

// Forum.

func (srv *Server) AddForum(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.AddForumParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.AddForumResult
	r, re = srv.addForum(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) ChangeForumName(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.ChangeForumNameParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.ChangeForumNameResult
	r, re = srv.changeForumName(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) ChangeForumSection(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.ChangeForumSectionParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.ChangeForumSectionResult
	r, re = srv.changeForumSection(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetForum(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.GetForumParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.GetForumResult
	r, re = srv.getForum(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) MoveForumUp(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.MoveForumUpParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.MoveForumUpResult
	r, re = srv.moveForumUp(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) MoveForumDown(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.MoveForumDownParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.MoveForumDownResult
	r, re = srv.moveForumDown(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) DeleteForum(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.DeleteForumParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.DeleteForumResult
	r, re = srv.deleteForum(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

// Thread.

func (srv *Server) AddThread(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.AddThreadParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.AddThreadResult
	r, re = srv.addThread(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) ChangeThreadName(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.ChangeThreadNameParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.ChangeThreadNameResult
	r, re = srv.changeThreadName(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) ChangeThreadForum(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.ChangeThreadForumParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.ChangeThreadForumResult
	r, re = srv.changeThreadForum(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetThread(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.GetThreadParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.GetThreadResult
	r, re = srv.getThread(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetThreadNamesByIds(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.GetThreadNamesByIdsParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.GetThreadNamesByIdsResult
	r, re = srv.getThreadNamesByIds(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) MoveThreadUp(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.MoveThreadUpParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.MoveThreadUpResult
	r, re = srv.moveThreadUp(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) MoveThreadDown(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.MoveThreadDownParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.MoveThreadDownResult
	r, re = srv.moveThreadDown(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) DeleteThread(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.DeleteThreadParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.DeleteThreadResult
	r, re = srv.deleteThread(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) ThreadExistsS(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.ThreadExistsSParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.ThreadExistsSResult
	r, re = srv.threadExistsS(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

// Message.

func (srv *Server) AddMessage(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.AddMessageParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.AddMessageResult
	r, re = srv.addMessage(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) ChangeMessageText(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.ChangeMessageTextParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.ChangeMessageTextResult
	r, re = srv.changeMessageText(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) ChangeMessageThread(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.ChangeMessageThreadParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.ChangeMessageThreadResult
	r, re = srv.changeMessageThread(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetMessage(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.GetMessageParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.GetMessageResult
	r, re = srv.getMessage(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetLatestMessageOfThread(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.GetLatestMessageOfThreadParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.GetLatestMessageOfThreadResult
	r, re = srv.getLatestMessageOfThread(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) DeleteMessage(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.DeleteMessageParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.DeleteMessageResult
	r, re = srv.deleteMessage(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

// Composite objects.

func (srv *Server) ListThreadAndMessages(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.ListThreadAndMessagesParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.ListThreadAndMessagesResult
	r, re = srv.listThreadAndMessages(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) ListThreadAndMessagesOnPage(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.ListThreadAndMessagesOnPageParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.ListThreadAndMessagesOnPageResult
	r, re = srv.listThreadAndMessagesOnPage(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) ListForumAndThreads(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.ListForumAndThreadsParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.ListForumAndThreadsResult
	r, re = srv.listForumAndThreads(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) ListForumAndThreadsOnPage(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.ListForumAndThreadsOnPageParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.ListForumAndThreadsOnPageResult
	r, re = srv.listForumAndThreadsOnPage(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) ListSectionsAndForums(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.ListSectionsAndForumsParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.ListSectionsAndForumsResult
	r, re = srv.listSectionsAndForums(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

// Other.

func (srv *Server) GetDKey(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.GetDKeyParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.GetDKeyResult
	r, re = srv.getDKey(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) ShowDiagnosticData(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.ShowDiagnosticDataParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.ShowDiagnosticDataResult
	r, re = srv.showDiagnosticData()
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) Test(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *mm.TestParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *mm.TestResult
	r, re = srv.test(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}
