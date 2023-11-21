package shema

type Tsv struct {
	Number       string `tsv:"n"`
	MQTT         string `tsv:"mqtt"`
	InventoryID  string `tsv:"invid"`
	UnitGUID     string `tsv:"unit_guid" gorm:"not null"`
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

//CREATE TABLE checkedFiles (
//name VARCHAR(255) PRIMARY KEY,
//error VARCHAR(255)
//);
//
//CREATE TABLE Occurrence (
//ID           SERIAL PRIMARY KEY,
//Number       int,
//MQTT         VARCHAR(255),
//InventoryID  VARCHAR(255),
//UnitGUID     VARCHAR(255),
//MessageID    VARCHAR(255),
//MessageText  TEXT,
//Context      VARCHAR(255),
//MessageClass VARCHAR(255),
//Level        INTEGER,
//Area         VARCHAR(255),
//Address      VARCHAR(255),
//Block        BOOLEAN,
//Type         VARCHAR(255),
//Bit          INTEGER,
//InvertBit    INTEGER
//);
