package internal

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4/database/postgres"

	"github.com/lib/pq"

	"github.com/friendsofgo/errors"

	"github.com/samber/lo"

	"github.com/qustavo/sqlhooks/v2"

	"github.com/jmoiron/sqlx"

	// source/file import is required for migration files to read
	_ "github.com/golang-migrate/migrate/v4/source/file"

	migrator "github.com/golang-migrate/migrate/v4"
)

type Database struct {
	*sqlx.DB
	env *EnvConfig
}

type hooks struct {
	logger *Logger
	env    *EnvConfig
	sqlhooks.OnErrorer
}

type dbCtxKey string

const (
	begin dbCtxKey = "begin"
)

func (h *hooks) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	return context.WithValue(ctx, begin, time.Now()), nil
}
func (h *hooks) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	queryStart, okay := ctx.Value(begin).(time.Time)
	if !okay || time.Since(queryStart).Seconds() > float64(h.env.DatabaseMaxQuerySecond) {
		cleanQuery := strings.ReplaceAll(query, "\t", "")
		cleanQuery = strings.ReplaceAll(cleanQuery, "\n", "")
		cleanQuery = strings.Join(strings.Fields(cleanQuery), " ")
		timeTaken := lo.Ternary[string](okay, fmt.Sprintf("%d", time.Since(queryStart).Milliseconds()), "")
		h.logger.WithCtx(ctx).
			With("query", cleanQuery).
			With("query args", args).
			With("time-ms", timeTaken).
			Error("long running query")
	}
	return ctx, nil
}
func (h *hooks) OnError(ctx context.Context, err error, query string, args ...interface{}) error {
	queryStart, okay := ctx.Value(begin).(time.Time)
	cleanQuery := strings.ReplaceAll(query, "\t", " ")
	cleanQuery = strings.ReplaceAll(cleanQuery, "\n", "")
	cleanQuery = strings.Join(strings.Fields(cleanQuery), " ")
	timeTaken := lo.Ternary[string](okay, fmt.Sprintf("%d", time.Since(queryStart).Milliseconds()), "")
	h.logger.WithCtx(ctx).
		With("query", cleanQuery).
		With("query args", args).
		With("error", err).
		With("time-ms", timeTaken).
		Debug("error executing query")
	return err
}

func newHook(logger *Logger, env *EnvConfig) hooks {
	return hooks{
		logger: logger,
		env:    env,
	}
}

func (h *hooks) DriverName() string {
	return "postgresWithHook"
}

func (h *hooks) Driver(d driver.Driver) driver.Driver {
	return sqlhooks.Wrap(d, h)
}

func NewDatabase(env *EnvConfig, logger *Logger, tracer *Tracer) (*Database, error) {
	appName := env.Environment.String() + "-" + env.ServiceName
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s application_name=%s",
		env.DatabaseHost, env.DatabasePort, env.DatabaseUserName, env.DatabaseUserPassword, env.DatabaseName,
		env.DatabaseSSLMode, appName)
	var db *sql.DB
	hook := newHook(logger, env)
	sql.Register(hook.DriverName(), hook.Driver(pq.Driver{}))
	if tracer.enabled {
		hookedAndTracedDB, connectErr := tracer.WrapDB(hook.DriverName(), env.DatabaseName, connStr)
		if connectErr != nil {
			return nil, errors.WithMessage(connectErr, "failed to connect with database")
		}
		db = hookedAndTracedDB
	} else {
		hookedOnlyDB, connectErr := sql.Open(hook.DriverName(), connStr)
		if connectErr != nil {
			return nil, errors.WithMessage(connectErr, "failed to connect with database")
		}
		db = hookedOnlyDB
	}
	db.SetMaxOpenConns(env.DatabaseMaxConnection)
	db.SetMaxIdleConns(env.DatabaseMaxIdleConnection)
	db.SetConnMaxLifetime(time.Minute * env.DatabaseMaxConnectionLifeTime)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	pingErr := db.PingContext(ctx)
	if pingErr != nil && pingErr != driver.ErrSkip {
		return nil, errors.WithMessage(pingErr, "failed to ping database")
	}
	_, appNameErr := db.Exec(fmt.Sprintf("SET application_name = '%s'", appName))
	if appNameErr != nil {
		return nil, errors.WithMessage(appNameErr, "failed to set app name with database")
	}
	logger.Info("database connected")
	return &Database{DB: sqlx.NewDb(db, hook.DriverName()), env: env}, nil
}

func (database *Database) Close() error {
	return database.DB.Close()
}

func (database *Database) MigrateUp(path string) error {
	dbDriver, err := postgres.WithInstance(database.DB.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrator.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", path),
		"postgres", dbDriver)

	if err != nil {
		return err
	}
	if migrateUpErr := m.Up(); migrateUpErr != nil && migrateUpErr != migrator.ErrNoChange {
		return migrateUpErr
	}
	return nil
}
