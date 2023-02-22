package models

import "github.com/lib/pq"

type Standard struct {
	Id     string         `db:"id"`
	Title  string         `db:"title"`
	Format pq.StringArray `db:"format"`
}
