package service

import (
	"encoding/csv"
	"errors"
	"go.uber.org/zap"
	"goTSV/config"
	"goTSV/internal/constants"
	"goTSV/internal/domains/mocks"
	"goTSV/internal/shema"
	"os"
	"reflect"
	"testing"
)

func TestService_ParseFile(t *testing.T) {
	type args struct {
		file string
		dir  string
	}
	tests := []struct {
		name      string
		args      args
		wantTsv   []shema.Tsv
		wantGuids []string
		wantErr   error
	}{
		{
			name: "OK#1",
			args: args{dir: "testDirectory", file: "OK1.tsv"},
			wantTsv: []shema.Tsv{
				{
					Number:       "5",
					InventoryID:  "G-044325",
					UnitGUID:     "01749246-9617-585e-9e19-157ccad61ee2",
					MessageID:    "cold78_Defrost_status",
					MessageText:  "Разморозка",
					MessageClass: "waiting",
					Level:        "100",
					Area:         "LOCAL",
					Address:      "cold78_status.Defrost_status",
				},
			},
			wantGuids: []string{
				"01749246-9617-585e-9e19-157ccad61ee2",
			},
			wantErr: nil,
		},
		{
			name: "OK#2",
			args: args{dir: "testDirectory", file: "OK2.tsv"},
			wantTsv: []shema.Tsv{
				{
					Number:       "5",
					InventoryID:  "G-044325",
					UnitGUID:     "01749246-9617-585e-9e19-157ccad61ee2",
					MessageID:    "cold78_Defrost_status",
					MessageText:  "Разморозка",
					MessageClass: "waiting",
					Level:        "100",
					Area:         "LOCAL",
					Address:      "cold78_status.Defrost_status",
				},
				{
					Number:       "6",
					InventoryID:  "G-044325",
					UnitGUID:     "01749246-9617-585e-9e19-157ccad61ee2",
					MessageID:    "cold78_VentSK_status",
					MessageText:  "Вентилятор",
					MessageClass: "working",
					Level:        "100",
					Area:         "LOCAL",
					Address:      "cold78_status.VentSK_status",
				},
			},
			wantGuids: []string{
				"01749246-9617-585e-9e19-157ccad61ee2",
			},
			wantErr: nil,
		},
		{
			name:      "BAD#1",
			args:      args{dir: "testDirectory", file: "BAD1"},
			wantTsv:   nil,
			wantGuids: nil,
			wantErr:   constants.ErrNotTSV,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := createTempDir(tt.args.dir, t)
			if err != nil {
				t.Errorf("not creating temp dir: %v", err)
				return
			}

			defer func(dir string) {
				err := removeTempDir(dir)
				if err != nil {
					t.Errorf("not delete temp dir: %v", err)
					return
				}
			}(tempDir)

			file, err := createTempFile(tempDir, tt.args.file)
			if err != nil {
				t.Errorf("not creating temp file: %v", err)
				return
			}

			err = writeDataToFile(file, tt.wantTsv)
			if err != nil {
				t.Errorf("not writing temp file: %v", err)
				return
			}
			logger, err := zap.NewProduction()
			s := &Service{
				config: config.Config{DirectoryFrom: tempDir},
				logger: logger,
			}
			tsv, guids, err := s.ParseFile(tt.args.file)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ParseFile error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tsv, tt.wantTsv) {
				t.Errorf("Parse() got = %v, want %v", tsv, tt.wantTsv)
			}
			if !reflect.DeepEqual(guids, tt.wantGuids) {
				t.Errorf("Parse() got1 = %v, want %v", guids, tt.wantGuids)
			}
		})
	}
}

func createTempDir(dir string, t *testing.T) (string, error) {
	tempDir, err := os.MkdirTemp(".", dir)
	if err != nil {
		t.Errorf("not created directory")
		return "", err
	}
	return tempDir, nil
}

func createTempFile(dir, name string) (*os.File, error) {
	tempFile, err := os.Create(dir + "/" + name)
	if err != nil {
		return nil, err
	}
	return tempFile, nil
}

func writeDataToFile(file *os.File, data []shema.Tsv) error {
	writer := csv.NewWriter(file)
	writer.Comma = '\t'

	for _, d := range data {
		record := []string{d.Number, d.MQTT, d.InventoryID, d.UnitGUID, d.MessageID, d.MessageText, d.Context,
			d.MessageClass, d.Level, d.Area, d.Address, d.Block, d.Type, d.Bit, d.InvertBit}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func removeTempDir(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	return nil
}

type storageMock[A any] func(c *mocks.Storage, args A)

func TestService_GetAll(t *testing.T) {
	tests := []struct {
		name        string
		args        shema.Request
		storageMock storageMock[shema.Request]
		wantArr     [][]shema.Tsv
		wantErr     error
	}{
		{
			name: "OK1",
			args: shema.Request{
				UnitGUID: "01749246-9617-585e-9e19-157ccad61ee2",
				Limit:    1,
				Page:     1,
			},

			storageMock: func(c *mocks.Storage, r shema.Request) {
				c.Mock.On("GetAllGuids", r.UnitGUID).Return([]shema.Tsv{
					{
						Number:       "5",
						InventoryID:  "G-044325",
						UnitGUID:     "01749246-9617-585e-9e19-157ccad61ee2",
						MessageID:    "cold78_Defrost_status",
						MessageText:  "Разморозка",
						MessageClass: "waiting",
						Level:        "100",
						Area:         "LOCAL",
						Address:      "cold78_status.Defrost_status",
					},
					{
						Number:       "6",
						InventoryID:  "G-044325",
						UnitGUID:     "01749246-9617-585e-9e19-157ccad61ee2",
						MessageID:    "cold78_VentSK_status",
						MessageText:  "Вентилятор",
						MessageClass: "working",
						Level:        "100",
						Area:         "LOCAL",
						Address:      "cold78_status.VentSK_status",
					},
				}, nil).Times(1)
			},
			wantArr: [][]shema.Tsv{{
				{
					Number:       "6",
					InventoryID:  "G-044325",
					UnitGUID:     "01749246-9617-585e-9e19-157ccad61ee2",
					MessageID:    "cold78_VentSK_status",
					MessageText:  "Вентилятор",
					MessageClass: "working",
					Level:        "100",
					Area:         "LOCAL",
					Address:      "cold78_status.VentSK_status",
				}}},
			wantErr: nil,
		},
		{
			name: "OK2",
			args: shema.Request{
				UnitGUID: "01749246-9617-585e-9e19-157ccad61ee2",
				Limit:    1,
				Page:     2,
			},

			storageMock: func(c *mocks.Storage, r shema.Request) {
				c.Mock.On("GetAllGuids", r.UnitGUID).Return([]shema.Tsv{
					{
						Number:       "5",
						InventoryID:  "G-044325",
						UnitGUID:     "01749246-9617-585e-9e19-157ccad61ee2",
						MessageID:    "cold78_Defrost_status",
						MessageText:  "Разморозка",
						MessageClass: "waiting",
						Level:        "100",
						Area:         "LOCAL",
						Address:      "cold78_status.Defrost_status",
					},
					{
						Number:       "6",
						InventoryID:  "G-044325",
						UnitGUID:     "01749246-9617-585e-9e19-157ccad61ee2",
						MessageID:    "cold78_VentSK_status",
						MessageText:  "Вентилятор",
						MessageClass: "working",
						Level:        "100",
						Area:         "LOCAL",
						Address:      "cold78_status.VentSK_status",
					},
					{
						Number: "7",
					},
					{
						Number: "8",
					},
				}, nil).Times(1)
			},
			wantArr: [][]shema.Tsv{{
				{
					Number: "7",
				},
			}, {
				{
					Number: "8",
				},
			}},
			wantErr: nil,
		},
		{
			name: "BAD1",
			args: shema.Request{
				UnitGUID: "ahsyuiflgh",
				Limit:    1,
				Page:     1,
			},

			storageMock: func(c *mocks.Storage, r shema.Request) {
				c.Mock.On("GetAllGuids", r.UnitGUID).Return(nil, errors.New("error getting")).Times(1)
			},
			wantArr: nil,
			wantErr: constants.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := mocks.NewStorage(t)
			tt.storageMock(storage, tt.args)
			logger, err := zap.NewProduction()

			service := Service{
				storage: storage,
				logger:  logger,
			}
			tsvs, err := service.GetAll(tt.args)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tsvs, tt.wantArr) {
				t.Errorf("got %v, want %v", tsvs, tt.wantArr)
			}
		})
	}
}
