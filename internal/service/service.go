package service

import (
	"encoding/csv"
	"fmt"
	"github.com/signintech/gopdf"
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
}

func NewService(storage storage.DBStorage, watcher *watcher.Watcher) *Service {
	return &Service{storage: storage, watcher: watcher}
}

func (s *Service) Scanner() error {
	out := s.watcher.Scan()
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

	}
}

func (s *Service) ParseFile(fileName string) ([]shema.Tsv, []string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening: %w", err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	var data []shema.Tsv
	var array []string
	for {
		str, err := reader.Read()
		if err != nil {
			log.Fatal(err)
		}
		for _, s := range array {
			if s == strings.TrimSpace(str[3]) {
				break
			} else {
				array = append(array, strings.TrimSpace(str[3]))
			}
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
	for _, v := range unitGuid {
		pdf := gopdf.GoPdf{}
		pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
		pdf.AddPage()
		err := pdf.AddTTFFont("LiberationSerif-Regular", "resources/LiberationSerif-Regular.ttf")
		if err != nil {
			return fmt.Errorf("bad adding font: %w", err)
		}
		err = pdf.SetFont("LiberationSerif-Regular", "", 14)
		if err != nil {
			return fmt.Errorf("bad setting font: %w", err)
		}
		pdf.Cell(nil, "您好")
		pdf.WritePdf("hello.pdf")
	}
}
