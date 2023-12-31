package service

import (
	"encoding/csv"
	"github.com/signintech/gopdf"
	"go.uber.org/zap"
	"goTSV/config"
	"goTSV/internal/constants"
	"goTSV/internal/domains"
	"goTSV/internal/shema"
	"goTSV/internal/watcher"
	"io"
	"os"
	"strings"
)

type Service struct {
	storage domains.Storage
	watcher *watcher.Watcher
	config  config.Config
	logger  *zap.Logger
}

func NewService(storage domains.Storage, watcher *watcher.Watcher, config config.Config) *Service {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil
	}
	return &Service{storage: storage, watcher: watcher, config: config, logger: logger}
}

func (s *Service) Scanner() error {
	out := make(chan string)
	go s.watcher.Scan(out)
	for file := range out {
		tsv, unitGuid, err := s.ParseFile(file)
		if err != nil {
			f := shema.Files{
				File: file,
				Err:  err.Error(),
			}
			err = s.storage.SaveFiles(f)
			if err != nil {
				s.logger.Info("failed to save file in db")
				return err
			}
			s.logger.Info("failed to parse")
			continue
		} else {
			f := shema.Files{
				File: file,
				Err:  "",
			}
			err := s.storage.SaveFiles(f)
			if err != nil {
				s.logger.Info("failed to save files in db")
				return err
			}
		}

		for _, ts := range tsv {
			err = s.storage.Save(ts)
			if err != nil {
				s.logger.Info("failed to save in db")
				return err
			}
		}
		err = s.WritePDF(tsv, unitGuid)
		if err != nil {
			s.logger.Info("failed to write pdf")
			return err
		}

	}
	return nil
}

func (s *Service) ParseFile(fileName string) ([]shema.Tsv, []string, error) {
	file, err := os.Open(s.config.DirectoryFrom + "/" + fileName)
	if err != nil {
		s.logger.Info("error opening")
		return nil, nil, err
	}

	if !strings.HasSuffix(file.Name(), ".tsv") {
		s.logger.Info("not a tsv file")
		return nil, nil, constants.ErrNotTSV
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	var data []shema.Tsv
	var array []string
	for {
		for _, d := range data {
			if array == nil {
				array = append(array, d.UnitGUID)
			}
		}
		str, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				return data, array, nil
			}
			s.logger.Info("not a tsv strings")
			return nil, nil, err
		}
		if str == nil {
			break
		}
		if len(strings.TrimSpace(str[3])) < 10 {
			continue
		}
	loop:
		for _, s := range data {
			for _, guid := range array {
				if s.UnitGUID == guid || guid == strings.TrimSpace(str[3]) {
					continue loop
				}
			}
			array = append(array, strings.TrimSpace(str[3]))

		}

		t := shema.Tsv{
			Number:       strings.TrimSpace(str[0]),
			MQTT:         strings.TrimSpace(str[1]),
			InventoryID:  strings.TrimSpace(str[2]),
			UnitGUID:     strings.TrimSpace(str[3]),
			MessageID:    strings.TrimSpace(str[4]),
			MessageText:  strings.TrimSpace(str[5]),
			Context:      strings.TrimSpace(str[6]),
			MessageClass: strings.TrimSpace(str[7]),
			Level:        strings.TrimSpace(str[8]),
			Area:         strings.TrimSpace(str[9]),
			Address:      strings.TrimSpace(str[10]),
			Block:        strings.TrimSpace(str[11]),
			Type:         strings.TrimSpace(str[12]),
			Bit:          strings.TrimSpace(str[13]),
			InvertBit:    strings.TrimSpace(str[14]),
		}
		data = append(data, t)
	}
	return data, array, nil
}

func (s *Service) WritePDF(tsv []shema.Tsv, unitGuid []string) error {
	for _, guid := range unitGuid {
		pdf := gopdf.GoPdf{}
		pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
		pdf.AddPage()

		defer pdf.Close()

		err := pdf.AddTTFFont("LiberationSerif-Regular", "resources/LiberationSerif-Regular.ttf")
		if err != nil {
			s.logger.Info("can't add font")
			return err
		}

		err = pdf.SetFont("LiberationSerif-Regular", "", 14)
		if err != nil {
			s.logger.Info("can't set font")
			return err
		}

		for _, t := range tsv {
			var resultArray []string
			if guid == t.UnitGUID {
				pdf.AddPage()
				resultArray = append(resultArray, "n: "+strings.TrimSpace(t.Number))
				resultArray = append(resultArray, "mqtt: "+strings.TrimSpace(t.MQTT))
				resultArray = append(resultArray, "invid: "+strings.TrimSpace(t.InventoryID))
				resultArray = append(resultArray, "unit_guid: "+strings.TrimSpace(t.UnitGUID))
				resultArray = append(resultArray, "msg_id: "+strings.TrimSpace(t.MessageID))
				resultArray = append(resultArray, "text: "+strings.TrimSpace(t.MessageText))
				resultArray = append(resultArray, "context: "+strings.TrimSpace(t.Context))
				resultArray = append(resultArray, "class: "+strings.TrimSpace(t.MessageClass))
				resultArray = append(resultArray, "level: "+strings.TrimSpace(t.Level))
				resultArray = append(resultArray, "area: "+strings.TrimSpace(t.Area))
				resultArray = append(resultArray, "addr: "+strings.TrimSpace(t.Address))
				resultArray = append(resultArray, "block: "+strings.TrimSpace(t.Block))
				resultArray = append(resultArray, "type: "+strings.TrimSpace(t.Type))
				resultArray = append(resultArray, "bit: "+strings.TrimSpace(t.Bit))
				resultArray = append(resultArray, "invert_bit: "+strings.TrimSpace(t.InvertBit))

				y := 20
				for _, str := range resultArray {
					pdf.SetXY(10, float64(y))
					err := pdf.Text(str)
					if err != nil {
						s.logger.Info("can't write string")
						return err
					}
					y += 20
				}
			}
		}
		resultFile := s.config.DirectoryTo + "/" + guid + ".pdf"
		err = pdf.WritePdf(resultFile)
		if err != nil {
			s.logger.Info("can't write pdf")
			return err
		}

	}
	return nil
}

func (s *Service) GetAll(r shema.Request) ([][]shema.Tsv, error) {
	if r.Limit <= 0 || r.UnitGUID == "" || r.Page < 0 {
		s.logger.Info("bad request in json")
		return nil, constants.ErrBadRequest
	}
	tsvFromDB, err := s.storage.GetAllGuids(r.UnitGUID)
	if err != nil {
		s.logger.Info("can't get tsv from db")
		return nil, constants.ErrNotFound
	}
	arrayWithPage := SubArray(r.Page, tsvFromDB)
	var resultArray [][]shema.Tsv
	for i := 0; i < len(arrayWithPage); i += r.Limit {
		end := i + r.Limit
		if end > len(arrayWithPage) {
			end = len(arrayWithPage)
		}
		resultArray = append(resultArray, arrayWithPage[i:end])
	}
	return resultArray, nil
}

func SubArray(startIndex int, data []shema.Tsv) []shema.Tsv {
	if startIndex < 0 || startIndex >= len(data) {
		return nil
	}
	return data[startIndex:]
}
