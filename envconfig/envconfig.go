package envconfig

import (
	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	AppPort int `envconfig:"APP_PORT"    default:"8080"`
	// DBHost         string `envconfig:"DB_HOST"     default:"localhost"`
	// DBPort         int    `envconfig:"DB_PORT"     default:"5555"`
	// DBName         string `envconfig:"DB_NAME"     default:"test"`
	// DBUser         string `envconfig:"DB_USER"     default:"test"`
	// DBPassword     string `envconfig:"DB_PASSWORD" default:"test"`
}

func Process() (env Env, err error) {
	err = envconfig.Process("", &env)
	return
}
