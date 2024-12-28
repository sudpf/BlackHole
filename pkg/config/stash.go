package config

import (
	"bytes"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"gopkg.in/yaml.v2"
)

var (
	GlobalStashConfig StashConfig
	appStashName      string
	appStashBaseDir   string
	appStashLogFile   string
	apiStashLogFile   string
)

type (
	AppConf struct {
		RunMode  string `yaml:"RunMode" json:",default=release"`
		HttpAddr string `yaml:"HttpAddr" json:",default=127.0.0.1:10002"`
	}

	LogConf struct {
		Level string `yaml:"Level" json:",default=info"`
		Size  string `yaml:"Size" json:",default=256m"`
		Dir   string `yaml:"Dir" json:",optional"`
		Gin   string `yaml:"Gin" json:",optional"`
	}

	ConditionConf struct {
		Key   string
		Value string
		Type  string `json:",default=match,options=match|contains"`
		Op    string `json:",default=and,options=and|or"`
	}

	ElasticSearchConf struct {
		Hosts         []string
		Index         string
		DocType       string `json:",default=doc"`
		TimeZone      string `json:",optional"`
		MaxChunkBytes int    `json:",default=15728640"` // default 15M
		Compress      bool   `json:",default=false"`
		Username      string `json:",optional"`
		Password      string `json:",optional"`
	}

	SyslogAddrConf struct {
		Protocol string   `yaml:"Protocol" json:",default=Udp"`
		Address  string   `yaml:"Address"`
		Port     int      `yaml:"Port"`
		Columns  []string `yaml:"Columns" json:",optional"`
	}

	SyslogOutputConf struct {
		Conditions  [][]*ConditionConf `yaml:"Conditions" json:",optional"`
		SyslogAddrs []*SyslogAddrConf  `yaml:"SyslogAddrs" json:",optional"`
	}

	FilterConf struct {
		Action     string          `json:",options=drop|remove_field|transfer"`
		Conditions []ConditionConf `json:",optional"`
		Fields     []string        `json:",optional"`
		Field      string          `json:",optional"`
		Target     string          `json:",optional"`
	}

	SyslogServiceConf struct {
		Protocol   string `yaml:"Protocol" json:",default=Udp"`
		Ssl        string `yaml:"Ssl,omitempty" json:"Ssl,options=on|off,default=off"`
		Address    string `yaml:"Address,omitempty" json:",optional"`
		Port       int    `yaml:"Port,omitempty" json:",optional"`
		Processors int    `yaml:"Processors" json:",default=2"`
	}

	KafkaConf struct {
		service.ServiceConf
		Brokers    []string
		Group      string
		Topics     []string
		Offset     string `json:",options=first|last,default=last"`
		Conns      int    `json:",default=1"`
		Consumers  int    `json:",default=8"`
		Processors int    `json:",default=8"`
		MinBytes   int    `json:",default=10240"`    // 10K
		MaxBytes   int    `json:",default=10485760"` // 10M
		Username   string `json:",optional"`
		Password   string `json:",optional"`
	}

	ClickHouseAuthConf struct {
		Database string `yaml:"Database" json:"Database"`
		Username string `yaml:"Username" json:"Username"`
		Password string `yaml:"Password" json:"Password"`
	}

	ClickHouseConf struct {
		Addr           []string           `yaml:"Addr" json:"Addr"`
		Auth           ClickHouseAuthConf `yaml:"Auth" json:"Auth"`
		Table          string             `yaml:"Table" json:"Table"`
		Columns        []string           `yaml:"Columns" json:",optional"`
		FillNoneColumn bool               `yaml:"FillNoneColumn" json:",default=true"`
		Interval       int64              `yaml:"Interval" json:"interval,default=15"`
		MaxChunkBytes  int                `yaml:"MaxChunkBytes" json:",default=15728640"` // default 15M
	}

	InputConf struct {
		Syslogs []*SyslogServiceConf `yaml:"Syslogs,omitempty" json:",optional"`
		Kafka   *KafkaConf           `yaml:"Kafka,omitempty" json:",optional"`
	}

	OutputConf struct {
		ElasticSearch *ElasticSearchConf  `yaml:"ElasticSearch,omitempty" json:",optional"`
		Syslogs       []*SyslogOutputConf `yaml:"Syslogs,omitempty" json:",optional"`
		Clickhouse    *ClickHouseConf     `yaml:"Clickhouse,omitempty" json:",optional"`
	}

	ClusterConf struct {
		Input   *InputConf   `yaml:"Input,omitempty" json:"Input,optional"`
		Filters []FilterConf `yaml:"Filters,omitempty" json:",optional"`
		Output  *OutputConf  `yaml:"Output,omitempty" json:"Output,optional"`
	}

	StashConfig struct {
		App         appConfig      `yaml:"App" json:"App"`
		Log         logConfig      `yaml:"Log" json:"Log"`
		Clusters    []*ClusterConf `yaml:"Clusters" json:"Clusters"`
		GracePeriod time.Duration  `yaml:"GracePeriod" json:",default=10s"`
	}
)

func (c *StashConfig) String() string {
	buf := new(bytes.Buffer)
	if err := yaml.NewEncoder(buf).Encode(c); err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

func (c *StashConfig) AppLogFile() string {
	return appStashLogFile
}

func (c *StashConfig) ApiLogFile() string {
	return apiStashLogFile
}

func (c *StashConfig) LogLevel() string {
	return c.Log.Level
}

func (c *StashConfig) LogDir() string {
	return c.Log.Dir
}

func GetStashConfig() *StashConfig {
	return &GlobalStashConfig
}

func ParseStashConfig(file string) error {
	// 获取当前执行文件的路径
	appPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	appStashName = filepath.Base(appPath)

	// 获取绝对路径
	absPath, err := filepath.Abs(appPath)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	// 获取目录路径
	appStashBaseDir := filepath.Dir(absPath)

	if !filepath.IsAbs(file) {
		file = appStashBaseDir + "/../conf/" + file
	}

	conf.MustLoad(file, &GlobalStashConfig)

	if !filepath.IsAbs(GlobalStashConfig.Log.Dir) {
		GlobalStashConfig.Log.Dir = appStashBaseDir + "/../" + GlobalStashConfig.Log.Dir
	}

	appStashLogFile = GlobalStashConfig.Log.Dir + "/" + appStashName + ".log"
	apiStashLogFile = GlobalStashConfig.Log.Dir + "/" + appStashName + "api.log"

	return nil
}
