package config

import "os"

type config struct {
	MysqlUser   string
	MysqlPass   string
	MysqlHost   string
	MysqlPort   string
	MysqlDriver string

	RedisUser string
	RedisPass string
	RedisHost string
	RedisPort string
}

func GetConfig() config {
	//Hmm, my naming conventions kinda suck
	mysqlDriver := os.Getenv("DB_DRIVER")
	if mysqlDriver == "" {
		mysqlDriver = "mysql"
	}
	mysqlUser := os.Getenv("MYSQL_USER")
	if mysqlUser == "" {
		mysqlUser = "root"
	}
	mysqlPass := os.Getenv("MYSQL_PASSWORD")
	if mysqlPass == "" {
		mysqlPass = "testing"
	}
	mysqlHost := os.Getenv("MYSQL_HOST")
	if mysqlHost == "" {
		mysqlHost = "localhost"
	}
	mysqlPort := os.Getenv("MYSQL_PORT")
	if mysqlPort == "" {
		mysqlPort = "3306"
	}

	//Redis Username and Password are intentionally blank on Localhost
	redisUser := os.Getenv("REDIS_USER")
	redisPass := os.Getenv("REDIS_PASS")
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	return config{
		MysqlUser:   mysqlUser,
		MysqlPass:   mysqlPass,
		MysqlHost:   mysqlHost,
		MysqlPort:   mysqlPort,
		MysqlDriver: mysqlDriver,

		RedisUser: redisUser,
		RedisPass: redisPass,
		RedisHost: redisHost,
		RedisPort: redisPort,
	}
}
