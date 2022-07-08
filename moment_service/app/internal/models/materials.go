package models

/*

	Mark Table
		id sereal
		title text


	Сделаю три таблицы. Запрос буду делать с джоином или чем-то подобным (вроде должно сработать)
	надо еще сортировать по температуре

	^ - at

	может болты тоже отдельно в таблицу вынести (можно будет сразу забирать диаметр)
*/

//? //TODO Хз правильно это или нет
type Materials struct {
	// возможно нужно сделать *float32 или разделить структуру на 3 части
	MarkId     string  `db:"mark_id"`
	Temp       float64 `db:"temp"`
	Voltage    float64 `db:"voltage"`
	Elasticity float64 `db:"elasticity"`
	Alpha      float64 `db:"alpha"`
}

type MaterialsResult struct {
	AlphaF      float64
	EpsilonAt20 float64
	Epsilon     float64
	SigmaAt20   float64
	Sigma       float64
}
