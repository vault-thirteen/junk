package c

import (
	"fmt"
	"time"

	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/dbc"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
)

// Functions for scheduler component.

func (c *Controller) RemoveOutdatedRegistrationRequests() (err error) {
	rrTtl := c.far.systemSettings.GetParameterAsInt(ccp.RegistrationRequestTtl)
	edgeTime := time.Now().Add(-time.Duration(rrTtl) * time.Second)
	dbC := dbc.NewDbController(c.GetDb())
	var isDebugMode = c.far.systemSettings.GetParameterAsBool(ccp.IsDebugMode)

	var rr *cm.RegistrationRequest
	for {
		rr, err = c.getNextOutdatedRegistrationRequest(edgeTime)
		if err != nil {
			return err
		}

		if rr == nil {
			break
		}

		if isDebugMode {
			fmt.Println("removing outdated registration request. RequestId:", rr.RequestId)
		}
		err = dbC.DeleteRegistrationRequestNRFA(rr)
		if err != nil {
			return err
		}
	}

	return nil
}
func (c *Controller) getNextOutdatedRegistrationRequest(edgeTime time.Time) (rr *cm.RegistrationRequest, err error) {
	dbC := dbc.NewDbController(c.GetDb())

	var rrs []cm.RegistrationRequest
	rrs, err = dbC.GetFirstOutdatedRegistrationRequest(edgeTime)
	if err != nil {
		return nil, err
	}

	if len(rrs) != 1 {
		return nil, nil
	}

	return &rrs[0], nil
}

func (c *Controller) RemoveOutdatedLogInRequests() (err error) {
	lirTtl := c.far.systemSettings.GetParameterAsInt(ccp.LogInRequestTtl)
	edgeTime := time.Now().Add(-time.Duration(lirTtl) * time.Second)
	dbC := dbc.NewDbController(c.GetDb())
	var isDebugMode = c.far.systemSettings.GetParameterAsBool(ccp.IsDebugMode)

	var lir *cm.LogInRequest
	for {
		lir, err = c.getNextOutdatedLogInRequest(edgeTime)
		if err != nil {
			return err
		}

		if lir == nil {
			break
		}

		if isDebugMode {
			fmt.Println("removing outdated log-in request. RequestId:", lir.RequestId)
		}
		err = dbC.DeleteOldLogInRequest(lir)
		if err != nil {
			return err
		}
	}

	return nil
}
func (c *Controller) getNextOutdatedLogInRequest(edgeTime time.Time) (lir *cm.LogInRequest, err error) {
	dbC := dbc.NewDbController(c.GetDb())

	var lirs []cm.LogInRequest
	lirs, err = dbC.GetFirstOutdatedLogInRequest(edgeTime)
	if err != nil {
		return nil, err
	}

	if len(lirs) != 1 {
		return nil, nil
	}

	return &lirs[0], nil
}

func (c *Controller) RemoveOutdatedLogOutRequests() (err error) {
	lorTtl := c.far.systemSettings.GetParameterAsInt(ccp.LogOutRequestTtl)
	edgeTime := time.Now().Add(-time.Duration(lorTtl) * time.Second)
	dbC := dbc.NewDbController(c.GetDb())
	var isDebugMode = c.far.systemSettings.GetParameterAsBool(ccp.IsDebugMode)

	var lor *cm.LogOutRequest
	for {
		lor, err = c.getNextOutdatedLogOutRequest(edgeTime)
		if err != nil {
			return err
		}

		if lor == nil {
			break
		}

		if isDebugMode {
			fmt.Println("removing outdated log-out request. RequestId:", lor.RequestId)
		}
		err = dbC.DeleteOldLogOutRequest(lor)
		if err != nil {
			return err
		}
	}

	return nil
}
func (c *Controller) getNextOutdatedLogOutRequest(edgeTime time.Time) (lor *cm.LogOutRequest, err error) {
	dbC := dbc.NewDbController(c.GetDb())

	var lors []cm.LogOutRequest
	lors, err = dbC.GetFirstOutdatedLogOutRequest(edgeTime)
	if err != nil {
		return nil, err
	}

	if len(lors) != 1 {
		return nil, nil
	}

	return &lors[0], nil
}

func (c *Controller) RemoveOutdatedEmailChangeRequests() (err error) {
	ecrTtl := c.far.systemSettings.GetParameterAsInt(ccp.EmailChangeRequestTtl)
	edgeTime := time.Now().Add(-time.Duration(ecrTtl) * time.Second)
	dbC := dbc.NewDbController(c.GetDb())
	var isDebugMode = c.far.systemSettings.GetParameterAsBool(ccp.IsDebugMode)

	var ecr *cm.EmailChangeRequest
	for {
		ecr, err = c.getNextOutdatedEmailChangeRequest(edgeTime)
		if err != nil {
			return err
		}

		if ecr == nil {
			break
		}

		if isDebugMode {
			fmt.Println("removing outdated e-mail change request. RequestId:", ecr.RequestId)
		}
		err = dbC.DeleteOldEmailChangeRequest(ecr)
		if err != nil {
			return err
		}
	}

	return nil
}
func (c *Controller) getNextOutdatedEmailChangeRequest(edgeTime time.Time) (ecr *cm.EmailChangeRequest, err error) {
	dbC := dbc.NewDbController(c.GetDb())

	var ecrs []cm.EmailChangeRequest
	ecrs, err = dbC.GetFirstOutdatedEmailChangeRequest(edgeTime)
	if err != nil {
		return nil, err
	}

	if len(ecrs) != 1 {
		return nil, nil
	}

	return &ecrs[0], nil
}

func (c *Controller) RemoveOutdatedPasswordChangeRequests() (err error) {
	pcrTtl := c.far.systemSettings.GetParameterAsInt(ccp.PasswordChangeRequestTtl)
	edgeTime := time.Now().Add(-time.Duration(pcrTtl) * time.Second)
	dbC := dbc.NewDbController(c.GetDb())
	var isDebugMode = c.far.systemSettings.GetParameterAsBool(ccp.IsDebugMode)

	var pcr *cm.PasswordChangeRequest
	for {
		pcr, err = c.getNextOutdatedPasswordChangeRequest(edgeTime)
		if err != nil {
			return err
		}

		if pcr == nil {
			break
		}

		if isDebugMode {
			fmt.Println("removing outdated password change request. RequestId:", pcr.RequestId)
		}
		err = dbC.DeleteOldPasswordChangeRequest(pcr)
		if err != nil {
			return err
		}
	}

	return nil
}
func (c *Controller) getNextOutdatedPasswordChangeRequest(edgeTime time.Time) (pcr *cm.PasswordChangeRequest, err error) {
	dbC := dbc.NewDbController(c.GetDb())

	var pcrs []cm.PasswordChangeRequest
	pcrs, err = dbC.GetFirstOutdatedPasswordChangeRequest(edgeTime)
	if err != nil {
		return nil, err
	}

	if len(pcrs) != 1 {
		return nil, nil
	}

	return &pcrs[0], nil
}

func (c *Controller) RemoveOutdatedSessions() (err error) {
	sessionTtl := c.far.systemSettings.GetParameterAsInt(ccp.SessionMaxDuration)
	edgeTime := time.Now().Add(-time.Duration(sessionTtl) * time.Second)
	dbC := dbc.NewDbController(c.GetDb())
	var isDebugMode = c.far.systemSettings.GetParameterAsBool(ccp.IsDebugMode)

	var s *cm.Session
	for {
		s, err = c.getNextOutdatedSession(edgeTime)
		if err != nil {
			return err
		}

		if s == nil {
			break
		}

		if isDebugMode {
			fmt.Println("removing outdated session. Id:", s.Id)
		}

		// Delete session.
		err = dbC.DeleteOldSession(s)
		if err != nil {
			return err
		}

		// Journaling.
		logEvent := cm.NewLogEvent(cm.LogEvent_Type_LogOutByTimeout, s.UserId, nil, nil)

		err = dbC.CreateLogEvent(logEvent)
		if err != nil {
			return err
		}
	}

	return nil
}
func (c *Controller) getNextOutdatedSession(edgeTime time.Time) (s *cm.Session, err error) {
	dbC := dbc.NewDbController(c.GetDb())

	var ss []cm.Session
	ss, err = dbC.GetFirstOutdatedSession(edgeTime)
	if err != nil {
		return nil, err
	}

	if len(ss) != 1 {
		return nil, nil
	}

	return &ss[0], nil
}
