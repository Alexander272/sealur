package models

type RingSize struct {
	Id    string  `db:"id"`
	Outer float64 `db:"outer"`
	Inner float64 `db:"inner"`
}
