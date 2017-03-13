package main

import (
	"fmt"

	"github.com/jhunt/go-db"
)

func SetupSchema(d db.DB) error {
	s := db.NewSchema()
	s.Version(1, func(d *db.DB) error {
		var err error
		switch d.Driver {
		case "mysql", "sqlite3":
			err = d.Exec(`CREATE TABLE repos (
			   id       integer      not null primary key,
			   org      varchar(100) not null,
			   name     varchar(200) not null,
			   included smallint     not null
			)`)
		case "postgres":
			err = d.Exec(`CREATE TABLE repos (
			   id       serial       not null primary key,
			   org      varchar(100) not null,
			   name     varchar(200) not null,
			   included smallint     not null
			)`)
		default:
			err = fmt.Errorf("unsupported database driver '%s'", d.Driver)
		}
		if err != nil {
			return fmt.Errorf("could not create repos table '%v'", err)
		}

		switch d.Driver {
		case "postgres", "mysql", "sqlite3":
			err = d.Exec(`CREATE TABLE pulls (
			   id          integer not null,
			   repo_id     integer not null,
			   created_at  integer NOT NULL,
			   updated_at  integer NOT NULL,
			   reporter    text NOT NULL,
			   assignees   text NOT NULL,
			   title       text NOT NULL,
			   primary key (id, repo_id)
			)`)
		default:
			err = fmt.Errorf("unsupported database driver '%s'", d.Driver)
		}
		if err != nil {
			return fmt.Errorf("could not create pulls table '%v'", err)
		}

		switch d.Driver {
		case "postgres", "mysql", "sqlite3":
			err = d.Exec(`CREATE TABLE issues (
			  id          integer not null,
			  repo_id     integer not null,
			  created_at  integer NOT NULL,
			  updated_at  integer NOT NULL,
			  reporter    text NOT NULL,
			  assignees   text NOT NULL,
			  title       text NOT NULL,
			  primary key (id, repo_id)
			)`)
		default:
			err = fmt.Errorf("unsupported database driver '%s'", d.Driver)
		}
		if err != nil {
			return fmt.Errorf("could not create issues table '%v'", err)
		}

		return nil
	})
	err := s.Migrate(&d, db.Latest)
	return err
}
