package server

import (
	"errors"
	"fmt"
	"github.com/vault-thirteen/SimpleBB/pkg/MM/dbo"
	mm "github.com/vault-thirteen/SimpleBB/pkg/MM/models"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived2"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/SectionChildType"
	c "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
)

const (
	Err_TooManyRootSections = "too many root sections"
	ErrF_SectionIsNotFound  = "section is not found, ID=%v"
	ErrF_SectionIsDamaged   = "section is damaged, ID=%v"
	ErrF_ForumIsNotFound    = "forum is not found, ID=%v"
	ErrF_ForumIsDamaged     = "forum is damaged, ID=%v"
	ErrF_ThreadIsNotFound   = "thread is not found, ID=%v"
	ErrF_ThreadIsDamaged    = "thread is damaged, ID=%v"
	ErrF_MessageIsNotFound  = "message is not found, ID=%v"
	ErrF_MessageIsDamaged   = "message is damaged, ID=%v"
)

// checkDatabaseConsistency checks consistency of sections, forums, threads and
// messages. This function is used in the scheduler and is also run once during
// the server's start.
func (srv *Server) checkDatabaseConsistency() (err error) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	fmt.Print(c.MsgDatabaseConsistencyCheck)

	// Sections.
	var sections []derived2.ISection
	sections, err = srv.dbo.ReadSections()
	if err != nil {
		return err
	}

	sectionsMap := make(map[cmb.Id]derived2.ISection)
	for _, section := range sections {
		sectionsMap[section.GetId()] = section
	}

	err = checkSections(sections, sectionsMap)
	if err != nil {
		return err
	}

	// Forums.
	var forums []derived2.IForum
	forums, err = srv.dbo.ReadForums()
	if err != nil {
		return err
	}

	forumsMap := make(map[cmb.Id]derived2.IForum)
	for _, forum := range forums {
		forumsMap[forum.GetId()] = forum
	}

	err = checkForums(sections, sectionsMap, forums, forumsMap)
	if err != nil {
		return err
	}

	// Threads.
	var threads []mm.ThreadLink
	threads, err = srv.dbo.ReadThreadLinks()
	if err != nil {
		return err
	}

	threadsMap := make(map[cmb.Id]mm.ThreadLink)
	for _, thread := range threads {
		threadsMap[thread.Id] = thread
	}

	err = checkThreads(forums, forumsMap, threads, threadsMap)
	if err != nil {
		return err
	}

	// Messages.
	err = checkMessages(srv.dbo, threads)
	if err != nil {
		return err
	}

	fmt.Println(c.MsgOK)

	return nil
}

func checkSections(sections []derived2.ISection, sectionsMap map[cmb.Id]derived2.ISection) (err error) {
	// Step I. Downward check (parent to child).
	var childSection derived2.ISection
	var childIds []cmb.Id
	var ok bool
	for _, section := range sections {
		if section.GetChildType().AsInt() != sct.SectionChildType_Section {
			continue
		}

		if section.GetChildren().Size() == 0 {
			continue
		}

		childIds = section.GetChildren().AsArray()

		for _, childId := range childIds {
			childSection, ok = sectionsMap[childId]
			if !ok {
				return fmt.Errorf(ErrF_SectionIsNotFound, childId)
			}

			if childSection.GetParent() == nil {
				return fmt.Errorf(ErrF_SectionIsDamaged, childId)
			}

			if *childSection.GetParent() != section.GetId() {
				return fmt.Errorf(ErrF_SectionIsDamaged, childId)
			}
		}
	}

	// Step II. Root section.
	var rootSectionsCount = 0
	for _, section := range sections {
		if section.GetParent() != nil {
			continue
		}
		rootSectionsCount++
	}
	if rootSectionsCount > 1 {
		return errors.New(Err_TooManyRootSections)
	}

	// Step III. Upward check (child to parent).
	var parentId cmb.Id
	var parentSection derived2.ISection
	for _, section := range sections {
		if section.GetParent() == nil {
			continue
		}

		parentId = *section.GetParent()

		parentSection, ok = sectionsMap[parentId]
		if !ok {
			return fmt.Errorf(ErrF_SectionIsNotFound, parentId)
		}

		if parentSection.GetChildType().AsInt() != sct.SectionChildType_Section {
			return fmt.Errorf(ErrF_SectionIsDamaged, parentId)
		}

		if parentSection.GetChildren().Size() == 0 {
			return fmt.Errorf(ErrF_SectionIsDamaged, parentId)
		}

		if !parentSection.GetChildren().HasItem(section.GetId()) {
			return fmt.Errorf(ErrF_SectionIsDamaged, parentId)
		}
	}

	return nil
}

func checkForums(sections []derived2.ISection, sectionsMap map[cmb.Id]derived2.ISection, forums []derived2.IForum, forumsMap map[cmb.Id]derived2.IForum) (err error) {
	// Step I. Downward check (parent to child).
	var childIds []cmb.Id
	var ok bool
	for _, section := range sections {
		if section.GetChildType().AsInt() != sct.SectionChildType_Forum {
			continue
		}

		if section.GetChildren().Size() == 0 {
			continue
		}

		childIds = section.GetChildren().AsArray()

		var forum derived2.IForum
		for _, childId := range childIds {
			forum, ok = forumsMap[childId]
			if !ok {
				return fmt.Errorf(ErrF_ForumIsNotFound, childId)
			}

			if forum.GetSectionId() != section.GetId() {
				return fmt.Errorf(ErrF_ForumIsDamaged, childId)
			}
		}
	}

	// Step II. Upward check (child to parent).
	var parentId cmb.Id
	var parentSection derived2.ISection
	for _, forum := range forums {
		parentId = forum.GetSectionId()

		parentSection, ok = sectionsMap[parentId]
		if !ok {
			return fmt.Errorf(ErrF_SectionIsNotFound, parentId)
		}

		if parentSection.GetChildType().AsInt() != sct.SectionChildType_Forum {
			return fmt.Errorf(ErrF_SectionIsDamaged, parentId)
		}

		if parentSection.GetChildren().Size() == 0 {
			return fmt.Errorf(ErrF_SectionIsDamaged, parentId)
		}

		if !parentSection.GetChildren().HasItem(forum.GetId()) {
			return fmt.Errorf(ErrF_SectionIsDamaged, parentId)
		}
	}

	return nil
}

func checkThreads(forums []derived2.IForum, forumsMap map[cmb.Id]derived2.IForum, threads []mm.ThreadLink, threadsMap map[cmb.Id]mm.ThreadLink) (err error) {
	// Step I. Downward check (parent to child).
	var childIds []cmb.Id
	var ok bool
	for _, forum := range forums {
		if forum.GetThreads().Size() == 0 {
			continue
		}

		childIds = forum.GetThreads().AsArray()

		var thread mm.ThreadLink
		for _, childId := range childIds {
			thread, ok = threadsMap[childId]
			if !ok {
				return fmt.Errorf(ErrF_ThreadIsNotFound, childId)
			}

			if thread.ForumId != forum.GetId() {
				return fmt.Errorf(ErrF_ThreadIsDamaged, childId)
			}
		}
	}

	// Step II. Upward check (child to parent).
	var parentId cmb.Id
	var parentForum derived2.IForum
	for _, thread := range threads {
		parentId = thread.ForumId

		parentForum, ok = forumsMap[parentId]
		if !ok {
			return fmt.Errorf(ErrF_ForumIsNotFound, parentId)
		}

		if parentForum.GetThreads().Size() == 0 {
			return fmt.Errorf(ErrF_ForumIsDamaged, parentId)
		}

		if !parentForum.GetThreads().HasItem(thread.Id) {
			return fmt.Errorf(ErrF_ForumIsDamaged, parentId)
		}
	}

	return nil
}

func checkMessages(dbo *dbo.DatabaseObject, threads []mm.ThreadLink) (err error) {
	// Step I. Downward check (parent to child).
	var messages []mm.MessageLink
	for _, thread := range threads {
		if thread.Messages.Size() == 0 {
			continue
		}

		messages, err = dbo.ReadMessageLinksById(thread.Messages)
		if err != nil {
			return err
		}

		messagesMap := make(map[cmb.Id]mm.MessageLink)
		for _, message := range messages {
			messagesMap[message.Id] = message
		}

		var ok bool
		var message mm.MessageLink
		for _, messageId := range thread.Messages.AsArray() {
			message, ok = messagesMap[messageId]
			if !ok {
				return fmt.Errorf(ErrF_MessageIsNotFound, messageId)
			}

			if message.ThreadId != thread.Id {
				return fmt.Errorf(ErrF_MessageIsDamaged, message.Id)
			}
		}
	}

	// Step II. Upward check (child to parent).
	// This kind of check requires huge amount of time.
	// It is not implemented.

	return nil
}
