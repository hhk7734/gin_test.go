package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type config struct {
	Debug bool      `env:"DEBUG" envDefault:"false"`
	Port  int       `env:"PORT" envDefault:"8080"`
	K8s   k8sConfig `envPrefix:"K8S_"`
}

type k8sConfig struct {
	PodName      string `env:"POD_NAME,required"`
	PodNamespace string `env:"POD_NAMESPACE,required"`
}

var c config

func Init() {
	var err error
	err = godotenv.Load()
	if err != nil {
		fmt.Println("failed to load .env", err)
	}

	err = env.Parse(&c)
	if err != nil {
		fmt.Println("failed to parse env", err)
	}
}

func Config() config {
	return c
}
