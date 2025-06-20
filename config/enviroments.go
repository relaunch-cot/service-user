package config

import "os"

var (
	PORT = os.Getenv("PORT")

	/////////////////////////////////////////// MYSQL INSTANCE
	MYSQL_HOST   = os.Getenv("MYSQL_HOST")
	MYSQL_PORT   = os.Getenv("MYSQL_PORT")
	MYSQL_USER   = os.Getenv("MYSQL_USER")
	MYSQL_PASS   = os.Getenv("MYSQL_PASS")
	MYSQL_DBNAME = os.Getenv("MYSQL_DBNAME")
)
