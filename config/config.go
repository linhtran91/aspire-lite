package config

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server *Server `yaml:"server" mapstructure:"server"`
	DB     *DB     `yaml:"db"  mapstructure:"db"`
	JWT    *JWT    `yaml:"jwt"  mapstructure:"jwt"`
}

type Server struct {
	Host string `yaml:"host" mapstructure:"host"`
	Port string `yaml:"port" mapstructure:"port"`
}

type DB struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     string `yaml:"port" mapstructure:"port"`
	Database string `yaml:"database" mapstructure:"database"`
	User     string `yaml:"user" mapstructure:"user"`
	Password string `yaml:"password" mapstructure:"password"`
}

type JWT struct {
	Secret   string        `json:"secret" mapstructure:"secret"`
	Duration time.Duration `json:"duration" mapstructure:"duration"`
}

func (cfg *Config) BuildDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Database)
}

func NewConfig(configPath string) (*Config, error) {
	// config := &Config{}
	// file, err := os.Open(configPath)
	// if err != nil {
	// 	return nil, err
	// }
	// defer file.Close()

	// d := yaml.NewDecoder(file)
	// if err := d.Decode(&config); err != nil {
	// 	return nil, err
	// }
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func ParseFlags() (string, error) {
	var configPath string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", "./config/", "path to config file")

	flag.Parse()

	return configPath, nil
}
