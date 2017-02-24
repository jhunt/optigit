package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jhunt/go-db"
	_ "github.com/mattn/go-sqlite3"
)

func split(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}

func database() (db.DB, error) {
	env := os.Getenv("DATABASE")
	if env == "" {
		return db.DB{}, fmt.Errorf("no DATABASE env var set; which database do you want to use?")
	}

	dsn := strings.Split(env, ":")
	if len(dsn) != 2 {
		return db.DB{}, fmt.Errorf("failed to determine database from DATABASE '%s' env var", os.Getenv("DATABASE"))
	}
	d := db.DB{
		Driver: dsn[0],
		DSN:    dsn[1],
	}

	err := d.Connect()
	if err != nil {
		return d, err
	}
	if !d.Connected() {
		return d, fmt.Errorf("not connected")
	}
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
