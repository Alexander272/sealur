package models

type OrderNew struct {
	Id     string `db:"id"`
	Date   string `db:"date"`
	Count  int64  `db:"count_position"`
	Number int64  `db:"number"`
	Status string `db:"status"`
	Info   string `db:"info"`
}

type ManagerOrder struct {
	Id      string `db:"id"`
	Date    string `db:"date"`
	UserId  string `db:"user_id"`
	Company string `db:"company"`
	Count   int64  `db:"count_position"`
	Number  int64  `db:"number"`
	Status  string `db:"status"`
	Info    string `db:"info"`
}

type OrderWithPosition struct {
	Id            string `db:"id"`
	Date          string `db:"date"`
	Count         int64  `db:"count_position"`
	Number        int64  `db:"number"`
	PositionId    string `db:"position_id"`
	Title         string `db:"title"`
	Amount        string `db:"amount"`
	PositionCount int64  `db:"position_count"`
	Info          string `db:"info"`
}

type OrderAnalytics struct {
	UserId           string `db:"user_id"`
	Company          string `db:"user_company"`
	ManagerId        string `db:"manager_id"`
	Manager          string `db:"manager"`
	OrderCount       int64  `db:"order_count"`
	PositionCount    int64  `db:"position_count"`
	PositionSnpCount int64  `db:"position_snp_count"`
}

type OrderCount struct {
	UserId              string `db:"user_id"`
	Company             string `db:"company"`
	Name                string `db:"name"`
	Inn                 string `db:"inn"`
	Count               int64  `db:"count"`
	OrderCount          int64  `db:"order_count"`
	SnpOrderCount       int64  `db:"order_snp_count"`
	PutgOrderCount      int64  `db:"order_putg_count"`
	PositionCount       int64  `db:"position_count"`
	PositionSnpCount    int64  `db:"position_snp_count"`
	PositionPutgCount   int64  `db:"position_putg_count"`
	AveragePosition     int64  `db:"average_position"`
	AverageSnpPosition  int64  `db:"average_snp_position"`
	AveragePutgPosition int64  `db:"average_putg_position"`
}

type FullOrderAnalytics struct {
	Id        string `db:"id"`
	UserId    string `db:"user_id"`
	ManagerId string `db:"manager_id"`
	Number    string `db:"number"`
	Date      string `db:"date"`
	Status    string `db:"status"`
	Company   string `db:"user_company"`
	Name      string `db:"user_name"`
	Manager   string `db:"manager"`
}
