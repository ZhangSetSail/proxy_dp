package public

import (
	"github.com/alecthomas/kingpin/v2"
)

type DBConfig struct {
	DBType   string
	DBHost   string
	DBPort   string
	DBUser   string
	DBPasswd string
	DBName   string
}

// AddDBFlags 注册数据库相关的命令行标志
func AddDBFlags(puc *DBConfig) {
	kingpin.Flag("db-type", "Database type (e.g., mysql, sqlite)").
		Default(puc.DBType).
		Envar("DB_TYPE").
		StringVar(&puc.DBType)

	kingpin.Flag("db-host", "Database host address").
		Default(puc.DBHost).
		Envar("DB_HOST").
		StringVar(&puc.DBHost)

	kingpin.Flag("db-port", "Database port").
		Default(puc.DBPort).
		Envar("DB_PORT").
		StringVar(&puc.DBPort)

	kingpin.Flag("db-user", "Database username").
		Default(puc.DBUser).
		Envar("DB_USER").
		StringVar(&puc.DBUser)

	kingpin.Flag("db-passwd", "Database password").
		Default(puc.DBPasswd).
		Envar("DB_PASSWD").
		StringVar(&puc.DBPasswd)

	kingpin.Flag("db-name", "Database name").
		Default(puc.DBName).
		Envar("DB_NAME").
		StringVar(&puc.DBName)
}
