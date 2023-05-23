package database

import (
	"database/sql"
	stdLog "log"
	"mocker/common"
	mockerlog "mocker/util/log"
	"os"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	_ "gorm.io/driver/sqlite"
)

type QueryAble interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

func DatabaseInit(username, password, protocol, host, port, dbname string, timeout,
	idleConCount, maxConCount int, logger *log.Entry,
	logfile *os.File, apiMode common.APIModeType) (*gorm.DB, error) {
	var err error
	conString := getConnectionString(username, password, protocol, host, port, dbname)
	var Orm *gorm.DB
	if Orm, err = gorm.Open("mysql", conString); err != nil {
		logger.WithFields(log.Fields{mockerlog.LOG_COMPONENT: "Database"}).Fatal(err)
		return nil, err
	}
	Orm.DB().SetConnMaxLifetime(time.Second * time.Duration(timeout))
	Orm.SetLogger(stdLog.New(logfile, "", stdLog.LstdFlags|stdLog.Lshortfile))
	if apiMode != common.APIMODE_PROD {
		Orm.LogMode(true)
	}
	Orm.DB().SetMaxIdleConns(idleConCount)
	Orm.DB().SetMaxOpenConns(maxConCount)
	return Orm, err
}

func SqlLiteInit(path string) (*gorm.DB, error) {
	Orm, err := gorm.Open("sqlite3", path)
	return Orm, err
}

func getConnectionString(username, password, protocol, host, port, dbname string) string {
	conString := ""
	if len(username) > 0 {
		conString = username
		if len(password) > 0 {
			conString += ":" + password
		}
		conString += "@"
	}
	if len(protocol) > 0 {
		conString += protocol
	}
	conString += "("
	if len(host) > 0 {
		conString += host
		if len(port) > 0 {
			conString += ":" + port
		}
	}
	conString += ")/"
	if len(dbname) > 0 {
		conString += dbname
	}
	extraParams := "?charset=utf8&parseTime=true&loc=Local"
	conString += extraParams
	return conString
}

func DBNow() time.Time {
	var datetime = time.Now()
	datetime.Format(time.RFC3339)
	return datetime
}

func GetStringFromNullString(nullStr sql.NullString) string {
	if nullStr.Valid {
		return nullStr.String
	}
	return ""
}

func GetStringFromNullFloat(nullFloat sql.NullFloat64) string {
	if nullFloat.Valid {
		return strconv.FormatFloat(nullFloat.Float64, 'f', -1, 64)
	}
	return ""
}

func GetFloatFromNullFloat(nullFloat sql.NullFloat64) float64 {
	if nullFloat.Valid {
		return nullFloat.Float64
	}
	return 0
}

func GetIntFromNullInt(nullInt sql.NullInt64) int64 {
	if nullInt.Valid {
		return nullInt.Int64
	}
	return 0
}

func GetNullableString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}
	return sql.NullString{
		String: value,
		Valid:  true,
	}
}

func IsNullStringEmpty(value sql.NullString) bool {
	return value.Valid == false || value.String == ""
}

func GetNullableBool(value *bool) sql.NullBool {
	if value == nil {
		return sql.NullBool{}
	}
	return sql.NullBool{
		Bool:  *value,
		Valid: true,
	}
}

func GetNullableInt64(value *int64) sql.NullInt64 {
	if value == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: *value,
		Valid: true,
	}
}

func GetNullableFloat64(value *float64) sql.NullFloat64 {
	if value == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{
		Float64: *value,
		Valid:   true,
	}
}

func GetNullableTime(value *time.Time) mysql.NullTime {
	if value == nil {
		return mysql.NullTime{}
	}
	return mysql.NullTime{
		Time:  *value,
		Valid: true,
	}
}

func GetBoolFromNullBool(nullBool sql.NullBool) bool {
	return nullBool.Valid && nullBool.Bool
}
