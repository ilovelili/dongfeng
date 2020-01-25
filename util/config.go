package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
)

var once sync.Once
var instance *Config

// GetConfig get config defined in config.json
func GetConfig() *Config {
	once.Do(func() {
		env := os.Getenv("DF_ENVIROMENT")
		if env == "" {
			env = "dev"
		}

		var config *Config
		pwd, _ := os.Getwd()
		filepath := path.Join(pwd, fmt.Sprintf("config.%s.json", strings.ToLower(env)))
		configFile, err := os.Open(filepath)
		defer configFile.Close()
		if err != nil {
			panic(err)
		}

		jsonParser := json.NewDecoder(configFile)
		err = jsonParser.Decode(&config)
		if err != nil {
			panic(err)
		}

		instance = config
	})

	return instance
}

// DataBase database config
type DataBase struct {
	Host     string `json:"host"`
	DataBase string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
	Debug    bool   `json:"debug"`
}

// OSS sso sevice config
type OSS struct {
	APIKey     string `json:"api_key"`
	APISecret  string `json:"api_secret"`
	Endpoint   string `json:"endpoint"`
	BucketName string `json:"bucket_name"`
}

// Aliyun aliyun services
type Aliyun struct {
	OSS `json:"oss"`
}

// Auth authing fields
type Auth struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// Ebook ebook related config
type Ebook struct {
	Width            float64 `json:"width"`
	Height           float64 `json:"height"`
	OriginDir        string  `json:"origin_dir"`
	PDFDestDir       string  `json:"pdf_dest_dir"`
	ImageDestDir     string  `json:"img_dest_dir"`
	MergeTargetDir   string  `json:"merge_target_dir"`
	MergeDestDir     string  `json:"merge_dest_dir"`
	ImageLoadTimeout int64   `json:"image_load_timeout"`
}

// Config config entry
type Config struct {
	Auth     `json:"auth"`
	DataBase `json:"database"`
	Aliyun   `json:"aliyun"`
	Ebook    `json:"ebook"`
}
