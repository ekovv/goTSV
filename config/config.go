package config

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Host            string
	DirectoryFrom   string
	DirectoryTo     string
	DB              string
	RefreshInterval int
}

type F struct {
	host            *string
	directoryFrom   *string
	directoryTo     *string
	db              *string
	refreshInterval *int
}

var f F

const addr = "localhost:8080"

func init() {
	f.host = flag.String("a", addr, "-a=")
	f.directoryFrom = flag.String("f", "", "-f=from")
	f.db = flag.String("d", "", "-d=db")
	f.directoryTo = flag.String("t", "", "-t=to")
	f.refreshInterval = flag.Int("r", 10, "interval of check")
}

func New() (c Config) {
	flag.Parse()
	if envHost := os.Getenv("HOST"); envHost != "" {
		f.host = &envHost
	}
	if envRunDirectoryFrom := os.Getenv("DIRECTORY_FROM"); envRunDirectoryFrom != "" {
		f.directoryFrom = &envRunDirectoryFrom
	}
	if envRunDirectoryTo := os.Getenv("DIRECTORY_TO"); envRunDirectoryTo != "" {
		f.directoryTo = &envRunDirectoryTo
	}
	if envDB := os.Getenv("DATABASE_DSN"); envDB != "" {
		f.db = &envDB
	}
	envRefresh := os.Getenv("REFRESH_INTERVAL")
	if refreshInterval, _ := strconv.Atoi(envRefresh); refreshInterval != 0 {
		f.refreshInterval = &refreshInterval
	}
	c.Host = *f.host
	c.DirectoryFrom = *f.directoryFrom
	c.DB = *f.db
	c.DirectoryTo = *f.directoryTo
	c.RefreshInterval = *f.refreshInterval
	return c

}
