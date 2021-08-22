package sago

import (
	"fmt"

	"git.coryptex.com/lib/sago/sagolog"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SagaInstanceRepositoryPostgreImpl struct {
	db *gorm.DB
}

type sagaInstancePgSchema struct {
	SagaID        string `gorm:"primaryKey"`
	SagaType      string `gorm:"primaryKey"`
	StateName     string
	LastRequestID string
	SagaData      []byte
	// EndState      bool
	// Compensating  bool
	CreatedAt int64 `gorm:"autoCreateTime"`
	UpdatedAt int64 `gorm:"autoUpdateTime"`
}

func (sagaInstancePgSchema) TableName() string {
	return "saga_instances"
}

type PostgreConfig struct {
	Host     string
	User     string
	Password string
	DB       string
	Port     string
	SSLMode  string
	TimeZone string
}

func NewSagaInstanceRepositoryPostgreImpl(conf PostgreConfig) SagaInstanceRepository {
	if conf.TimeZone == "" {
		conf.TimeZone = "GMT"
	}
	if conf.SSLMode == "" {
		conf.SSLMode = "disable"
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		conf.Host,
		conf.User,
		conf.Password,
		conf.DB,
		conf.Port,
		conf.SSLMode,
		conf.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&sagaInstancePgSchema{})
	if err != nil {
		panic(err)
	}
	return &SagaInstanceRepositoryPostgreImpl{db}
}

func (r *SagaInstanceRepositoryPostgreImpl) Save(si SagaInstance) (string, error) {
	const op string = "sago.instance_repository_postgre.Save"

	data := sagaInstancePgSchema{
		SagaID:        si.ID(),
		SagaType:      si.SagaType(),
		StateName:     si.StateName(),
		LastRequestID: si.LastRequestID(),
		SagaData:      si.SerializedSagaData(),
		// EndState:      si.IsEndState(),
		// Compensating:  si.IsCompensating(),
	}

	result := r.db.Create(&data)
	if result.Error != nil {
		return "", errors.Wrapf(
			result.Error,
			"failed to store sagaInstance of %s:%s saga\n",
			si.SagaType(), si.ID(),
		)
	}

	sagolog.Log(sagolog.DEBUG,
		fmt.Sprintf("%s: SagaInstance %s:%s was saved\n", op, si.SagaType(), data.SagaID),
	)

	return data.SagaID, nil
}

func (r *SagaInstanceRepositoryPostgreImpl) Find(sagaType, sagaID string) (*SagaInstance, error) {
	const op string = "sago.instance_repository_postgre.Save"

	sagolog.Log(sagolog.DEBUG,
		fmt.Sprintf("%s: finding SagaInstance %s:%s", op, sagaType, sagaID),
	)

	var data sagaInstancePgSchema
	result := r.db.Where("saga_id = ? AND saga_type = ?", sagaID, sagaType).First(&data)

	if result.Error != nil {
		return nil, errors.Wrapf(
			result.Error,
			"%s: failed to find SagaInstance %s:%s",
			op, sagaType, sagaID,
		)
	}

	si := NewSagaInstance(
		data.SagaID, data.SagaType,
		data.StateName, data.LastRequestID,
		data.SagaData, nil,
	)
	// si.SetEndState(data.EndState)
	// si.SetCompensating(data.Compensating)

	return si, nil
}

func (r *SagaInstanceRepositoryPostgreImpl) Update(si SagaInstance) error {
	if si.ID() == "" || si.SagaType() == "" {
		return errors.Errorf(
			"Saga id and type should not be nil, id:%s, type:%s\n",
			si.SagaType(), si.ID(),
		)
	}

	data := sagaInstancePgSchema{
		SagaID:        si.ID(),
		SagaType:      si.SagaType(),
		StateName:     si.StateName(),
		LastRequestID: si.LastRequestID(),
		SagaData:      si.SerializedSagaData(),
		// EndState:      si.IsEndState(),
		// Compensating:  si.IsCompensating(),
	}

	result := r.db.Model(&data).Updates(sagaInstancePgSchema{
		StateName:     data.StateName,
		LastRequestID: data.LastRequestID,
		SagaData:      data.SagaData,
		// EndState:      data.EndState,
		// Compensating:  data.Compensating,
	})

	if result.Error != nil || result.RowsAffected != 1 {
		return errors.Wrapf(
			result.Error,
			"failed to update SagaInstance %s:%s\n",
			data.SagaType, data.SagaID,
		)
	}

	return nil
}
