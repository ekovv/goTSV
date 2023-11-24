package domains

import "goTSV/internal/shema"

//go:generate go run github.com/vektra/mockery/v3 --name=Storage
type Storage interface {
	Save(sh shema.Tsv) error
	SaveFiles(sh shema.Files) error
	GetAllGuids(unitGuid string) ([]shema.Tsv, error)
}
