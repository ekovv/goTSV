package service

import (
	"encoding/csv"
	"fmt"
	"github.com/signintech/gopdf"
	"goTSV/config"
	"goTSV/internal/shema"
	"goTSV/internal/storage"
	"goTSV/internal/watcher"
	"io"
	"log"
	"net/http"
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
		err = s.WritePDF(tsv, unitGuid)
		if err != nil {
			return fmt.Errorf("failed to write pdf: %w", err)
		}

	}
	return nil
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
	for _, guid := range unitGuid {
		pdf := gopdf.GoPdf{}
		pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
		pdf.AddPage()

		fontUrl := "https://github.com/google/fonts/blob/main/ofl/actor/Actor-Regular.ttf"
		if err := DownloadFile("Actor-Regular.ttf", fontUrl); err != nil {
			return fmt.Errorf("don't download font: %w", err)
		}

		err := pdf.AddTTFFont("Actor-Regular", "Actor-Regular.ttf")
		if err != nil {
			return fmt.Errorf("don't add font: %w", err)
		}

		err = pdf.SetFont("Actor-Regular", "", 20)
		if err != nil {
			return fmt.Errorf("don't set font: %w", err)
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
			}
			for _, str := range resultArray {
				err = pdf.Text(str)
				if err != nil {
					return fmt.Errorf("can't write string: %w", err)
				}
			}
		}
		resultFile := s.config.DirectoryTo + guid + ".pdf"
		err = pdf.WritePdf(resultFile)
		if err != nil {
			return fmt.Errorf("can't write pdf: %w", err)
		}

	}
	return nil
}

func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
