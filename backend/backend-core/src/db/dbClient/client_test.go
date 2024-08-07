package dbClient

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
)

func setupStubDatabaseConnection(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *relationalDatabaseClientImpl) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error occurred while opening a stub database connection: %s", err.Error())
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), new(gorm.Config))
	if err != nil {
		t.Fatalf("an error occurred while opening a GORM database session: %s", err.Error())
	}
	session := new(gorm.Session)
	session.Logger = logger.Default.LogMode(logger.Silent)
	gormDB = gormDB.Session(session)
	return db, mock, &relationalDatabaseClientImpl{db: gormDB}
}

func TestPersistSDType(t *testing.T) {
	db, mock, client := setupStubDatabaseConnection(t)
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	mock.ExpectBegin()
	mock.ExpectQuery(`^INSERT INTO "sd_types"`).
		WithArgs("shelly1pro").
		WillReturnRows(sqlmock.NewRows(sharedUtils.SliceOf("id")).AddRow(1))
	mock.ExpectQuery(`^INSERT INTO "sd_parameters"`).
		WithArgs(1, "relay_0_temperature", "number", 1, "relay_0_output", "boolean").
		WillReturnRows(sqlmock.NewRows(sharedUtils.SliceOf("id")).AddRow(1).AddRow(2))
	mock.ExpectCommit()
	sdType := dllModel.SDType{
		ID:         sharedUtils.NewEmptyOptional[uint32](),
		Denotation: "shelly1pro",
		Parameters: []dllModel.SDParameter{{
			ID:         sharedUtils.NewEmptyOptional[uint32](),
			Denotation: "relay_0_temperature",
			Type:       dllModel.SDParameterTypeNumber,
		}, {
			ID:         sharedUtils.NewEmptyOptional[uint32](),
			Denotation: "relay_0_output",
			Type:       dllModel.SDParameterTypeTypeBoolean,
		}},
	}
	if sdTypePersistResult := client.PersistSDType(sdType); sdTypePersistResult.IsFailure() {
		t.Fatalf("an error occurred while trying to persist the SD type: %s", sdTypePersistResult.GetError().Error())
	}
	if smUnfulfilledExpectationsError := mock.ExpectationsWereMet(); smUnfulfilledExpectationsError != nil {
		t.Fatalf("sqlmock: not all expectations were fulfilled: %s", smUnfulfilledExpectationsError.Error())
	}
}
