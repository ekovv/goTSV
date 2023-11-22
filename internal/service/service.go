package service

import (
	"encoding/csv"
	"fmt"
	"github.com/signintech/gopdf"
	"goTSV/config"
	"goTSV/internal/shema"
	"goTSV/internal/storage"
	"goTSV/internal/watcher"
	"log"
	"os"
	"strings"
)

type Service struct {
	storage storage.DBStorage
	watcher *watcher.Watcher
	config  config.Config
}

func NewService(storage storage.DBStorage, watcher *watcher.Watcher, config config.Config) *Service {
	return &Service{storage: storage, watcher: watcher, config: config}
}

func (s *Service) Scanner() error {
	out := make(chan string)
	go s.watcher.Scan(out)
	for file := range out {
		tsv, unitGuid, err := s.ParseFile(file)
		if err != nil {
			f := shema.Files{
				File: file,
				Err:  err,
			}
			err = s.storage.SaveFiles(f)
			return fmt.Errorf("failed to parse: %w", err)
		} else {
			f := shema.Files{
				File: file,
				Err:  nil,
			}
			err := s.storage.SaveFiles(f)
			if err != nil {
				return fmt.Errorf("failed to save files in db: %w", err)
			}
		}

		for _, ts := range tsv {
			err = s.storage.Save(ts)
			if err != nil {
				return fmt.Errorf("failed to save in db: %w", err)
			}
		}
		err = s.WritePDF(tsv, unitGuid)
		if err != nil {
			return fmt.Errorf("failed to write pdf: %w", err)
		}

	}
	return nil
}

func (s *Service) ParseFile(fileName string) ([]shema.Tsv, []string, error) {
	file, err := os.Open(s.config.DirectoryFrom + "/" + fileName)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening: %w", err)
	}

	if !strings.HasSuffix(file.Name(), ".tsv") {
		return nil, nil, fmt.Errorf("not a tsv file: %s", fileName)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	var data []shema.Tsv
	var array []string
	for {
		str, err := reader.Read()
		if str == nil {
			break
		}
		if len(strings.TrimSpace(str[3])) < 10 {
			continue
		}
		if err != nil {
			log.Fatal(err)
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

		err := pdf.AddTTFFont("SANS-SERIF", "resources/Actor-Regular.ttf")
		if err != nil {
			return err
		}

		err = pdf.SetFont("SANS-SERIF", "", 14)
		if err != nil {
			return err
		}

		for _, t := range tsv {
			var resultArray []string
			if guid == t.UnitGUID {
				pdf.AddPage()
				resultArray = append(resultArray, "Number: "+strings.TrimSpace(t.Number))
				resultArray = append(resultArray, "MQTT: "+strings.TrimSpace(t.MQTT))
				resultArray = append(resultArray, "InventoryID: "+strings.TrimSpace(t.InventoryID))
				resultArray = append(resultArray, "UnitGUID: "+strings.TrimSpace(t.UnitGUID))
				resultArray = append(resultArray, "MessageID: "+strings.TrimSpace(t.MessageID))
				resultArray = append(resultArray, "MessageText: "+strings.TrimSpace(t.MessageText))
				resultArray = append(resultArray, "Context: "+strings.TrimSpace(t.Context))
				resultArray = append(resultArray, "MessageClass: "+strings.TrimSpace(t.MessageClass))
				resultArray = append(resultArray, "Level: "+strings.TrimSpace(t.Level))
				resultArray = append(resultArray, "Area: "+strings.TrimSpace(t.Area))
				resultArray = append(resultArray, "Address: "+strings.TrimSpace(t.Address))
				resultArray = append(resultArray, "Block: "+strings.TrimSpace(t.Block))
				resultArray = append(resultArray, "Type: "+strings.TrimSpace(t.Type))
				resultArray = append(resultArray, "Bit: "+strings.TrimSpace(t.Bit))
				resultArray = append(resultArray, "InvertBit: "+strings.TrimSpace(t.InvertBit))

				y := 20
				for _, str := range resultArray {
					pdf.SetXY(10, float64(y))
					err := pdf.Text(str)
					if err != nil {
						return fmt.Errorf("can't write string: %w", err)
					}
					y += 20
				}
			}
		}
		resultFile := s.config.DirectoryTo + "/" + guid + ".pdf"
		err = pdf.WritePdf(resultFile)
		if err != nil {
			return fmt.Errorf("can't write pdf: %w", err)
		}

	}
	return nil
}

func (s *Service) GetAll(req shema.Request) ([][]shema.Tsv, error) {
	tsvFromDB, err := s.storage.GetAllGuids(req.UnitGUID)
	if err != nil {
		return nil, fmt.Errorf("can't get tsvFromDB from db: %w", err)
	}
	arrayWithPage := SubArray(req.Page, tsvFromDB)
	var resultArray [][]shema.Tsv
	for i := 0; i < len(arrayWithPage); i += req.Limit {
		end := i + req.Limit
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
