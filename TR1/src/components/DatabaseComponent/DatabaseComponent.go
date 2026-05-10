package dc

import (
	"database/sql"
	"fmt"
	"net"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/vault-thirteen/TR1/src/interfaces"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	ErrF_DatabaseTypeIsNotSupported = "database type is not supported: %s"
)

const (
	Msg_ConnectingDatabase = "Connecting with database server ... "
	Msg_Failure            = "Failure"
	Msg_OK                 = "OK"
)

const (
	DbType_MySQL = "mysql"
)

type DatabaseComponent struct {
	cfg    interfaces.IConfiguration
	dbType string
	sqlDb  *sql.DB
	gormDb *gorm.DB
}

func (c *DatabaseComponent) Init(cfg interfaces.IConfiguration, controller interfaces.IController) (sc interfaces.IServiceComponent, err error) {
	systemSettings := cfg.GetComponent(cm.Component_System, cm.Protocol_None)
	databaseType := systemSettings.GetParameterAsString(ccp.DatabaseType)

	switch databaseType {
	case DbType_MySQL:
		return c.initWithMysql(cfg)
	}

	return nil, fmt.Errorf(ErrF_DatabaseTypeIsNotSupported, databaseType)
}
func (c *DatabaseComponent) GetConfiguration() interfaces.IConfiguration {
	return c.cfg
}

func (c *DatabaseComponent) initWithMysql(cfg interfaces.IConfiguration) (sc interfaces.IServiceComponent, err error) {
	fmt.Print(Msg_ConnectingDatabase)
	defer func() {
		if err != nil {
			fmt.Print(Msg_Failure)
		} else {
			fmt.Println(Msg_OK)
		}
	}()

	dc := &DatabaseComponent{
		cfg:    cfg,
		dbType: DbType_MySQL,
	}

	dbSettings := cfg.GetComponent(cm.Component_Database, cm.Protocol_MySQL)

	driverName := dbSettings.GetParameterAsString(ccp.DriverName)
	host := dbSettings.GetParameterAsString(ccp.Host)
	port := dbSettings.GetParameterAsInt(ccp.Port)
	addr := net.JoinHostPort(host, strconv.Itoa(port))

	mc := mysql.Config{
		Net:                  dbSettings.GetParameterAsString(ccp.Net),
		Addr:                 addr,
		DBName:               dbSettings.GetParameterAsString(ccp.DatabaseName),
		User:                 dbSettings.GetParameterAsString(ccp.User),
		Passwd:               dbSettings.GetParameterAsString(ccp.Password),
		AllowNativePasswords: dbSettings.GetParameterAsBool(ccp.AllowNativePasswords),
		CheckConnLiveness:    dbSettings.GetParameterAsBool(ccp.CheckConnLiveness),
		MaxAllowedPacket:     dbSettings.GetParameterAsInt(ccp.MaxAllowedPacket),
		Params:               dbSettings.GetParameterAsMap(ccp.Params),
	}

	dc.sqlDb, err = sql.Open(driverName, mc.FormatDSN())
	if err != nil {
		return nil, err
	}

	err = dc.sqlDb.Ping()
	if err != nil {
		return nil, err
	}

	dc.gormDb, err = gorm.Open(gmysql.New(gmysql.Config{Conn: dc.sqlDb}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return dc, nil
}

func (c *DatabaseComponent) Start(s interfaces.IService) (err error) {
	var sqlDb *sql.DB
	sqlDb, err = c.gormDb.DB()
	if err != nil {
		return err
	}

	err = sqlDb.Ping()
	if err != nil {
		return err
	}

	return nil
}
func (c *DatabaseComponent) Stop(s interfaces.IService) (err error) {
	wg := s.GetSubRoutinesWG()
	defer wg.Done()

	var sqlDb *sql.DB
	sqlDb, err = c.gormDb.DB()
	if err != nil {
		return err
	}

	err = sqlDb.Close()
	if err != nil {
		return err
	}

	c.ReportStop()

	return nil
}

func (c *DatabaseComponent) ReportStart() {
	fmt.Println("DatabaseComponent has started")
}
func (c *DatabaseComponent) ReportStop() {
	fmt.Println("DatabaseComponent has stopped")
}

// Other methods.

func FromAny(x any) (c *DatabaseComponent) {
	return x.(*DatabaseComponent)
}

// Non-standard methods.

func (c *DatabaseComponent) GetGormDb() (gormDb *gorm.DB) {
	return c.gormDb
}
