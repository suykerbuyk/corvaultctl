package main

import (
	//"errors"
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/rs/zerolog/log"

	//"log"
	"os"

	"gopkg.in/yaml.v2"
)

const default_config_file_name string = ".stx_corvault.yaml"

var configFile string //       `yaml:"config_file"`

// Set up default configuration file path
func init() {
	userhome, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	configFile = userhome + "/" + default_config_file_name
}

// Where we will try to pull our yaml settings file from

type CorvaultCredential struct {
	Host string `yaml:"host"`
	User string `yaml:"user"`
	Auth string `yaml:"auth"`
	Pass string `yaml:"pass"`
	Key  string `yaml:"key"`
}

func (c *CorvaultCredential) SetAuth(password string) {
	c.Auth = base64.StdEncoding.EncodeToString([]byte(c.User + ":" + password))
}

type CorvaultConfig struct {
	Targets map[string]CorvaultCredential `yaml:"targets"`
}

func (c *CorvaultConfig) String() (s string) {
	b, err := c.PrettyPrint()
	if err != nil {
		log.Fatal().Err(err).Msg("CorvaultConfig to string")
	}
	return string(b)
}
func (c *CorvaultConfig) PrettyPrint() ([]byte, error) {
	//b, err := yaml.Marshal(c.Target)
	b, err := yaml.Marshal(c)
	if err != nil {
		err = fmt.Errorf("Error decoding config data: " + err.Error())
		return b, err
	}
	return b, err
}

func NewCorvaultConfig() (cfg *CorvaultConfig) {
	cfg = &CorvaultConfig{}
	cfg.Targets = make(map[string]CorvaultCredential)
	return cfg
}

func validateConfigPath(path string) (err error) {
	log.Info().Msg("validate config path")
	s, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Config Error %w", err)
		return
	}
	if s.IsDir() {
		return fmt.Errorf("Config: '%s' is a directory, not a normal file", path)
	}
	return nil
}

func SaveCvtConfig(cfg *CorvaultConfig) (err error) {
	pp, err := cfg.PrettyPrint()
	err = ioutil.WriteFile(configFile, pp, 0666)
	return
}
func GetCvtConfig() (cfg *CorvaultConfig, err error) {
	cfg = NewCorvaultConfig()
	err = validateConfigPath(configFile)
	if err != nil {
		err = SaveCvtConfig(cfg)
		if err != nil {
			err = fmt.Errorf("Could not create config file: %s %w", configFile, err)
			return
		}
	}
	file, err := os.Open(configFile)
	if err != nil {
		err = fmt.Errorf("Config Error: %w", err)
		return
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	if err = d.Decode(&cfg); err != nil {
		err = fmt.Errorf("Config Error %w", err)
		return
	}
	return
}
