package todo_store

import (
	"RobertTC32/example-demo_hello/src/commons"
	"errors"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// "store" (aka repository) contains adapter for MariaDB database in the "Hexagonal Architecture"

type MysqlStore struct {
	Dsn string
	Db  *sqlx.DB
}

func NewMysqlStore() (*MysqlStore, error) {
	slog.Debug("store::NewMysqlStore() - Executing")
	//dsnFormat := "%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local"
	//dsn := fmt.Sprintf(dsnFormat, Username, Password, Host, Port, DbName)
	dbdriver := os.Getenv("DB_DRIVER")
	dsn := os.Getenv("DB_DSN")
	if commons.IsRunningInDockerContainer() {
		// change database localhost to host.docker.internal
		dsn = strings.Replace(dsn, "localhost", "host.docker.internal", 1)
	}
	db, err := sqlx.Open(dbdriver, dsn)
	if err != nil {
		slog.Error("store::NewMysqlStore() - Failed to open SQL connection", "error", err)
		return &MysqlStore{}, commons.NewWrappingError(errors.New("Failed to open SQL connection"), err)
	}
	// Pool configuration
	poolFlag, err := strconv.Atoi(os.Getenv("DB_POOL_MAX_OPEN_CONNS"))
	if err == nil {
		slog.Debug("store::NewMysqlStore() - Set", "maxOpenConns", poolFlag)
		db.SetMaxOpenConns(poolFlag) // max simultaneous connections
	}
	poolFlag, err = strconv.Atoi(os.Getenv("DB_POOL_MAX_IDLE_CONNS"))
	if err == nil {
		slog.Debug("store::NewMysqlStore() - Set", "maxIdleConns", poolFlag)
		db.SetMaxIdleConns(poolFlag) // connections kept idle in pool
	}
	poolFlag, err = strconv.Atoi(os.Getenv("DB_POOL_CONN_MAX_LIFETIME"))
	if err == nil {
		slog.Debug("store::NewMysqlStore() - Set", "connMaxLifetime", poolFlag)
		db.SetConnMaxLifetime(time.Duration(poolFlag) * time.Minute) // max connection age
	}
	poolFlag, err = strconv.Atoi(os.Getenv("DB_POOL_CONN_MAX_IDLE_TIME"))
	if err == nil {
		slog.Debug("store::NewMysqlStore() - Set", "connMaxIdleTime", poolFlag)
		db.SetConnMaxIdleTime(time.Duration(poolFlag) * time.Minute) // evict idle connections after this
	}
	err = db.Ping()
	if err != nil {
		slog.Error("store::NewMysqlStore() - Failed to ping db", "error", err)
		return &MysqlStore{}, commons.NewWrappingError(errors.New("Failed to ping db"), err)
	}
	slog.Debug("store::NewMysqlStore() - Database opened", "dsn", dsn)
	return &MysqlStore{
		Dsn: dsn,
		Db:  db,
	}, nil
}

func (this *MysqlStore) Destroy() error {
	slog.Debug("store::Destroy() - Executing")
	return this.Db.Close()
}
