package im

import (
	"context"
	"errors"
	"fmt"
	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	gc "github.com/vault-thirteen/SimpleBB/pkg/GWM/client"
	gm "github.com/vault-thirteen/SimpleBB/pkg/GWM/rpc"
	"github.com/vault-thirteen/SimpleBB/pkg/NM/dbo"
	s "github.com/vault-thirteen/SimpleBB/pkg/NM/settings"
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived2"
	cc "github.com/vault-thirteen/SimpleBB/pkg/common/models/Client"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/app"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/avm"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	ev "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/EnumValue"
	inc "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/Incident"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/IncidentType"
	m "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/Module"
	server2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	"log"
	"net"
	"sync"
)

const (
	TaskChannelSize = 4
)

type IncidentManager struct {
	ssp                    *avm.SSP
	wg                     *sync.WaitGroup
	tasks                  chan derived2.IIncident
	isTableOfIncidentsUsed bool
	dbo                    *dbo.DatabaseObject
	gwmClient              *cc.Client

	// Block time in seconds for each incident type.
	blockTimePerIncidentType [it.IncidentTypeMax + 1]cmb.Count
}

func NewIncidentManager(
	isTableOfIncidentsUsed bool,
	dbo *dbo.DatabaseObject,
	gwmClient *cc.Client,
	blockTimePerIncident *s.BlockTimePerIncident,
) (im *IncidentManager) {
	im = &IncidentManager{
		ssp:                      avm.NewSSP(),
		wg:                       new(sync.WaitGroup),
		tasks:                    make(chan derived2.IIncident, TaskChannelSize),
		isTableOfIncidentsUsed:   isTableOfIncidentsUsed,
		dbo:                      dbo,
		gwmClient:                gwmClient,
		blockTimePerIncidentType: initBlockTimePerIncidentType(blockTimePerIncident),
	}

	return im
}

func initBlockTimePerIncidentType(blockTimePerIncident *s.BlockTimePerIncident) (blockTimePerIncidentType [it.IncidentTypeMax + 1]cmb.Count) {
	// The "zero"-indexed element is empty because it is not used.
	blockTimePerIncidentType[it.IncidentType_IllegalAccessAttempt] = blockTimePerIncident.IllegalAccessAttempt
	blockTimePerIncidentType[it.IncidentType_ReadingNotificationOfOtherUsers] = blockTimePerIncident.ReadingNotificationOfOtherUsers
	blockTimePerIncidentType[it.IncidentType_WrongDKey] = blockTimePerIncident.WrongDKey

	return blockTimePerIncidentType
}

// Start starts the incident manager.
func (im *IncidentManager) Start() (err error) {
	im.ssp.Lock()
	defer im.ssp.Unlock()

	err = im.ssp.BeginStart()
	if err != nil {
		return err
	}

	im.wg.Add(1)
	go im.run()

	im.ssp.CompleteStart()

	return nil
}

// run is the main work loop of the incident manager.
func (im *IncidentManager) run() {
	defer im.wg.Done()

	var err error
	var re *jrm1.RpcError
	for i := range im.tasks {
		if im.isTableOfIncidentsUsed {
			err = i.Check()
			if err != nil {
				log.Println(err)
				continue
			}

			err = im.saveIncident(i)
			im.logError(err)

			re = im.informGateway(i)
			// This is why Go language is a complete Scheiße (utter trash):
			// https://github.com/golang/go/issues/40442
			if re != nil {
				err = re.AsError()
			} else {
				err = nil
			}
			im.logError(err)
		}
	}

	log.Println(server2.MsgIncidentManagerHasStopped)
}

// Stop stops the incident manager.
func (im *IncidentManager) Stop() (err error) {
	im.ssp.Lock()
	defer im.ssp.Unlock()

	err = im.ssp.BeginStop()
	if err != nil {
		return err
	}

	close(im.tasks)
	im.wg.Wait()

	im.ssp.CompleteStop()

	return nil
}

func (im *IncidentManager) ReportIncident(itype cmi.IEnumValue, email simple.Email, userIPA net.IP) {
	incident := inc.NewIncidentWithFields(itype, email, userIPA)
	im.tasks <- incident
}

func (im *IncidentManager) logError(err error) {
	if err == nil {
		return
	}

	log.Println(err)
}

func (im *IncidentManager) saveIncident(inc derived2.IIncident) (err error) {
	if inc.GetUserIPA() == nil {
		err = im.dbo.SaveIncidentWithoutUserIPA(m.NewModuleWithValue(ev.NewEnumValue(app.ModuleId_NM)), inc.GetType(), inc.GetEmail())
	} else {
		err = im.dbo.SaveIncident(m.NewModuleWithValue(ev.NewEnumValue(app.ModuleId_NM)), inc.GetType(), inc.GetEmail(), inc.GetUserIPA())
	}
	if err != nil {
		return err
	}

	return nil
}

func (im *IncidentManager) informGateway(inc derived2.IIncident) (re *jrm1.RpcError) {
	blockTime := im.blockTimePerIncidentType[inc.GetType().GetValue().RawValue()]

	// Some incidents are only statistical.
	if blockTime == 0 {
		return nil
	}

	// Some incidents may have an empty IP address.
	// By the way, Go language does not even check anything and returns the
	// `<nil>` string if the IP address is empty. This is a very serious bug in
	// the language, but developers are too stupid to understand this.
	// https://github.com/golang/go/issues/39516
	if inc.GetUserIPA() == nil {
		return nil
	}

	// Other incidents must be directed to the Gateway module.
	var params = gm.BlockIPAddressParams{
		UserIPA:      simple.IPAS(inc.GetUserIPA().String()),
		BlockTimeSec: blockTime,
	}

	var result = new(gm.BlockIPAddressResult)
	var err error
	re, err = im.gwmClient.MakeRequest(context.Background(), gc.FuncBlockIPAddress, params, result)
	if err != nil {
		im.logError(err)
		return jrm1.NewRpcErrorByUser(server2.RpcErrorCode_RPCCall, server2.RpcErrorMsg_RPCCall, nil)
	}
	if re != nil {
		return re
	}
	if !result.OK {
		err = errors.New(fmt.Sprintf(server2.MsgFModuleIsBroken, app.ServiceShortName_GWM))
		im.logError(err)
		return jrm1.NewRpcErrorByUser(server2.RpcErrorCode_RPCCall, server2.RpcErrorMsg_RPCCall, nil)
	}

	return nil
}
