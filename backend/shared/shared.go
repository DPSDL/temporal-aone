package shared

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	Config Configuration
	Logger = logrus.New()
)

// Configuration struct for application configuration
type Configuration struct {
	Database DatabaseConfig
	Server   ServerConfig
	Log      LogConfig
}

// DatabaseConfig struct for database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

// ServerConfig struct for server configuration
type ServerConfig struct {
	Port int
}

// LogConfig struct for log configuration
type LogConfig struct {
	Level  string
	Format string
}

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("backend/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Error reading config file: %s", err))
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(fmt.Sprintf("Unable to decode into struct: %v", err))
	}
}

func InitLogger() {
	level, err := logrus.ParseLevel(Config.Log.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	if Config.Log.Format == "json" {
		Logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		Logger.SetFormatter(&logrus.TextFormatter{})
	}

	Logger.Out = os.Stdout
}

func InitDatabase() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Config.Database.User, Config.Database.Password, Config.Database.Host, Config.Database.Port, Config.Database.Name)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the database: " + err.Error())
	}
}

func GetDB() *gorm.DB {
	return DB
}
