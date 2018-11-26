package conf

import (
	"fmt"
	"log"

	"github.com/toolkits/file"
	yaml "gopkg.in/yaml.v2"
)

type Nagios struct {
	Addr        string `yaml:"addr"`
	Timeout     int    `yaml:"timeout"`
	TemplateDir string `yaml:"tempaltedir"`
}

type AlertBaet struct {
	Debug  bool     `yaml:"debug"`
	Web    string   `yaml:"web"`
	Ignore []string `yaml:"ignore"`
	Nagios *Nagios  `yaml:"nagios"`
}

type AlertConf struct {
	AlertBaet *AlertBaet `yaml:"alertbeat"`
}

var (
	Config *AlertBaet
)

func Parse(cfg string) error {
	if cfg == "" {
		return fmt.Errorf("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		return fmt.Errorf("configuration file %s is nonexistent", cfg)
	}

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		return fmt.Errorf("read configuration file %s fail %s", cfg, err.Error())
	}

	var c AlertConf
	err = yaml.Unmarshal([]byte(configContent), &c)
	if err != nil {
		return fmt.Errorf("parse configuration file %s fail %s", cfg, err.Error())
	}

	Config = c.AlertBaet
	log.Println("[INFO] load configuration file", cfg, "successfully")
	return nil
}
