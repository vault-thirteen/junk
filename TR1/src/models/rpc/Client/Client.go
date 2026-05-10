package rmc

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"

	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationServiceEntry"
	"github.com/vault-thirteen/auxie/number"
)

const (
	Err_ShortNameIsNotSet = "short name is not set"
)

const (
	MsgF_PingingService  = "Pinging the service %s ..."
	Msg_PingAttempt      = "."
	MsgF_ServiceIsBroken = "service %s is broken"
	Msg_Ok               = "OK"
)

type Client struct {
	shortName string
	jc        *jrm1.Client
}

func NewClientFromSettings(settings *ccse.CommonConfigurationServiceEntry, shortName string) (serviceClient *Client, err error) {
	schema := settings.GetParameterAsString(ccp.Schema)
	host := settings.GetParameterAsString(ccp.Host)
	port := settings.GetParameterAsInt(ccp.Port)
	path := settings.GetParameterAsString(ccp.Path)
	enableSelfSignedCertificate := settings.GetParameterAsBool(ccp.EnableSelfSignedCertificate)
	dsn := fmt.Sprintf("%s://%s:%d%s", schema, host, port, path)

	serviceClient, err = newClient(shortName, dsn, enableSelfSignedCertificate)
	if err != nil {
		return nil, err
	}

	return serviceClient, nil
}

// NewClient is a constructor of an RPC client.
// Port in DSN must be explicitly set.
func newClient(shortName string, dsn string, enableSelfSignedCertificate bool) (client *Client, err error) {
	if len(shortName) == 0 {
		return nil, errors.New(Err_ShortNameIsNotSet)
	}

	var dsnUrl *url.URL
	dsnUrl, err = url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	var customHttpClient *http.Client
	if (dsnUrl.Scheme == rm.UrlSchemeHttps) && enableSelfSignedCertificate {
		customHttpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	}

	return newCustomClient(shortName, dsnUrl, customHttpClient)
}

func newCustomClient(shortName string, dsnUrl *url.URL, customHttpClient *http.Client) (client *Client, err error) {
	var port uint16
	port, err = number.ParseUint16(dsnUrl.Port())
	if err != nil {
		return nil, err
	}

	path := dsnUrl.RequestURI()

	var clientSettings *jrm1.ClientSettings
	clientSettings, err = jrm1.NewClientSettings(dsnUrl.Scheme, dsnUrl.Hostname(), port, path, customHttpClient, nil, true)
	if err != nil {
		return nil, err
	}

	var rpcClient *jrm1.Client
	rpcClient, err = jrm1.NewClient(clientSettings)
	if err != nil {
		return nil, err
	}

	client = &Client{
		shortName: shortName,
		jc:        rpcClient,
	}

	return client, nil
}

func (cli *Client) MakeRequest(ctx context.Context, method string, params any, result any) (re *jrm1.RpcError, err error) {
	return cli.jc.Call(ctx, method, params, result)
}

func (cli *Client) Ping(verbose bool) (err error) {
	if verbose {
		fmt.Print(fmt.Sprintf(MsgF_PingingService, cli.shortName))
	}

	var params = rm.PingParams{}
	var iMax = int(math.Ceil(float64(rm.ServicePingAttemptsDurationMinutes) * float64(60) / float64(rm.ServiceNextPingAttemptDelaySec)))

	var result = new(rm.PingResult)
	var re *jrm1.RpcError
	for i := 1; i <= iMax; i++ {
		re, err = cli.MakeRequest(context.Background(), rm.Func_Ping, params, result)
		if (err == nil) && (re == nil) {
			break
		}

		if verbose {
			fmt.Print(Msg_PingAttempt)
		}

		if i < iMax {
			time.Sleep(time.Second * time.Duration(rm.ServiceNextPingAttemptDelaySec))
		}
	}

	if err != nil {
		return err
	}
	if re != nil {
		return re.AsError()
	}
	if !result.OK {
		return errors.New(fmt.Sprintf(MsgF_ServiceIsBroken, cli.shortName))
	}

	if verbose {
		fmt.Println(Msg_Ok)
	}

	return nil
}
