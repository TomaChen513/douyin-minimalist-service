package lib

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

var Cfg *ini.File

//服务端配置数据结构
type ServerConfig struct {
	RunMode         string
	HTTPPort        string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration

	Type            string
	User            string
	Password        string
	Host            string
	DbName          string
	TablePrefix     string

	RedisHost       string
	RedisIndex      string

}

//加载服务端配置
func LoadServerConfig() ServerConfig {

	var err error

	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatal(2, "Fail to parse 'conf/app.ini': %v", err)
	}


	//database配置节点读取
	database, err := Cfg.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}
	
	//redis 配置节点读取
	redis, err := Cfg.GetSection("redis")
	if err != nil {
		log.Fatal(2, "Fail to get section 'redis': %v", err)
	}


	Config := ServerConfig{
		RunMode:         Cfg.Section("").Key("RUN_MODE").MustString("debug"),
		Type:            database.Key("TYPE").MustString(""),
		User:            database.Key("USER").MustString(""),
		Password:        database.Key("PASSWORD").MustString(""),
		Host:            database.Key("HOST").MustString(""),
		DbName:          database.Key("NAME").MustString(""),
		TablePrefix:     database.Key("TABLE_PREFIX").MustString(""),
		RedisHost:       redis.Key("HOST").MustString(""),
		RedisIndex:      redis.Key("INDEX").MustString(""),
	}

	return Config
}
