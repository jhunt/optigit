package db

import (
	"fmt"
	"sort"
)

type Migrator func(*DB) error
type Schema struct {
	latest     int
	migrations map[int]Migrator
}

const Latest = -1

/*

  s := db.NewSchema()
  s.Version(1, v1schema)
  s.Version(2, v2schema)
  s.Version(3, func (db *DB) error {
    // implement database migration
    // logic here, inline
  })
  // etc ...

  err := s.Migrate(db, db.Latest)
  // ...

*/

func NewSchema() *Schema {
	return &Schema{
		latest:     0,
		migrations: make(map[int]Migrator),
	}
}

func (s *Schema) Latest() int {
	return s.latest
}

func (s *Schema) Version(number int, fn Migrator) {
	s.migrations[number] = fn
	if number > s.latest {
		s.latest = number
	}
}

func (s *Schema) Current(d *DB) (int, error) {
	r, err := d.Query(`SELECT version FROM schema_info LIMIT 1`)
	if err != nil {
		if err.Error() == "no such table: schema_info" {
			return 0, nil
		}
		if err.Error() == `pq: relation "schema_info" does not exist` {
			return 0, nil
		}
		return 0, err
	}
	defer r.Close()

	// no records = no schema
	if !r.Next() {
		return 0, nil
	}

	var v int
	err = r.Scan(&v)
	// failed unmarshall is an actual error
	if err != nil {
		return 0, err
	}

	// invalid (negative) schema version is an actual error
	if v < 0 {
		return 0, fmt.Errorf("Invalid schema version %d found", v)
	}

	return int(v), nil
}

func (s *Schema) IsAt(d *DB, want int) bool {
	have, err := s.Current(d)
	return err == nil && have == want
}

func (s *Schema) Migrate(d *DB, to int) error {
	if to == Latest {
		to = s.latest
	}

	at, err := s.Current(d)
	if err != nil {
		return err
	}

	if at > to {
		return fmt.Errorf("Schema (v%d) is newer than requested (v%d)", at, to)
	}

	var vv []int
	for k, _ := range s.migrations {
		vv = append(vv, k)
	}
	sort.Ints(vv)

	for _, v := range vv {
		if at < v {
			err = s.migrations[v](d)
			if err != nil {
				return err
			}
		}
	}

	// set up the schema_info table
	d.Exec(`CREATE TABLE schema_info (version INTEGER)`)
	d.Exec(`TRUNCATE TABLE schema_info`)
	d.Exec(`INSERT INTO schema_info VALUES ($1)`, to)

	return nil
}
