package main

import "os"

type config struct {
	mysqlUser string
	mysqlPass string
	mysqlHost string
	mysqlPort string

	redisUser string
	redisPass string
	redisHost string
}

func getConfig() config {
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

	redisUser := os.Getenv("REDIS_USER") //Redis Username and Password are intentionally blank on
	redisPass := os.Getenv("REDIS_PASS")
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	return config{
		mysqlUser: mysqlUser,
		mysqlPass: mysqlPass,
		mysqlHost: mysqlHost,
		mysqlPort: mysqlPort,

		redisUser: redisUser,
		redisPass: redisPass,
		redisHost: redisHost,
	}
}
