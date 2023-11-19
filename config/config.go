package config

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	DirectoryFrom   string
	DirectoryTo     string
	DB              string
	RefreshInterval int
}

type F struct {
	directoryFrom   *string
	directoryTo     *string
	db              *string
	refreshInterval *int
}

var f F

func init() {
	f.directoryFrom = flag.String("f", "", "-f=from")
	f.db = flag.String("d", "", "-d=db")
	f.directoryTo = flag.String("t", "", "-t=to")
	f.refreshInterval = flag.Int("r", 10, "interval of check")
}

func New() (c Config) {
	flag.Parse()
	envRefresh := os.Getenv("REPORT_INTERVAL")
	if envRunDirectoryFrom := os.Getenv("TOKEN"); envRunDirectoryFrom != "" {
		f.directoryFrom = &envRunDirectoryFrom
	}
	if envRunDirectoryTo := os.Getenv("TOKEN"); envRunDirectoryTo != "" {
		f.directoryTo = &envRunDirectoryTo
	}
	if envDB := os.Getenv("DATABASE_DSN"); envDB != "" {
		f.db = &envDB
	}
	if refreshInterval, _ := strconv.Atoi(envRefresh); refreshInterval != 0 {
		f.refreshInterval = &refreshInterval
	}
	c.DirectoryFrom = *f.directoryFrom
	c.DB = *f.db
	c.DirectoryTo = *f.directoryTo
	c.RefreshInterval = *f.refreshInterval
	return c

}
