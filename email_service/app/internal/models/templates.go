package models

type JoinTemplate struct {
	Name     string
	Login    string
	Password string
	Services string
	Link     string
	Email    string
}

type OrderTemplate struct {
	Name     string
	Position string
	Company  string
	Address  string
	Email    string
	Phone    string
	Link     string
}

type ConfirmTemplate struct {
	Name         string
	Organization string
	Position     string
	Link         string
}

type ConfirmTemplateNew struct {
	Name  string
	Link  string
	Email string
}

type RejectTemplate struct {
	Name  string
	Email string
}

type BlockedTemplate struct {
	Ip    string
	Login string
}
