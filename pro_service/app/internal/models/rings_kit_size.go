package models

type RingsKitSize struct {
	Id        string  `db:"id"`
	Outer     float64 `db:"outer"`
	Inner     float64 `db:"inner"`
	Thickness float64 `db:"thickness"`
}
