package dbo

// Due to the large number of methods, they are sorted alphabetically.

import (
	"database/sql"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived2"
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/Resource"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/SystemEvent"
	dbo2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/dbo"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	cms "github.com/vault-thirteen/SimpleBB/pkg/common/models/sql"
	"net"

	nm "github.com/vault-thirteen/SimpleBB/pkg/NM/models"
	ae "github.com/vault-thirteen/auxie/errors"
)

func (dbo *DatabaseObject) AddResource(r derived2.IResource) (lastInsertedId base2.Id, err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_AddResource).Exec(r.GetType(), r.GetText(), r.GetNumber())
	if err != nil {
		return dbo2.LastInsertedIdOnError, err
	}

	return dbo2.CheckRowsAffectedAndGetLastInsertedId(result, 1)
}

func (dbo *DatabaseObject) CountAllNotificationsByUserId(userId base2.Id) (n base2.Count, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountAllNotificationsByUserId).QueryRow(userId)

	n, err = cms.NewNonNullValueFromScannableSource[base2.Count](row)
	if err != nil {
		return dbo2.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountAllResources() (n base2.Count, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountAllResources).QueryRow()

	n, err = cms.NewNonNullValueFromScannableSource[base2.Count](row)
	if err != nil {
		return dbo2.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountUnreadNotificationsByUserId(userId base2.Id) (n base2.Count, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountUnreadNotificationsByUserId).QueryRow(userId)

	n, err = cms.NewNonNullValueFromScannableSource[base2.Count](row)
	if err != nil {
		return dbo2.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) DeleteNotificationById(notificationId base2.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_DeleteNotificationById).Exec(notificationId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) DeleteResourceById(resourceId base2.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_DeleteResourceById).Exec(resourceId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) GetAllNotificationsByUserId(userId base2.Id) (notifications []nm.Notification, err error) {
	var rows *sql.Rows
	rows, err = dbo.DatabaseObject.PreparedStatement(DbPsid_GetAllNotificationsByUserId).Query(userId)
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := rows.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	return nm.NewNotificationArrayFromRows(rows)
}

func (dbo *DatabaseObject) GetNotificationById(notificationId base2.Id) (notification *nm.Notification, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetNotificationById).QueryRow(notificationId)

	notification, err = nm.NewNotificationFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return notification, nil
}

func (dbo *DatabaseObject) GetNotificationsByUserIdOnPage(userId base2.Id, pageNumber base2.Count, pageSize base2.Count) (notifications []nm.Notification, err error) {
	var rows *sql.Rows
	rows, err = dbo.PreparedStatement(DbPsid_GetNotificationsByUserIdOnPage).Query(userId, pageSize, (pageNumber-1)*pageSize)
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := rows.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	return nm.NewNotificationArrayFromRows(rows)
}

func (dbo *DatabaseObject) GetResourceById(resourceId base2.Id) (r derived2.IResource, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetResourceById).QueryRow(resourceId)

	r, err = res.NewResourceFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (dbo *DatabaseObject) GetResourceIdsOnPage(pageNumber base2.Count, pageSize base2.Count) (resourceIds []base2.Id, err error) {
	var rows *sql.Rows
	rows, err = dbo.PreparedStatement(DbPsid_ListAllResourceIdsOnPage).Query(pageSize, (pageNumber-1)*pageSize)
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := rows.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	return cms.NewArrayFromScannableSource[base2.Id](rows)
}

func (dbo *DatabaseObject) GetSystemEventById(systemEventId base2.Id) (se derived2.ISystemEvent, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetSystemEventById).QueryRow(systemEventId)

	se, err = cm.NewSystemEventFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return se, nil
}

func (dbo *DatabaseObject) GetUnreadNotifications(userId base2.Id) (notifications []nm.Notification, err error) {
	var rows *sql.Rows
	rows, err = dbo.DatabaseObject.PreparedStatement(DbPsid_GetUnreadNotificationsByUserId).Query(userId)
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := rows.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	return nm.NewNotificationArrayFromRows(rows)
}

func (dbo *DatabaseObject) InsertNewNotification(userId base2.Id, text base2.Text) (lastInsertedId base2.Id, err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_InsertNewNotification).Exec(userId, text)
	if err != nil {
		return dbo2.LastInsertedIdOnError, err
	}

	return dbo2.CheckRowsAffectedAndGetLastInsertedId(result, 1)
}

func (dbo *DatabaseObject) MarkNotificationAsRead(notificationId base2.Id, userId base2.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_MarkNotificationAsRead).Exec(notificationId, userId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SaveIncident(module derived1.IModule, incidentType derived1.IIncidentType, email simple.Email, userIPAB net.IP) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SaveIncident).Exec(module, incidentType, email, userIPAB)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SaveIncidentWithoutUserIPA(module derived1.IModule, incidentType derived1.IIncidentType, email simple.Email) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SaveIncidentWithoutUserIPA).Exec(module, incidentType, email)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SaveSystemEvent(se derived2.ISystemEvent) (err error) {
	var result sql.Result
	sed := se.GetSystemEventData()

	result, err = dbo.PreparedStatement(DbPsid_SaveSystemEvent).Exec(sed.GetType(), sed.GetThreadId(), sed.GetMessageId(), sed.GetUserId())
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}
