package config

import (
	"skripsi-sso/utils"
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
)

const JWT_EXPIRED = time.Hour * 24
const JWT_REFRESH_EXPIRED = time.Hour * 24 * 3

type config struct {
	App       app
	Databases []database
}

type app struct {
	Env           string
	ServerPort    string `toml:"serverport"`
	JwtSecret     string `toml:"jwtsecret"`
	LogPath       string `toml:"log_path"`
	ExportPath    string `toml:"export_path"`
	ApiSettlement string `toml:"api_settlement"`
	UrlStorage    string `toml:"url_storage"`
}

type database struct {
	Name              string `toml:"name"`
	DriverName        string `toml:"driver_name"`
	ConnectionString  string `toml:"connection_string"`
	MaxConnectionOpen int    `toml:"max_connection_open"`
}

var App app
var Databases []database

// mapping json config
// params : dev or prod or local
func Init(env string) {
	logrus.Info("Config in ", env, " mode")
	logrus.SetLevel(logrus.DebugLevel)

	//read configuration file
	var config config
	if _, err := toml.DecodeFile("./config/config."+env+".toml", &config); err != nil {
		logrus.Fatal(err)
	}

	App = config.App
	App.Env = env
	Databases = config.Databases

	if env == "prod" || env == "dev" {
		//write log to file
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
		logrus.SetOutput(&lumberjack.Logger{
			Filename:   App.LogPath + "/skripsi-sso.log",
			MaxSize:    500, // megabytes
			MaxBackups: 10,
			MaxAge:     30,   //days
			Compress:   true, // disabled by default
		})

		utils.LogSummary.Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
		utils.LogSummary.Logger.SetOutput(&lumberjack.Logger{
			Filename:   App.LogPath + "/summary.log",
			MaxSize:    500, // megabytes
			MaxBackups: 10,
			MaxAge:     30,   //days
			Compress:   true, // disabled by default
		})
	}
}
