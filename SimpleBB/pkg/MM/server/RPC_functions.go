package server

import (
	"fmt"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/rpc"
	rpc2 "github.com/vault-thirteen/SimpleBB/pkg/MM/rpc"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived2"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/UidList"
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	ev "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/EnumValue"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/ForumAndThreads"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/SectionChildType"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/SystemEvent"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/SystemEventData"
	set "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/SystemEventType"
	tam "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/ThreadAndMessages"
	rpc3 "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	c "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	"sync"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	mm "github.com/vault-thirteen/SimpleBB/pkg/MM/models"
)

// RPC functions.

// Section.

// addSection inserts a new section as a root section or as a sub-section.
func (srv *Server) addSection(p *rpc2.AddSectionParams) (result *rpc2.AddSectionResult, re *jrm1.RpcError) {
	// Check parameters.
	if len(p.Name) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionNameIsNotSet, RpcErrorMsg_SectionNameIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// If parent is not set, the new section is a root section.
	// Only a single root section may exist.
	var err error
	var n base2.Count
	if p.Parent == nil {
		n, err = srv.dbo.CountRootSections()
		if err != nil {
			return nil, srv.databaseError(err)
		}

		if n > 0 {
			return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_RootSectionAlreadyExists, RpcErrorMsg_RootSectionAlreadyExists, nil)
		}

		var insertedSectionId base2.Id
		insertedSectionId, err = srv.dbo.InsertNewSection(p.Parent, p.Name, userRoles.User.GetUserParameters().GetId())
		if err != nil {
			return nil, srv.databaseError(err)
		}

		result = &rpc2.AddSectionResult{
			SectionId: insertedSectionId,
		}

		return result, nil
	}

	// Insert a sub-section.
	// Ensure that a parent exists.
	n, err = srv.dbo.CountSectionsById(*p.Parent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	// Check compatibility.
	var childType derived1.ISectionChildType
	childType, err = srv.dbo.GetSectionChildTypeById(*p.Parent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if childType.AsInt() == sct.SectionChildType_Forum {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_IncompatibleChildType, RpcErrorMsg_IncompatibleChildType, nil)
	}

	if childType.AsInt() == sct.SectionChildType_None {
		err = srv.dbo.SetSectionChildTypeById(*p.Parent, sct.NewSectionChildTypeWithValue(ev.NewEnumValue(sct.SectionChildType_Section)))
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	// Insert a section and link it with its parent.
	var parentChildren *ul.UidList
	parentChildren, err = srv.dbo.GetSectionChildrenById(*p.Parent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var insertedSectionId base2.Id
	insertedSectionId, err = srv.dbo.InsertNewSection(p.Parent, p.Name, userRoles.User.GetUserParameters().GetId())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = parentChildren.AddItem(insertedSectionId, false)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetSectionChildrenById(*p.Parent, parentChildren)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.AddSectionResult{
		SectionId: insertedSectionId,
	}

	return result, nil
}

// changeSectionName renames a section.
func (srv *Server) changeSectionName(p *rpc2.ChangeSectionNameParams) (result *rpc2.ChangeSectionNameResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.SectionId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIdIsNotSet, RpcErrorMsg_SectionIdIsNotSet, nil)
	}

	if len(p.Name) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionNameIsNotSet, RpcErrorMsg_SectionNameIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var n base2.Count
	var err error
	n, err = srv.dbo.CountSectionsById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	err = srv.dbo.SetSectionNameById(p.SectionId, p.Name, userRoles.User.GetUserParameters().GetId())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.ChangeSectionNameResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// changeSectionParent moves a section from an old parent to a new parent.
func (srv *Server) changeSectionParent(p *rpc2.ChangeSectionParentParams) (result *rpc2.ChangeSectionParentResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.SectionId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIdIsNotSet, RpcErrorMsg_SectionIdIsNotSet, nil)
	}

	if p.Parent == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIdIsNotSet, RpcErrorMsg_SectionIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var n base2.Count
	var err error
	n, err = srv.dbo.CountSectionsById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	// Ensure that an old parent exists.
	var oldParent *base2.Id
	oldParent, err = srv.dbo.GetSectionParentById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if oldParent == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_RootSectionCanNotBeMoved, RpcErrorMsg_RootSectionCanNotBeMoved, nil)
	}

	n, err = srv.dbo.CountSectionsById(*oldParent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	// Ensure that a new parent exists.
	n, err = srv.dbo.CountSectionsById(p.Parent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	// Check compatibility.
	var childType derived1.ISectionChildType
	childType, err = srv.dbo.GetSectionChildTypeById(p.Parent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if childType.AsInt() == sct.SectionChildType_Forum {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_IncompatibleChildType, RpcErrorMsg_IncompatibleChildType, nil)
	}

	if childType.AsInt() == sct.SectionChildType_None {
		err = srv.dbo.SetSectionChildTypeById(p.Parent, sct.NewSectionChildTypeWithValue(ev.NewEnumValue(sct.SectionChildType_Section)))
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	// Update the moved section.
	err = srv.dbo.SetSectionParentById(p.SectionId, p.Parent, userRoles.User.GetUserParameters().GetId())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Update the new link.
	var childrenR *ul.UidList
	childrenR, err = srv.dbo.GetSectionChildrenById(p.Parent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = childrenR.AddItem(p.SectionId, false)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetSectionChildrenById(p.Parent, childrenR)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Update the old link.
	var childrenL *ul.UidList
	childrenL, err = srv.dbo.GetSectionChildrenById(*oldParent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = childrenL.RemoveItem(p.SectionId)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetSectionChildrenById(*oldParent, childrenL)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Clear the child type if the old parent becomes empty.
	if childrenL.Size() == 0 {
		err = srv.dbo.SetSectionChildTypeById(*oldParent, sct.NewSectionChildTypeWithValue(ev.NewEnumValue(sct.SectionChildType_None)))
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	result = &rpc2.ChangeSectionParentResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// getSection reads a section.
func (srv *Server) getSection(p *rpc2.GetSectionParams) (result *rpc2.GetSectionResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.SectionId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIdIsNotSet, RpcErrorMsg_SectionIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read the section.
	var section derived2.ISection
	var err error
	section, err = srv.dbo.GetSectionById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if section == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	result = &rpc2.GetSectionResult{
		Section: section,
	}

	return result, nil
}

// moveSectionUp moves a section up by one position if possible.
func (srv *Server) moveSectionUp(p *rpc2.MoveSectionUpParams) (result *rpc2.MoveSectionUpResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.SectionId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIdIsNotSet, RpcErrorMsg_SectionIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Check existence of the moved section.
	var n base2.Count
	var err error
	n, err = srv.dbo.CountSectionsById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	// Get the section which is being moved.
	var section derived2.ISection
	section, err = srv.dbo.GetSectionById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if section == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}
	if section.GetParent() == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_RootSectionCanNotBeMoved, RpcErrorMsg_RootSectionCanNotBeMoved, nil)
	}

	// Get the parent section.
	var parent derived2.ISection
	parent, err = srv.dbo.GetSectionById(*section.GetParent())
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if parent == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}
	if parent.GetChildren() == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	// Check compatibility.
	if parent.GetChildType().AsInt() != sct.SectionChildType_Section {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_IncompatibleChildType, RpcErrorMsg_IncompatibleChildType, nil)
	}

	err = parent.GetChildren().MoveItemUp(p.SectionId)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetSectionChildrenById(parent.GetId(), parent.GetChildren())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.MoveSectionUpResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// moveSectionDown moves a section down by one position if possible.
func (srv *Server) moveSectionDown(p *rpc2.MoveSectionDownParams) (result *rpc2.MoveSectionDownResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.SectionId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIdIsNotSet, RpcErrorMsg_SectionIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Check existence of the moved section.
	var n base2.Count
	var err error
	n, err = srv.dbo.CountSectionsById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	// Get the section which is being moved.
	var section derived2.ISection
	section, err = srv.dbo.GetSectionById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if section == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}
	if section.GetParent() == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_RootSectionCanNotBeMoved, RpcErrorMsg_RootSectionCanNotBeMoved, nil)
	}

	// Get the parent section.
	var parent derived2.ISection
	parent, err = srv.dbo.GetSectionById(*section.GetParent())
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if parent == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}
	if parent.GetChildren() == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	// Check compatibility.
	if parent.GetChildType().AsInt() != sct.SectionChildType_Section {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_IncompatibleChildType, RpcErrorMsg_IncompatibleChildType, nil)
	}

	err = parent.GetChildren().MoveItemDown(p.SectionId)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetSectionChildrenById(parent.GetId(), parent.GetChildren())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.MoveSectionDownResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// deleteSection removes a section.
func (srv *Server) deleteSection(p *rpc2.DeleteSectionParams) (result *rpc2.DeleteSectionResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.SectionId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIdIsNotSet, RpcErrorMsg_SectionIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Read the section.
	var section derived2.ISection
	var err error
	section, err = srv.dbo.GetSectionById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if section == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	var isRootSection = false
	if section.GetParent() == nil {
		isRootSection = true
	}

	// Check for derived1.
	if section.GetChildren().Size() > 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionHasChildren, RpcErrorMsg_SectionHasChildren, nil)
	}

	// Update the link.
	if !isRootSection {
		var linkSections *ul.UidList
		linkSections, err = srv.dbo.GetSectionChildrenById(*section.GetParent())
		if err != nil {
			return nil, srv.databaseError(err)
		}

		err = linkSections.RemoveItem(p.SectionId)
		if err != nil {
			srv.logError(err)
			return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
		}

		err = srv.dbo.SetSectionChildrenById(*section.GetParent(), linkSections)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		// Clear the child type if the old parent becomes empty.
		if linkSections.Size() == 0 {
			err = srv.dbo.SetSectionChildTypeById(*section.GetParent(), sct.NewSectionChildTypeWithValue(ev.NewEnumValue(sct.SectionChildType_None)))
			if err != nil {
				return nil, srv.databaseError(err)
			}
		}
	}

	// Delete the section.
	err = srv.dbo.DeleteSectionById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.DeleteSectionResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// Forum.

// addForum inserts a new forum into a section.
func (srv *Server) addForum(p *rpc2.AddForumParams) (result *rpc2.AddForumResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.SectionId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIdIsNotSet, RpcErrorMsg_SectionIdIsNotSet, nil)
	}

	if len(p.Name) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumNameIsNotSet, RpcErrorMsg_ForumNameIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Ensure that a section exists.
	n, err := srv.dbo.CountSectionsById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	// Check compatibility.
	var childType derived1.ISectionChildType
	childType, err = srv.dbo.GetSectionChildTypeById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if childType.AsInt() == sct.SectionChildType_Section {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_IncompatibleChildType, RpcErrorMsg_IncompatibleChildType, nil)
	}

	if childType.AsInt() == sct.SectionChildType_None {
		err = srv.dbo.SetSectionChildTypeById(p.SectionId, sct.NewSectionChildTypeWithValue(ev.NewEnumValue(sct.SectionChildType_Forum)))
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	// Insert a forum and link it with its section.
	var parentChildren *ul.UidList
	parentChildren, err = srv.dbo.GetSectionChildrenById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var insertedForumId base2.Id
	insertedForumId, err = srv.dbo.InsertNewForum(p.SectionId, p.Name, userRoles.User.GetUserParameters().GetId())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = parentChildren.AddItem(insertedForumId, false)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetSectionChildrenById(p.SectionId, parentChildren)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.AddForumResult{
		ForumId: insertedForumId,
	}

	return result, nil
}

// changeForumName renames a forum.
func (srv *Server) changeForumName(p *rpc2.ChangeForumNameParams) (result *rpc2.ChangeForumNameResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ForumId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIdIsNotSet, RpcErrorMsg_ForumIdIsNotSet, nil)
	}

	if len(p.Name) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumNameIsNotSet, RpcErrorMsg_ForumNameIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var n base2.Count
	var err error
	n, err = srv.dbo.CountForumsById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	err = srv.dbo.SetForumNameById(p.ForumId, p.Name, userRoles.User.GetUserParameters().GetId())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.ChangeForumNameResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// changeForumSection moves a forum from an old section to a new section.
func (srv *Server) changeForumSection(p *rpc2.ChangeForumSectionParams) (result *rpc2.ChangeForumSectionResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ForumId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIdIsNotSet, RpcErrorMsg_ForumIdIsNotSet, nil)
	}

	if p.SectionId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIdIsNotSet, RpcErrorMsg_SectionIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var n base2.Count
	var err error
	n, err = srv.dbo.CountForumsById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	// Ensure that an old section exists.
	var oldParent base2.Id
	oldParent, err = srv.dbo.GetForumSectionById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	n, err = srv.dbo.CountSectionsById(oldParent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	// Ensure that a new section exists.
	n, err = srv.dbo.CountSectionsById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	// Check compatibility.
	var childType derived1.ISectionChildType
	childType, err = srv.dbo.GetSectionChildTypeById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if childType.AsInt() == sct.SectionChildType_Section {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_IncompatibleChildType, RpcErrorMsg_IncompatibleChildType, nil)
	}

	if childType.AsInt() == sct.SectionChildType_None {
		err = srv.dbo.SetSectionChildTypeById(p.SectionId, sct.NewSectionChildTypeWithValue(ev.NewEnumValue(sct.SectionChildType_Forum)))
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	// Update the moved forum.
	err = srv.dbo.SetForumSectionById(p.ForumId, p.SectionId, userRoles.User.GetUserParameters().GetId())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Update the new link.
	var childrenR *ul.UidList
	childrenR, err = srv.dbo.GetSectionChildrenById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = childrenR.AddItem(p.ForumId, false)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetSectionChildrenById(p.SectionId, childrenR)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Update the old link.
	var childrenL *ul.UidList
	childrenL, err = srv.dbo.GetSectionChildrenById(oldParent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = childrenL.RemoveItem(p.ForumId)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetSectionChildrenById(oldParent, childrenL)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Clear the child type if the old section becomes empty.
	if childrenL.Size() == 0 {
		err = srv.dbo.SetSectionChildTypeById(oldParent, sct.NewSectionChildTypeWithValue(ev.NewEnumValue(sct.SectionChildType_None)))
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	result = &rpc2.ChangeForumSectionResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// getForum reads a forum.
func (srv *Server) getForum(p *rpc2.GetForumParams) (result *rpc2.GetForumResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ForumId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIdIsNotSet, RpcErrorMsg_ForumIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read the forum.
	var forum derived2.IForum
	var err error
	forum, err = srv.dbo.GetForumById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if forum == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	result = &rpc2.GetForumResult{
		Forum: forum,
	}

	return result, nil
}

// moveForumUp moves a forum up by one position if possible.
func (srv *Server) moveForumUp(p *rpc2.MoveForumUpParams) (result *rpc2.MoveForumUpResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ForumId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIdIsNotSet, RpcErrorMsg_ForumIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Check existence of the moved forum.
	var n base2.Count
	var err error
	n, err = srv.dbo.CountForumsById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	// Get the forum which is being moved.
	var forum derived2.IForum
	forum, err = srv.dbo.GetForumById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if forum == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	// Get the parent section.
	var parent derived2.ISection
	parent, err = srv.dbo.GetSectionById(forum.GetSectionId())
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if parent == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}
	if parent.GetChildren() == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	// Check compatibility.
	if parent.GetChildType().AsInt() != sct.SectionChildType_Forum {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_IncompatibleChildType, RpcErrorMsg_IncompatibleChildType, nil)
	}

	err = parent.GetChildren().MoveItemUp(p.ForumId)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetSectionChildrenById(parent.GetId(), parent.GetChildren())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.MoveForumUpResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// moveForumDown moves a forum down by one position if possible.
func (srv *Server) moveForumDown(p *rpc2.MoveForumDownParams) (result *rpc2.MoveForumDownResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ForumId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIdIsNotSet, RpcErrorMsg_ForumIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Check existence of the moved forum.
	var n base2.Count
	var err error
	n, err = srv.dbo.CountForumsById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	// Get the forum which is being moved.
	var forum derived2.IForum
	forum, err = srv.dbo.GetForumById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if forum == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	// Get the parent section.
	var parent derived2.ISection
	parent, err = srv.dbo.GetSectionById(forum.GetSectionId())
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if parent == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}
	if parent.GetChildren() == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	// Check compatibility.
	if parent.GetChildType().AsInt() != sct.SectionChildType_Forum {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_IncompatibleChildType, RpcErrorMsg_IncompatibleChildType, nil)
	}

	err = parent.GetChildren().MoveItemDown(p.ForumId)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetSectionChildrenById(parent.GetId(), parent.GetChildren())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.MoveForumDownResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// deleteForum removes a forum.
func (srv *Server) deleteForum(p *rpc2.DeleteForumParams) (result *rpc2.DeleteForumResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ForumId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIdIsNotSet, RpcErrorMsg_ForumIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Read the forum.
	var forum derived2.IForum
	var err error
	forum, err = srv.dbo.GetForumById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if forum == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	// Check for threads.
	if forum.GetThreads().Size() > 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumHasThreads, RpcErrorMsg_ForumHasThreads, nil)
	}

	// Update the link.
	var linkChildren *ul.UidList
	linkChildren, err = srv.dbo.GetSectionChildrenById(forum.GetSectionId())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = linkChildren.RemoveItem(p.ForumId)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetSectionChildrenById(forum.GetSectionId(), linkChildren)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Clear the child type if the old parent becomes empty.
	if linkChildren.Size() == 0 {
		err = srv.dbo.SetSectionChildTypeById(forum.GetSectionId(), sct.NewSectionChildTypeWithValue(ev.NewEnumValue(sct.SectionChildType_None)))
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	// Delete the forum.
	err = srv.dbo.DeleteForumById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.DeleteForumResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// Thread.

// addThread inserts a new thread into a forum.
func (srv *Server) addThread(p *rpc2.AddThreadParams) (result *rpc2.AddThreadResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ForumId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIdIsNotSet, RpcErrorMsg_ForumIdIsNotSet, nil)
	}

	if len(p.Name) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadNameIsNotSet, RpcErrorMsg_ThreadNameIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsAuthor {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Ensure that a forum exists.
	var err error
	var n base2.Count
	n, err = srv.dbo.CountForumsById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	// Insert a thread and link it with its forum.
	var parentThreads *ul.UidList
	parentThreads, err = srv.dbo.GetForumThreadsById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var insertedThreadId base2.Id
	insertedThreadId, err = srv.dbo.InsertNewThread(p.ForumId, p.Name, userRoles.User.GetUserParameters().GetId())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = parentThreads.AddItem(insertedThreadId, srv.settings.SystemSettings.NewThreadsAtTop.AsBool())
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetForumThreadsById(p.ForumId, parentThreads)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.AddThreadResult{
		ThreadId: insertedThreadId,
	}

	return result, nil
}

// changeThreadName renames a thread.
func (srv *Server) changeThreadName(p *rpc2.ChangeThreadNameParams) (result *rpc2.ChangeThreadNameResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	if len(p.Name) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadNameIsNotSet, RpcErrorMsg_ThreadNameIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	seData := sed.NewSystemEventDataWithValue(
		set.NewSystemEventTypeWithValue(ev.NewEnumValue(set.SystemEventType_ThreadNameChange)),
		&p.ThreadId,
		nil,
		userRoles.User.GetUserParameters().GetIdPtr(),
		nil,
	)

	se, err := cm.NewSystemEventWithData(seData)
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_SystemEvent, c.RpcErrorMsg_SystemEvent, nil)
	}

	re = srv.reportSystemEvent(se)
	if re != nil {
		return nil, re
	}

	re = srv.changeThreadNameH(p.ThreadId, p.Name, userRoles.User.GetUserParameters().GetId())
	if re != nil {
		return nil, re
	}

	result = &rpc2.ChangeThreadNameResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// changeThreadForum moves a thread from an old forum to a new forum.
func (srv *Server) changeThreadForum(p *rpc2.ChangeThreadForumParams) (result *rpc2.ChangeThreadForumResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	if p.ForumId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIdIsNotSet, RpcErrorMsg_ForumIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	seData := sed.NewSystemEventDataWithValue(
		set.NewSystemEventTypeWithValue(ev.NewEnumValue(set.SystemEventType_ThreadParentChange)),
		&p.ThreadId,
		nil,
		userRoles.User.GetUserParameters().GetIdPtr(),
		nil,
	)

	se, err := cm.NewSystemEventWithData(seData)
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_SystemEvent, c.RpcErrorMsg_SystemEvent, nil)
	}

	re = srv.reportSystemEvent(se)
	if re != nil {
		return nil, re
	}

	re = srv.changeThreadForumH(p.ThreadId, p.ForumId, userRoles.User.GetUserParameters().GetId())
	if re != nil {
		return nil, re
	}

	result = &rpc2.ChangeThreadForumResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// getThread reads a thread.
func (srv *Server) getThread(p *rpc2.GetThreadParams) (result *rpc2.GetThreadResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read the thread.
	var thread derived2.IThread
	var err error
	thread, err = srv.dbo.GetThreadById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if thread == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	result = &rpc2.GetThreadResult{
		Thread: thread,
	}

	return result, nil
}

// getThreadNamesByIds reads names of threads specified by their IDs.
func (srv *Server) getThreadNamesByIds(p *rpc2.GetThreadNamesByIdsParams) (result *rpc2.GetThreadNamesByIdsResult, re *jrm1.RpcError) {
	// Check parameters.
	if len(p.ThreadIds) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read thread names.
	var threadNames []simple.Name
	var err error
	threadNames, err = srv.dbo.ReadThreadNamesByIds(p.ThreadIds)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.GetThreadNamesByIdsResult{
		ThreadIds:   p.ThreadIds,
		ThreadNames: threadNames,
	}

	return result, nil
}

// moveThreadUp moves a thread up by one position if possible.
func (srv *Server) moveThreadUp(p *rpc2.MoveThreadUpParams) (result *rpc2.MoveThreadUpResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Check existence of the moved thread.
	var n base2.Count
	var err error
	n, err = srv.dbo.CountThreadsById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	// Get the thread which is being moved.
	var thread derived2.IThread
	thread, err = srv.dbo.GetThreadById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if thread == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	// Get the parent forum.
	var parent derived2.IForum
	parent, err = srv.dbo.GetForumById(thread.GetForumId())
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if parent == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}
	if parent.GetThreads() == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	err = parent.GetThreads().MoveItemUp(p.ThreadId)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetForumThreadsById(parent.GetId(), parent.GetThreads())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.MoveThreadUpResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// moveThreadDown moves a thread down by one position if possible.
func (srv *Server) moveThreadDown(p *rpc2.MoveThreadDownParams) (result *rpc2.MoveThreadDownResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Check existence of the moved thread.
	var n base2.Count
	var err error
	n, err = srv.dbo.CountThreadsById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	// Get the thread which is being moved.
	var thread derived2.IThread
	thread, err = srv.dbo.GetThreadById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if thread == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	// Get the parent forum.
	var parent derived2.IForum
	parent, err = srv.dbo.GetForumById(thread.GetForumId())
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if parent == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}
	if parent.GetThreads() == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	err = parent.GetThreads().MoveItemDown(p.ThreadId)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetForumThreadsById(parent.GetId(), parent.GetThreads())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.MoveThreadDownResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// deleteThread removes a thread.
func (srv *Server) deleteThread(p *rpc2.DeleteThreadParams) (result *rpc2.DeleteThreadResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	seData := sed.NewSystemEventDataWithValue(
		set.NewSystemEventTypeWithValue(ev.NewEnumValue(set.SystemEventType_ThreadDeletion)),
		&p.ThreadId,
		nil,
		userRoles.User.GetUserParameters().GetIdPtr(),
		nil,
	)

	se, err := cm.NewSystemEventWithData(seData)
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_SystemEvent, c.RpcErrorMsg_SystemEvent, nil)
	}

	re = srv.reportSystemEvent(se)
	if re != nil {
		return nil, re
	}

	re = srv.deleteThreadH(p.ThreadId)
	if re != nil {
		return nil, re
	}

	result = &rpc2.DeleteThreadResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// threadExistsS checks whether the specified thread exists or not. This method
// is used by the system.
func (srv *Server) threadExistsS(p *rpc2.ThreadExistsSParams) (result *rpc2.ThreadExistsSResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	re = srv.mustBeNoAuth(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check the DKey.
	if !srv.dKeyI.CheckString(p.DKey.ToString()) {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Count threads.
	var n base2.Count
	var err error
	n, err = srv.dbo.CountThreadsById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.ThreadExistsSResult{
		Exists: n == 1,
	}

	return result, nil
}

// Message.

// addMessage inserts a new message into a thread.
func (srv *Server) addMessage(p *rpc2.AddMessageParams) (result *rpc2.AddMessageResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	if len(p.Text) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_MessageTextIsNotSet, RpcErrorMsg_MessageTextIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions (Part I).
	if !userRoles.User.GetUserParameters().GetRoles().IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	result, re = srv.addMessageH(p.ThreadId, p.Text, userRoles)
	if re != nil {
		return nil, re
	}

	seData := sed.NewSystemEventDataWithValue(
		set.NewSystemEventTypeWithValue(ev.NewEnumValue(set.SystemEventType_ThreadNewMessage)),
		&p.ThreadId,
		&result.MessageId,
		userRoles.User.GetUserParameters().GetIdPtr(),
		nil,
	)

	se, err := cm.NewSystemEventWithData(seData)
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_SystemEvent, c.RpcErrorMsg_SystemEvent, nil)
	}

	re = srv.reportSystemEvent(se)
	if re != nil {
		return nil, re
	}

	return result, nil
}

// changeMessageText changes text of a message.
func (srv *Server) changeMessageText(p *rpc2.ChangeMessageTextParams) (result *rpc2.ChangeMessageTextResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.MessageId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_MessageIdIsNotSet, RpcErrorMsg_MessageIdIsNotSet, nil)
	}

	if len(p.Text) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_MessageTextIsNotSet, RpcErrorMsg_MessageTextIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	var initialMessage derived2.IMessage
	initialMessage, re = srv.changeMessageTextH(p.MessageId, p.Text, userRoles)
	if re != nil {
		return nil, re
	}

	seData := sed.NewSystemEventDataWithValue(
		set.NewSystemEventTypeWithValue(ev.NewEnumValue(set.SystemEventType_ThreadMessageEdit)),
		initialMessage.GetThreadIdPtr(),
		&p.MessageId,
		userRoles.User.GetUserParameters().GetIdPtr(),
		nil,
	)

	se, err := cm.NewSystemEventWithData(seData)
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_SystemEvent, c.RpcErrorMsg_SystemEvent, nil)
	}

	re = srv.reportSystemEvent(se)
	if re != nil {
		return nil, re
	}

	seData = sed.NewSystemEventDataWithValue(
		set.NewSystemEventTypeWithValue(ev.NewEnumValue(set.SystemEventType_MessageTextEdit)),
		initialMessage.GetThreadIdPtr(),
		&p.MessageId,
		userRoles.User.GetUserParameters().GetIdPtr(),
		initialMessage.GetEventData().GetCreatorUserIdPtr(),
	)

	se, err = cm.NewSystemEventWithData(seData)
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_SystemEvent, c.RpcErrorMsg_SystemEvent, nil)
	}

	re = srv.reportSystemEvent(se)
	if re != nil {
		return nil, re
	}

	result = &rpc2.ChangeMessageTextResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// changeMessageThread moves a message from an old thread to a new thread.
func (srv *Server) changeMessageThread(p *rpc2.ChangeMessageThreadParams) (result *rpc2.ChangeMessageThreadResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.MessageId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_MessageIdIsNotSet, RpcErrorMsg_MessageIdIsNotSet, nil)
	}

	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	var initialMessage derived2.IMessage
	initialMessage, re = srv.changeMessageThreadH(p.MessageId, p.ThreadId, userRoles)
	if re != nil {
		return nil, re
	}

	seData := sed.NewSystemEventDataWithValue(
		set.NewSystemEventTypeWithValue(ev.NewEnumValue(set.SystemEventType_MessageParentChange)),
		initialMessage.GetThreadIdPtr(),
		&p.MessageId,
		userRoles.User.GetUserParameters().GetIdPtr(),
		initialMessage.GetEventData().GetCreatorUserIdPtr(),
	)

	se, err := cm.NewSystemEventWithData(seData)
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_SystemEvent, c.RpcErrorMsg_SystemEvent, nil)
	}

	re = srv.reportSystemEvent(se)
	if re != nil {
		return nil, re
	}

	result = &rpc2.ChangeMessageThreadResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// getMessage reads a message.
func (srv *Server) getMessage(p *rpc2.GetMessageParams) (result *rpc2.GetMessageResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.MessageId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_MessageIdIsNotSet, RpcErrorMsg_MessageIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read the message.
	var message derived2.IMessage
	var err error
	message, err = srv.dbo.GetMessageById(p.MessageId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if message == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_MessageIsNotFound, RpcErrorMsg_MessageIsNotFound, nil)
	}

	result = &rpc2.GetMessageResult{
		Message: message,
	}

	return result, nil
}

// getLatestMessageOfThread reads the latest message of a thread.
func (srv *Server) getLatestMessageOfThread(p *rpc2.GetLatestMessageOfThreadParams) (result *rpc2.GetLatestMessageOfThreadResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	result = &rpc2.GetLatestMessageOfThreadResult{}

	result.Message, re = srv.getLatestMessageOfThreadH(p.ThreadId)
	if re != nil {
		return nil, re
	}

	return result, nil
}

// deleteMessage removes a message.
func (srv *Server) deleteMessage(p *rpc2.DeleteMessageParams) (result *rpc2.DeleteMessageResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.MessageId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_MessageIdIsNotSet, RpcErrorMsg_MessageIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	var initialMessage derived2.IMessage
	initialMessage, re = srv.deleteMessageH(p.MessageId)
	if re != nil {
		return nil, re
	}

	seData := sed.NewSystemEventDataWithValue(
		set.NewSystemEventTypeWithValue(ev.NewEnumValue(set.SystemEventType_ThreadMessageDeletion)),
		initialMessage.GetThreadIdPtr(),
		&p.MessageId,
		userRoles.User.GetUserParameters().GetIdPtr(),
		nil,
	)

	se, err := cm.NewSystemEventWithData(seData)
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_SystemEvent, c.RpcErrorMsg_SystemEvent, nil)
	}

	re = srv.reportSystemEvent(se)
	if re != nil {
		return nil, re
	}

	seData = sed.NewSystemEventDataWithValue(
		set.NewSystemEventTypeWithValue(ev.NewEnumValue(set.SystemEventType_MessageDeletion)),
		initialMessage.GetThreadIdPtr(),
		&p.MessageId,
		userRoles.User.GetUserParameters().GetIdPtr(),
		initialMessage.GetEventData().GetCreatorUserIdPtr(),
	)

	se, err = cm.NewSystemEventWithData(seData)
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_SystemEvent, c.RpcErrorMsg_SystemEvent, nil)
	}

	re = srv.reportSystemEvent(se)
	if re != nil {
		return nil, re
	}

	result = &rpc2.DeleteMessageResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// Composite objects.

// listThreadAndMessages reads a thread and all its messages.
func (srv *Server) listThreadAndMessages(p *rpc2.ListThreadAndMessagesParams) (result *rpc2.ListThreadAndMessagesResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read the thread.
	var thread derived2.IThread
	var err error
	thread, err = srv.dbo.GetThreadById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if thread == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	// Read messages.
	var allMessageIds = thread.GetMessages()

	var allMessages []derived2.IMessage
	allMessages, err = srv.dbo.ReadMessagesById(allMessageIds)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	taM := tam.NewThreadAndMessages()
	taM.SetThread(thread)
	taM.SetMessages(allMessages)

	result = &rpc2.ListThreadAndMessagesResult{
		ThreadAndMessages: taM,
	}

	return result, nil
}

// listThreadAndMessagesOnPage reads a thread and its messages on a selected page.
func (srv *Server) listThreadAndMessagesOnPage(p *rpc2.ListThreadAndMessagesOnPageParams) (result *rpc2.ListThreadAndMessagesOnPageResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	if p.Page == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_PageIsNotSet, RpcErrorMsg_PageIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read the thread.
	var thread derived2.IThread
	var err error
	thread, err = srv.dbo.GetThreadById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if thread == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	// Read messages.
	var allMessageIds = thread.GetMessages()
	var messageIdsOnPage = allMessageIds.OnPage(p.Page, srv.settings.SystemSettings.PageSize)

	var messagesOnPage []derived2.IMessage
	messagesOnPage, err = srv.dbo.ReadMessagesById(messageIdsOnPage)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	taM := tam.NewThreadAndMessages()
	taM.SetThread(thread)
	taM.SetMessages(messagesOnPage)
	taM.SetPageData(&rpc3.PageData{
		PageNumber:  p.Page,
		TotalPages:  base2.CalculateTotalPages(allMessageIds.Size(), srv.settings.SystemSettings.PageSize),
		PageSize:    srv.settings.SystemSettings.PageSize,
		ItemsOnPage: messageIdsOnPage.Size(),
		TotalItems:  allMessageIds.Size(),
	})
	thread.SetMessages(messageIdsOnPage)

	result = &rpc2.ListThreadAndMessagesOnPageResult{
		ThreadAndMessagesOnPage: taM,
	}

	return result, nil
}

// listForumAndThreads reads a forum and all its threads.
func (srv *Server) listForumAndThreads(p *rpc2.ListForumAndThreadsParams) (result *rpc2.ListForumAndThreadsResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ForumId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIdIsNotSet, RpcErrorMsg_ForumIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read the forum.
	var forum derived2.IForum
	var err error
	forum, err = srv.dbo.GetForumById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if forum == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	// Read threads.
	var allThreadIds = forum.GetThreads()

	var allThreads []derived2.IThread
	allThreads, err = srv.dbo.ReadThreadsById(allThreadIds)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	faT := fat.NewForumAndThreads()
	faT.SetForum(forum)
	faT.SetThreads(allThreads)

	result = &rpc2.ListForumAndThreadsResult{
		ForumAndThreads: faT,
	}

	return result, nil
}

// listForumAndThreadsOnPage reads a forum and its threads on a selected page.
func (srv *Server) listForumAndThreadsOnPage(p *rpc2.ListForumAndThreadsOnPageParams) (result *rpc2.ListForumAndThreadsOnPageResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ForumId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIdIsNotSet, RpcErrorMsg_ForumIdIsNotSet, nil)
	}

	if p.Page == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_PageIsNotSet, RpcErrorMsg_PageIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read the forum.
	var forum derived2.IForum
	var err error
	forum, err = srv.dbo.GetForumById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if forum == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	// Read threads.
	var allThreadIds = forum.GetThreads()
	var threadIdsOnPage = allThreadIds.OnPage(p.Page, srv.settings.SystemSettings.PageSize)

	var threadsOnPage []derived2.IThread
	threadsOnPage, err = srv.dbo.ReadThreadsById(threadIdsOnPage)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	faT := fat.NewForumAndThreads()
	faT.SetForum(forum)
	faT.SetThreads(threadsOnPage)
	faT.SetPageData(&rpc3.PageData{
		PageNumber:  p.Page,
		TotalPages:  base2.CalculateTotalPages(allThreadIds.Size(), srv.settings.SystemSettings.PageSize),
		PageSize:    srv.settings.SystemSettings.PageSize,
		ItemsOnPage: threadIdsOnPage.Size(),
		TotalItems:  allThreadIds.Size(),
	})
	forum.SetThreads(threadIdsOnPage)

	result = &rpc2.ListForumAndThreadsOnPageResult{
		ForumAndThreadsOnPage: faT,
	}

	return result, nil
}

// listSectionsAndForums reads all sections and forums.
func (srv *Server) listSectionsAndForums(p *rpc2.ListSectionsAndForumsParams) (result *rpc2.ListSectionsAndForumsResult, re *jrm1.RpcError) {
	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read all the sections.
	var sections []derived2.ISection
	var err error
	sections, err = srv.dbo.ReadSections()
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Read all the forums.
	var forums []derived2.IForum
	forums, err = srv.dbo.ReadForums()
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.ListSectionsAndForumsResult{
		SectionsAndForums: &mm.SectionsAndForums{
			Sections: sections,
			Forums:   forums,
		},
	}

	return result, nil
}

// Other.

func (srv *Server) getDKey(p *rpc2.GetDKeyParams) (result *rpc2.GetDKeyResult, re *jrm1.RpcError) {
	re = srv.mustBeNoAuth(p.Auth)
	if re != nil {
		return nil, re
	}

	result = &rpc2.GetDKeyResult{
		DKey: base2.Text(srv.dKeyI.GetString()),
	}

	return result, nil
}

func (srv *Server) showDiagnosticData() (result *rpc2.ShowDiagnosticDataResult, re *jrm1.RpcError) {
	trc, src := srv.js.GetRequestsCount()

	result = &rpc2.ShowDiagnosticDataResult{
		RequestsCount: rpc3.RequestsCount{
			TotalRequestsCount:      base2.Text(trc),
			SuccessfulRequestsCount: base2.Text(src),
		},
	}

	return result, nil
}

func (srv *Server) test(p *rpc2.TestParams) (result *rpc2.TestResult, re *jrm1.RpcError) {
	result = &rpc2.TestResult{}

	var wg = new(sync.WaitGroup)
	var errChan = make(chan error, p.N)

	for i := uint(1); i <= p.N; i++ {
		wg.Add(1)
		go srv.doTestA(wg, errChan)
	}
	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			srv.logError(err)
			return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_TestError, fmt.Sprintf(RpcErrorMsgF_TestError, err.Error()), nil)
		}
	}

	return result, nil
}
