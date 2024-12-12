package s

import (
	"database/sql"
	"errors"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base2"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived2"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/UidList"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/EventData"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/SectionChildType"
	cms "github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
)

type section struct {
	// Identifier of a section.
	Id cmb.Id `json:"id"`

	// Identifier of a parent section containing this section.
	// Null means that this section is a root section.
	// Only a single root section can exist.
	Parent *cmb.Id `json:"parent"`

	// Type of child elements: either sub-sections or forums.
	ChildType derived1.ISectionChildType `json:"childType"`

	// List of IDs of child elements (either sub-sections or forums).
	// Null means that this section has no derived1.
	Children *ul.UidList `json:"children"`

	// Name of this section.
	Name cms.Name `json:"name"`

	// Section meta-data.
	base2.IEventData
}

func NewSection() (sec derived2.ISection) {
	return &section{
		ChildType:  sct.NewSectionChildType(),
		IEventData: ed.NewEventData(),
	}
}

func NewSectionFromScannableSource(src base.IScannable) (sec derived2.ISection, err error) {
	sec = NewSection()
	var x = ul.New()
	var eventData = ed.NewEventData()

	err = src.Scan(
		sec.GetIdPtr(),
		sec.GetParentPtr(),
		sec.GetChildTypePtr(),
		x, //&sec.Children,
		sec.GetNamePtr(),
		eventData.GetCreatorUserIdPtr(),
		eventData.GetCreatorTimePtr(),
		eventData.GetEditorUserIdPtr(),
		eventData.GetEditorTimePtr(),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	sec.SetEventData(eventData)
	sec.SetChildren(x)
	return sec, nil
}

func NewSectionArrayFromRows(rows base.IScannableSequence) (sections []derived2.ISection, err error) {
	sections = []derived2.ISection{}
	var s derived2.ISection

	for rows.Next() {
		s, err = NewSectionFromScannableSource(rows)
		if err != nil {
			return nil, err
		}

		sections = append(sections, s)
	}

	return sections, nil
}

// Emulated class members.
func (s *section) GetIdPtr() (id *cmb.Id)                                  { return &s.Id }
func (s *section) GetId() (id cmb.Id)                                      { return s.Id }
func (s *section) GetParentPtr() (parent **cmb.Id)                         { return &s.Parent }
func (s *section) GetParent() (parent *cmb.Id)                             { return s.Parent }
func (s *section) GetChildTypePtr() (childType derived1.ISectionChildType) { return s.ChildType }
func (s *section) GetChildType() (childType derived1.ISectionChildType)    { return s.ChildType }
func (s *section) GetChildrenPtr() (children **ul.UidList)                 { return &s.Children }
func (s *section) GetChildren() (children *ul.UidList)                     { return s.Children }
func (s *section) GetNamePtr() (name *cms.Name)                            { return &s.Name }
func (s *section) GetEventDataPtr() base2.IEventData                       { return s.IEventData }
func (s *section) SetEventData(ed base2.IEventData) {
	s.IEventData = ed
}
func (s *section) SetChildren(children *ul.UidList) { s.Children = children }
