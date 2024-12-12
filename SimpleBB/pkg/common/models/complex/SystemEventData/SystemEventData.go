package sed

import (
	"fmt"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived2"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	set "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/SystemEventType"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
)

const (
	ErrSystemEventType       = "system event type error"
	ErrSystemEventParameters = "system event parameters error"
)

// systemEventData is a set of parameters of a system event.
type systemEventData struct {
	// Event type. This field is required. All other fields are optional.
	// The number and values of other fields depend on the event type.
	Type derived1.ISystemEventType `json:"type"`

	// ID of the thread mentioned in the event.
	// E.g., if a thread was renamed, it's ID is this field.
	ThreadId *cmb.Id `json:"threadId"`

	// ID of the message mentioned in the event.
	// E.g., if a message was added into a thread, it's ID is this field.
	MessageId *cmb.Id `json:"messageId"`

	// ID of the user mentioned in the event.
	// E.g., if a user has changed name of a thread, it's ID is this field.
	UserId *cmb.Id `json:"userId"`

	// Auxiliary fields.

	// ID of a user who initially created the object.
	// Some events deal with two UserId's, where one is the ID of the original
	// creator of the object and another is the user who modified the object.
	// In such cases, the original creator of the object is stored in the
	// 'Creator' field and the user who modified the object is stored in the
	// 'UserId' field. In cases where only a single user is important, this
	// field is not used, as this only user is stored in the 'UserId' field.
	Creator *cmb.Id `json:"creator,omitempty"`
}

func NewSystemEventData() derived2.ISystemEventData {
	return &systemEventData{
		Type: set.NewSystemEventType(),
	}
}

func NewSystemEventDataWithValue(t derived1.ISystemEventType, threadId *cmb.Id, messageId *cmb.Id, userId *cmb.Id, creator *cmb.Id) derived2.ISystemEventData {
	return &systemEventData{
		Type:      t,
		ThreadId:  threadId,
		MessageId: messageId,
		UserId:    userId,
		Creator:   creator,
	}
}

func (sed *systemEventData) isThreadIdSet() (ok bool) {
	return (sed.ThreadId != nil) && (*sed.ThreadId > 0)
}

func (sed *systemEventData) isMessageIdSet() (ok bool) {
	return (sed.MessageId != nil) && (*sed.MessageId > 0)
}

func (sed *systemEventData) isUserIdSet() (ok bool) {
	return (sed.UserId != nil) && (*sed.UserId > 0)
}

func (sed *systemEventData) isCreatorSet() (ok bool) {
	return (sed.Creator != nil) && (*sed.Creator > 0)
}

func (sed *systemEventData) CheckParameters() (ok bool, err error) {
	// Default requirements.
	var req = simple.SystemEventRequirements{
		IsThreadIdRequired: true,
		IsUserIdRequired:   true,
	}

	switch sed.Type.GetValue().RawValue() {
	case set.SystemEventType_ThreadParentChange:
		// Default requirements are used (TU).

	case set.SystemEventType_ThreadNameChange:
		// Default requirements are used (TU).

	case set.SystemEventType_ThreadDeletion:
		// Default requirements are used (TU).

	case set.SystemEventType_ThreadNewMessage:
		// TMU.
		req.IsMessageIdRequired = true

	case set.SystemEventType_ThreadMessageEdit:
		// TMU.
		req.IsMessageIdRequired = true

	case set.SystemEventType_ThreadMessageDeletion:
		// TMU.
		req.IsMessageIdRequired = true

	case set.SystemEventType_MessageTextEdit:
		// TMUC.
		req.IsMessageIdRequired = true
		req.IsCreatorRequired = true

	case set.SystemEventType_MessageParentChange:
		// TMUC.
		req.IsMessageIdRequired = true
		req.IsCreatorRequired = true

	case set.SystemEventType_MessageDeletion:
		// TMUC.
		req.IsMessageIdRequired = true
		req.IsCreatorRequired = true

	default:
		return false, fmt.Errorf(ErrSystemEventType)
	}

	// Check the required parameters.
	if req.IsThreadIdRequired {
		if !sed.isThreadIdSet() {
			return false, fmt.Errorf(ErrSystemEventParameters)
		}
	}

	if req.IsMessageIdRequired {
		if !sed.isMessageIdSet() {
			return false, fmt.Errorf(ErrSystemEventParameters)
		}
	}

	if req.IsUserIdRequired {
		if !sed.isUserIdSet() {
			return false, fmt.Errorf(ErrSystemEventParameters)
		}
	}

	if req.IsCreatorRequired {
		if !sed.isCreatorSet() {
			return false, fmt.Errorf(ErrSystemEventParameters)
		}
	}

	return true, nil
}

// Emulated class members.
func (sed *systemEventData) GetType() (t derived1.ISystemEventType) { return sed.Type }
func (sed *systemEventData) GetThreadIdPtr() (threadId **cmb.Id)    { return &sed.ThreadId }
func (sed *systemEventData) GetThreadId() (threadId *cmb.Id)        { return sed.ThreadId }
func (sed *systemEventData) GetMessageIdPtr() (messageId **cmb.Id)  { return &sed.MessageId }
func (sed *systemEventData) GetMessageId() (messageId *cmb.Id)      { return sed.MessageId }
func (sed *systemEventData) GetUserIdPtr() (userId **cmb.Id)        { return &sed.UserId }
func (sed *systemEventData) GetUserId() (userId *cmb.Id)            { return sed.UserId }
func (sed *systemEventData) GetCreator() (creator *cmb.Id)          { return sed.Creator }
