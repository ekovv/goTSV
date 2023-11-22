package shema

type Tsv struct {
	Number       string `tsv:"n"`
	MQTT         string `tsv:"mqtt"`
	InventoryID  string `tsv:"invid"`
	UnitGUID     string `tsv:"unit_guid"`
	MessageID    string `tsv:"msg_id"`
	MessageText  string `tsv:"text"`
	Context      string `tsv:"context"`
	MessageClass string `tsv:"class"`
	Level        string `tsv:"level"`
	Area         string `tsv:"area"`
	Address      string `tsv:"addr"`
	Block        string `tsv:"block"`
	Type         string `tsv:"type"`
	Bit          string `tsv:"bit"`
	InvertBit    string `tsv:"invert_bit"`
}

type Files struct {
	File string
	Err  error
}

type Request struct {
	UnitGUID string `json:"unit_guid"`
	Limit    int    `json:"limit"`
	Page     int    `json:"page"`
}
