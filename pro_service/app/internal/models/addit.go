package models

type Addit struct {
	Id          string `db:"id"`
	Materials   string `db:"materials"`
	Mod         string `db:"mod"`
	Temperature string `db:"temperature"`
	Mounting    string `db:"mounting"`
	Graphite    string `db:"graphite"`
	Fillers     string `db:"fillers"`
}
