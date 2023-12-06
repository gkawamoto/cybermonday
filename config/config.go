package config

import (
	"os"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr                string        `envconfig:"CYBERMONDAY_ADDR" default:":8080"`
	StaticDir           string        `envconfig:"CYBERMONDAY_STATIC_DIR" default:"./resources"`
	BasePath            string        `envconfig:"CYBERMONDAY_BASEPATH" default:"."`
	Template            string        `envconfig:"CYBERMONDAY_TEMPLATE"`
	TemplatePath        string        `envconfig:"CYBERMONDAY_TEMPLATE_PATH"`
	DefaultTemplatePath string        `envconfig:"CYBERMONDAY_DEFAULT_TEMPLATE_PATH" default:"./resources/index.tplt.html"`
	ShutdownTimeout     time.Duration `envconfig:"CYBERMONDAY_SHUTDOWN_TIMEOUT" default:"30s"`

	Envs map[string]string `envconfig:"-"`
}

func New() (*Config, error) {
	config := &Config{
		Envs: map[string]string{},
	}

	if err := envconfig.Process("", config); err != nil {
		return nil, err
	}

	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		config.Envs[pair[0]] = pair[1]
	}

	return config, nil
}
