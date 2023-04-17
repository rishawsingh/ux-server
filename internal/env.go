package internal

import (
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

var BuildNumber string

type Environment string

const (
	Local      Environment = "local"
	Dev        Environment = "dev"
	Stage      Environment = "stage"
	Production Environment = "prod"
)

func (e Environment) Valid() bool {
	return e == Local ||
		e == Dev ||
		e == Stage ||
		e == Production
}

func (e Environment) String() string {
	return string(e)
}

type SSLMode string

const (
	SSLModeEnable  SSLMode = "enable"
	SSLModeDisable SSLMode = "disable"
)

func (ssl SSLMode) Valid() bool {
	return ssl == SSLModeEnable ||
		ssl == SSLModeDisable
}

// EnvConfig has environment stored
type EnvConfig struct {
	Environment                   Environment   `env:"ENV"`
	ServerPort                    string        `env:"PORT"`
	DatabaseHost                  string        `env:"DB_HOST,required"`
	DatabasePort                  string        `env:"DB_PORT,required"`
	DatabaseName                  string        `env:"DB_NAME,required"`
	DatabaseUserName              string        `env:"DB_USERNAME,required"`
	DatabaseUserPassword          string        `env:"DB_USER_PASSWORD,required"`
	DatabaseMaxConnection         int           `env:"DB_MAX_CONN"`
	DatabaseMaxIdleConnection     int           `env:"DB_MAX_IDLE_CONN"`
	DatabaseMaxConnectionLifeTime time.Duration `env:"DB_MAX_CONN_LIFE"`
	DatabaseSSLMode               SSLMode       `env:"DB_SSL_MODE"`
	DatabaseMaxQuerySecond        int           `env:"DB_MAX_QUERY_SECOND"`
	LogLevel                      string        `env:"LOG_LEVEL"`
	Timezone                      string        `env:"TZ"`
	Trace                         bool          `env:"TRACE"`
	TracerEndpoint                string        `env:"TRACER_ENDPOINT"`
	TracerPort                    string        `env:"TRACER_PORT"`
	ServiceName                   string        `env:"SERVICE_NAME"`
}

// NewEnvConfig creates a new environment
func NewEnvConfig() *EnvConfig {
	config := EnvConfig{}
	if err := env.Parse(&config); err != nil {
		log.Fatalf("failed to load env config with error: %+v", err)
	}
	if !config.Environment.Valid() {
		// fallback to local
		config.Environment = Local
	}
	if config.ServerPort == "" {
		// fallback to 8081
		config.ServerPort = "8081"
	}
	if config.LogLevel == "" {
		// fallback to debug
		config.LogLevel = "debug"
	}
	if !config.DatabaseSSLMode.Valid() {
		// fallback to disable SSL
		config.DatabaseSSLMode = SSLModeDisable
	}
	if config.DatabaseMaxQuerySecond == 0 {
		// fallback to max 2 second
		config.DatabaseMaxQuerySecond = 2
	}
	if config.Timezone == "" {
		// fallback to UTC
		config.Timezone = time.UTC.String()
	}
	if config.Trace && (config.TracerEndpoint == "" || config.TracerPort == "") {
		log.Fatal("tracer endpoint and tracer port is required when tracing is enabled")
	}
	if config.DatabaseMaxIdleConnection > config.DatabaseMaxConnection {
		log.Fatal("max idle connection can not be greater than max connection")
	}
	if config.DatabaseMaxConnection == 0 {
		config.DatabaseMaxConnection = 10
	}
	if config.DatabaseMaxIdleConnection == 0 {
		config.DatabaseMaxIdleConnection = 3
	}
	if config.DatabaseMaxConnectionLifeTime == 0 {
		config.DatabaseMaxConnectionLifeTime = 10
	}
	return &config
}

func (e *EnvConfig) BuildNumber() string {
	if BuildNumber == "" {
		return string(Local)
	}
	return BuildNumber
}
