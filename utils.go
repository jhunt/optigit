package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jhunt/go-db"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func split(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}

func vcapdb(vcapenv string) ([]string, error) {
	var vcap map[string][]struct {
		Credentials    map[string]interface{} `json:"credentials"`
		Label          string                 `json:"label"`
		Name           string                 `json:"name"`
		Plan           string                 `json:"plan"`
		Provider       interface{}            `json:"provider"`
		SyslogDrainURL interface{}            `json:"syslog_drain_url"`
		Tags           []string               `json:"tags"`
	}
	var dsn, driver string
	err := json.Unmarshal([]byte(vcapenv), &vcap)
	if err != nil {
		return nil, fmt.Errorf("error: '%s'\n", err)
	}
	for _, t := range vcap {
		for _, st := range t {
			dsn = fmt.Sprintf("%s", st.Credentials["uri"])
			dsnsplit := strings.Split(dsn, ":")
			driver = fmt.Sprintf("%s", dsnsplit[0])
		}
	}
	return []string{driver, dsn}, nil
}

func database() (db.DB, error) {
	env := os.Getenv("VCAP_SERVICES")
	var dsn []string
	var err error

	if env == "" {
		env = os.Getenv("DATABASE")
		if env == "" {
			return db.DB{}, fmt.Errorf("no DATABASE or VCAP_SERVICES env var set; which database do you want to use?")
		}

		dsn = strings.Split(env, ":")
		if len(dsn) != 2 {
			return db.DB{}, fmt.Errorf("failed to determine database from DATABASE '%s' env var", os.Getenv("DATABASE"))
		}
	} else {
		dsn, err = vcapdb(env)
		if err != nil {
			return db.DB{}, fmt.Errorf("could not parse VCAP_SERVICES: '%v'", err)
		}
	}

	d := db.DB{
		Driver: dsn[0],
		DSN:    dsn[1],
	}

	err = d.Connect()
	if err != nil {
		return d, fmt.Errorf("could not connect: '%v'", err)
	}
	if !d.Connected() {
		return d, fmt.Errorf("not connected")
	}
	Setup(d)
	return d, nil
}

func bindto() string {
	s := os.Getenv("BIND")
	if s != "" {
		return s
	}
	s = os.Getenv("PORT")
	if s != "" {
		return ":" + s
	}
	return ":3000"
}
