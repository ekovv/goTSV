package storage

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"goTSV/config"
	"goTSV/internal/shema"
)

type DBStorage struct {
	conn *sql.DB
}

func NewDBStorage(config config.Config) (*DBStorage, error) {
	db, err := sql.Open("postgres", config.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db %w", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate driver, %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"tsv", driver)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("failed to do migrate %w", err)
	}
	s := &DBStorage{
		conn: db,
	}

	return s, s.CheckConnection()
}

func (s *DBStorage) CheckConnection() error {
	if err := s.conn.Ping(); err != nil {
		return fmt.Errorf("failed to connect to db %w", err)
	}
	return nil
}

func (s *DBStorage) Save(sh shema.Tsv) error {
	insertQuery := `INSERT INTO occurrence(number, mqtt, inventoryid, unitguid, messageid, messagetext, context, messageclass, 
                level, area, address, block, type, bit, invertbit) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`
	_, err := s.conn.Exec(insertQuery, sh.Number, sh.MQTT, sh.InventoryID, sh.UnitGUID, sh.MessageID, sh.MessageText, sh.Context, sh.MessageClass, sh.Level,
		sh.Area, sh.Address, sh.Block, sh.Type, sh.Bit, sh.InvertBit)

	if err != nil {
		fmt.Errorf("not save in db %w", err)
	}
	return nil
}

func (s *DBStorage) SaveFiles(sh shema.Files) error {
	insertQuery := `INSERT INTO checkedfiles(name, error) VALUES ($1, $2)`
	_, err := s.conn.Exec(insertQuery, sh.File, sh.Err)
	if err != nil {
		fmt.Errorf("not save in db %w", err)
	}
	return nil
}

func (s *DBStorage) GetAllGuids(unitGuid string) ([]shema.Tsv, error) {
	query := "SELECT number, mqtt, inventoryid, unitguid, messageid, messagetext, context, messageclass, level, area, address, block, type, bit, invertbit FROM occurrence WHERE unitguid = $1"
	rows, err := s.conn.Query(query, unitGuid)
	if err != nil {
		return nil, fmt.Errorf("error getting: %w", err)
	}
	defer rows.Close()

	var data []shema.Tsv
	for rows.Next() {
		var d shema.Tsv
		err = rows.Scan(&d.Number, &d.MQTT, &d.InventoryID, &d.UnitGUID, &d.MessageID, &d.MessageText, &d.Context, &d.MessageClass,
			&d.Level, &d.Area, &d.Address, &d.Block, &d.Type, &d.Bit, &d.InvertBit)
		if err != nil {
			return nil, fmt.Errorf("error put in struct: %w", err)
		}
		data = append(data, d)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error rows: %w", err)
	}
	return data, nil
}
