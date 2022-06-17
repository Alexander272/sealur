package models

type Interview struct {
	Organization  string   `json:"organization"`
	Name          string   `json:"name"`
	Email         string   `json:"email"`
	City          string   `json:"city"`
	Position      string   `json:"position"`
	Phone         string   `json:"phone"`
	Techprocess   string   `json:"techprocess"`
	Equipment     string   `json:"equipment"`
	Seal          string   `json:"seal"`
	Consumer      string   `json:"consumer"`
	Factory       string   `json:"factory"`
	Developer     string   `json:"developer"`
	Flange        string   `json:"flange"`
	TypeFl        string   `json:"typeFl"`
	Type          string   `json:"type"`
	DiffFrom      string   `json:"diffFrom"`
	DiffTo        string   `json:"diffTo"`
	PresWork      string   `json:"presWork"`
	PresTest      string   `json:"presTest"`
	Pressure      string   `json:"pressure"`
	Environ       string   `json:"environ"`
	TempWorkPipe  string   `json:"tempWorkPipe"`
	PresWorkPipe  string   `json:"presWorkPipe"`
	EnvironPipe   string   `json:"environPipe"`
	TempWorkAnn   string   `json:"tempWorkAnn"`
	PresWorkAnn   string   `json:"presWorkAnn"`
	EnvironAnn    string   `json:"environAnn"`
	Material      string   `json:"material"`
	BoltMaterial  string   `json:"boltMaterial"`
	Lubricant     bool     `json:"lubricant"`
	Along         string   `json:"along"`
	Across        string   `json:"across"`
	NonFlatness   string   `json:"nonFlatness"`
	Mounting      bool     `json:"mounting"`
	Condition     string   `json:"condition"`
	Period        string   `json:"period"`
	Abrasive      bool     `json:"abrasize"`
	Crystallized  bool     `json:"crystallized"`
	Penetrating   bool     `json:"penetrating"`
	DrawingNumber string   `json:"drawingNumber"`
	Info          string   `json:"info"`
	Drawing       *Drawing `json:"drawing"`
	Sizes         Size     `json:"size"`
}

type Drawing struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	OrigName string `json:"origName"`
	Group    string `json:"group"`
	Link     string `json:"link"`
}

type Size struct {
	Dy        string `json:"dy"`
	Py        string `json:"py"`
	DUp       string `json:"dUp"`
	D1        string `json:"d1"`
	D2        string `json:"d2"`
	D         string `json:"d"`
	H1        string `json:"h1"`
	H2        string `json:"h2"`
	Bolt      string `json:"bolt"`
	CountBolt int32  `json:"countBolt"`

	DIn  string `json:"dIn"`
	DOut string `json:"dOut"`
	H    string `json:"h"`
}
