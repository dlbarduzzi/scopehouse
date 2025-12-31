package scopehouse

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dlbarduzzi/scopehouse/core"
	"github.com/spf13/viper"
)

// Ensures that the ScopeHouse implements the App interface.
var _ core.App = (*ScopeHouse)(nil)

type ScopeHouse struct {
	core.App

	// Logger configs
	logLevel  string
	logFormat string

	// Server configs.
	serverPort         int
	serverIdleTimeout  time.Duration
	serverReadTimeout  time.Duration
	serverWriteTimeout time.Duration

	// Database configs.
	databaseUrl string
}

// Config is the ScopeHouse initialization config struct.
type Config struct {
	// Logger configs
	LogLevel  string
	LogFormat string

	// Server configs.
	ServerPort         int
	ServerIdleTimeout  time.Duration
	ServerReadTimeout  time.Duration
	ServerWriteTimeout time.Duration
}

func New() *ScopeHouse {
	return NewWithConfig(Config{})
}

func NewWithConfig(config Config) *ScopeHouse {
	sh := &ScopeHouse{
		logLevel:           config.LogLevel,
		logFormat:          config.LogFormat,
		serverPort:         config.ServerPort,
		serverIdleTimeout:  config.ServerIdleTimeout,
		serverReadTimeout:  config.ServerReadTimeout,
		serverWriteTimeout: config.ServerWriteTimeout,
	}

	sh.parseConfig(&config)

	sh.App = core.NewBaseApp(core.BaseAppConfig{
		LogLevel:  sh.logLevel,
		LogFormat: sh.logFormat,
	})

	return sh
}

func (sh *ScopeHouse) Start() error {
	if err := sh.Bootstrap(); err != nil {
		return err
	}

	return nil
}

func (sh *ScopeHouse) parseConfig(config *Config) {
	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvPrefix("SH")

	ec := new(envConfig)
	ec.register(v)
	sh.readConfigs(v, ec, config)

	es := new(envSecret)
	es.register(v)
	sh.readSecrets(v, es)
}

func (sh *ScopeHouse) readConfigs(v *viper.Viper, e *envConfig, config *Config) {
	v.AddConfigPath(e.filePath)
	v.SetConfigName(e.fileName)
	v.SetConfigType(e.fileType)

	// sanitizeConfig will register initial ScopeHouse config with some validation.
	sh.sanitizeConfig(config)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Fprintf(
				os.Stderr, "[warn] config file not found - using default values - %s\n",
				err,
			)
		} else {
			fmt.Fprintf(
				os.Stderr, "[error] failed to read configs - using default values - %s\n",
				err,
			)
		}
		return
	}

	// Logger configs.
	sh.logLevel = v.GetString("LOG_LEVEL")
	sh.logFormat = v.GetString("LOG_FORMAT")

	// Server configs.
	sh.serverPort = v.GetInt("SERVER_PORT")
	sh.serverIdleTimeout = v.GetDuration("SERVER_IDLE_TIMEOUT_SEC")
	sh.serverReadTimeout = v.GetDuration("SERVER_READ_TIMEOUT_SEC")
	sh.serverWriteTimeout = v.GetDuration("SERVER_WRITE_TIMEOUT_SEC")
}

func (sh *ScopeHouse) sanitizeConfig(config *Config) {
	sh.logLevel = strings.TrimSpace(config.LogLevel)
	if sh.logLevel == "" {
		sh.logLevel = "info"
	}

	sh.logFormat = strings.TrimSpace(config.LogFormat)
	if sh.logFormat == "" {
		sh.logFormat = "json"
	}

	sh.serverPort = config.ServerPort
	if sh.serverPort < 1 {
		sh.serverPort = 8090
	}

	sh.serverIdleTimeout = config.ServerIdleTimeout
	if sh.serverIdleTimeout < 1 {
		sh.serverIdleTimeout = 5
	}

	sh.serverReadTimeout = config.ServerReadTimeout
	if sh.serverReadTimeout < 1 {
		sh.serverReadTimeout = 5
	}

	sh.serverWriteTimeout = config.ServerWriteTimeout
	if sh.serverWriteTimeout < 1 {
		sh.serverWriteTimeout = 5
	}
}

func (sh *ScopeHouse) readSecrets(v *viper.Viper, e *envSecret) {
	v.AddConfigPath(e.filePath)
	v.SetConfigName(e.fileName)
	v.SetConfigType(e.fileType)

	// error handling is deferred to the upstream service when secrets are invalid.
	_ = v.ReadInConfig()

	sh.databaseUrl = v.GetString("DATABASE_URL")
}

type envConfig struct {
	filePath string
	fileName string
	fileType string
}

func (e *envConfig) register(v *viper.Viper) {
	e.filePath = strings.TrimSpace(v.GetString("CONFIG_PATH"))
	if e.filePath == "" {
		e.filePath = "/etc/scopehouse/configs"
	}

	e.fileName = strings.TrimSpace(v.GetString("CONFIG_NAME"))
	if e.fileName == "" {
		e.fileName = "config"
	}

	e.fileType = strings.TrimSpace(v.GetString("CONFIG_TYPE"))
	if e.fileType == "" {
		e.fileType = "yaml"
	}
}

type envSecret struct {
	filePath string
	fileName string
	fileType string
}

func (e *envSecret) register(v *viper.Viper) {
	e.filePath = strings.TrimSpace(v.GetString("SECRET_PATH"))
	if e.filePath == "" {
		e.filePath = "/etc/scopehouse/secrets"
	}

	e.fileName = strings.TrimSpace(v.GetString("SECRET_NAME"))
	if e.fileName == "" {
		e.fileName = "vault"
	}

	e.fileType = strings.TrimSpace(v.GetString("SECRET_TYPE"))
	if e.fileType == "" {
		e.fileType = "env"
	}
}
