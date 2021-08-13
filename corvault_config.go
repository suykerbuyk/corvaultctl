package main

import (
	//"errors"
	"encoding/base64"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

var configFile string //       `yaml:"config_file"`
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
		log.Fatal(fmt.Errorf("CorvaultConfig to string: %w", err))
	}
	return string(b)
}
func (c *CorvaultConfig) PrettyPrint() ([]byte, error) {
	//b, err := yaml.Marshal(c.Target)
	b, err := yaml.Marshal(c)
	if err != nil {
		err = fmt.Errorf("ResponseStatus to JSON string error: " + err.Error())
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

func parseFlags(cfg *CorvaultConfig) error {
	var configPath, nick, uri, user, pass string
	flag.StringVar(&configPath, "config", "/home/johns/.cvt.config.yml", "path to config file")
	flag.StringVar(&nick, "nickname", "corvault-2a", "Nickname of a corvault target node")
	flag.StringVar(&uri, "uri", "https://corvault-2a/", "URI of corvault target node")
	flag.StringVar(&user, "user", "manage", "user account name on corvault target")
	flag.StringVar(&pass, "pass", "Testit123!", "password of the named user account")

	// Actually parse the flags
	flag.Parse()

	haveAllTgtFlags := (nick != "" && uri != "" && user != "" && pass != "")
	haveSomeTgtFlags := (nick != "" || uri != "" || user != "" || pass != "")
	if haveSomeTgtFlags {
		if !haveAllTgtFlags {
			fmt.Println("Not All flags present")
			err := fmt.Errorf("Not all target flags supplied")
			return err
		}
		var cvtCredential = CorvaultCredential{Host: uri, User: user, Pass: pass}
		cfg.Targets[nick] = cvtCredential
	}
	fmt.Println(cfg.String())
	configFile = configPath
	return nil
}
func SaveCvtConfig(cfg *CorvaultConfig) (err error) {
	pp, err := cfg.PrettyPrint()
	err = ioutil.WriteFile(configFile, pp, 0666)
	return
}
func GetCvtConfig() (cfg *CorvaultConfig, err error) {
	cfg = NewCorvaultConfig()
	parseFlags(cfg)
	//fmt.Printf("%v\n", *cfg)
	err = validateConfigPath(configFile)
	if err != nil {
		err = SaveCvtConfig(cfg)
		if err != nil {
			err = fmt.Errorf("Could not create config file: %w", err)
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
