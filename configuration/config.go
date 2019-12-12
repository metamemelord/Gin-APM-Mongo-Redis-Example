package configuration

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Cache    CacheConfig `yaml:"cache"`
	Database DBConfig    `yaml:"database"`
	Port     int         `yaml:"port"`
}

type CacheConfig struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"database"`
}

type DBConfig struct {
	Hosts            []string `yaml:"hosts"`
	ConnectionString string   `yaml:"connection_string"`
	Database         string   `yaml:"database"`
	Collection       string   `yaml:"collection"`
}

func Get() *Configuration {

	defer func() {
		err := recover()
		if nil != err {
			log.Fatal(err)
		}
	}()

	env := os.Getenv("app_env")
	if env == "" {
		panic(fmt.Errorf("\"app_env\" not set"))
	}
	var file *os.File
	var err error

	cwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	env = strings.ToUpper(env)
	switch env {
	case "DEV":
		file, err = os.Open(strings.Join([]string{cwd, "configuration", "dev-config.yml"}, string(os.PathSeparator)))
	case "RELEASE":
		file, err = os.Open(strings.Join([]string{cwd, "configuration", "release-config.yml"}, string(os.PathSeparator)))
	}

	if err != nil {
		panic(fmt.Errorf("Error while opening the config file for %s env (%s-config.yml)\nERROR: %s", env, strings.ToLower(env), err.Error()))
	}

	config := &Configuration{}
	err = yaml.NewDecoder(file).Decode(config)

	if err != nil {
		panic(fmt.Errorf("Error while parsing the config file for %s env (%s-config.yml)\nERROR: %s", env, strings.ToLower(env), err.Error()))
	}

	log.Println(config)
	return config
}
