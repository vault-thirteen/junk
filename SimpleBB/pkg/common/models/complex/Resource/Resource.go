package res

import (
	"database/sql"
	"errors"
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived2"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/EnumValue"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/ResourceType"
	"time"
)

type resource struct {
	Id     cmb.Id                 `json:"id"`
	Type   derived1.IResourceType `json:"type"`
	Text   *cmb.Text              `json:"text"`
	Number *cmb.Count             `json:"number"`
	ToC    time.Time              `json:"toc"`
}

func NewResource() (r derived2.IResource) {
	return &resource{
		Type: rt.NewResourceType(),
	}
}

func NewResourceFromValue(value any) (r derived2.IResource) {
	switch value.(type) {
	case float64:
		return NewResourceFromNumber(cmb.Count(value.(float64)))

	case int64:
		return NewResourceFromNumber(cmb.Count(value.(int64)))

	case int:
		return NewResourceFromNumber(cmb.Count(value.(int)))

	case []byte:
		return NewResourceFromText(cmb.Text(value.([]byte)))

	case string:
		return NewResourceFromText(cmb.Text(value.(string)))
	}

	return nil
}

func NewResourceFromText(t cmb.Text) (r derived2.IResource) {
	return &resource{
		Type: rt.NewResourceTypeWithValue(ev.NewEnumValue(rt.ResourceType_Text)),
		Text: &t,
	}
}

func NewResourceFromNumber(n cmb.Count) (r derived2.IResource) {
	return &resource{
		Type:   rt.NewResourceTypeWithValue(ev.NewEnumValue(rt.ResourceType_Number)),
		Number: &n,
	}
}

func NewResourceFromScannableSource(src cmi.IScannable) (r derived2.IResource, err error) {
	r = NewResource()

	err = src.Scan(
		r.GetIdPtr(),
		r.GetTypePtr(),
		r.GetTextPtr(),
		r.GetNumberPtr(),
		r.GetTocPtr(),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return r, nil
}

func (r *resource) IsValid() bool {
	if r == nil {
		return false
	}

	switch r.Type.GetValue().RawValue() {
	case rt.ResourceType_Text:
		if r.Text == nil {
			return false
		}

	case rt.ResourceType_Number:
		if r.Number == nil {
			return false
		}

	default:
		return false
	}

	return true
}

func (r *resource) AsText() cmb.Text {
	return *r.Text
}

func (r *resource) AsNumber() cmb.Count {
	return *r.Number
}

func (r *resource) GetValue() any {
	switch r.Type.GetValue().RawValue() {
	case rt.ResourceType_Text:
		return *r.Text

	case rt.ResourceType_Number:
		return *r.Number

	default:
		return nil
	}
}

// Emulated class members.
func (r *resource) GetIdPtr() (id *cmb.Id) {
	return &r.Id
}
func (r *resource) GetTypePtr() (t derived1.IResourceType) { return r.Type }
func (r *resource) GetTextPtr() (text **cmb.Text) {
	return &r.Text
}
func (r *resource) GetNumberPtr() (number **cmb.Count) {
	return &r.Number
}
func (r *resource) GetTocPtr() (toc *time.Time) {
	return &r.ToC
}
func (r *resource) GetType() (t derived1.IResourceType) { return r.Type }
func (r *resource) GetText() (text *cmb.Text) {
	return r.Text
}
func (r *resource) GetNumber() (number *cmb.Count) {
	return r.Number
}
