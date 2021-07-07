package sago

import (
	"fmt"
	"log"
	"strconv"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SagaInstanceRepositoryPostgreImpl struct {
	db *gorm.DB
}

type sagaInstancePgSchema struct {
	SagaID        uint   `gorm:"primaryKey"`
	SagaType      string `gorm:"index"`
	StateName     string
	LastRequestID string
	SagaData      []byte
	EndState      bool
	Compensating  bool
	CreatedAt     int64 `gorm:"autoCreateTime"`
	UpdatedAt     int64 `gorm:"autoUpdateTime"`
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

func NewSagaInstanceRepositoryPostgreImpl(conf PostgreConfig) *SagaInstanceRepositoryPostgreImpl {
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
	data := sagaInstancePgSchema{
		SagaType:      si.SagaType(),
		StateName:     si.StateName(),
		LastRequestID: si.LastRequestID(),
		SagaData:      si.SerializedSagaData(),
		EndState:      si.IsEndState(),
		Compensating:  si.IsCompensating(),
	}

	result := r.db.Create(&data)
	if result.Error != nil {
		return "", errors.Wrap(result.Error, "Couldn't store sagaInstance, type:"+si.SagaType())
	}

	log.Printf("saving SagaInstance id:%s type:%s", si.ID(), si.SagaType())

	return strconv.Itoa(int(data.SagaID)), nil
}

func (r *SagaInstanceRepositoryPostgreImpl) Find(sagaType, sagaID string) (*SagaInstance, error) {
	log.Printf("finding SagaInstance id:%s type:%s", sagaID, sagaType)

	var data sagaInstancePgSchema
	result := r.db.Where("saga_id = ? AND saga_type = ?", sagaID, sagaType).First(&data)

	if result.Error != nil {
		return nil, errors.Wrapf(
			result.Error,
			"Couldn't find SagaInstance type:%s, id:%s",
			sagaType, sagaID,
		)
	}

	si := NewSagaInstance(
		strconv.Itoa(int(data.SagaID)), data.SagaType,
		data.StateName, data.LastRequestID,
		data.SagaData, nil,
	)
	si.SetEndState(data.EndState)
	si.SetCompensating(data.Compensating)

	return si, nil
}

func (r *SagaInstanceRepositoryPostgreImpl) Update(si SagaInstance) error {
	id, err := strconv.Atoi(si.ID())
	if err != nil {
		return errors.Wrapf(
			err,
			"Couldn't convert SagaInstance type:%s id:%s to uint",
			si.SagaType(), si.ID(),
		)
	}

	data := sagaInstancePgSchema{
		SagaID:        uint(id),
		SagaType:      si.SagaType(),
		StateName:     si.StateName(),
		LastRequestID: si.LastRequestID(),
		SagaData:      si.SerializedSagaData(),
		EndState:      si.IsEndState(),
		Compensating:  si.IsCompensating(),
	}

	result := r.db.Model(&data).Updates(sagaInstancePgSchema{
		StateName:     data.StateName,
		LastRequestID: data.LastRequestID,
		SagaData:      data.SagaData,
		EndState:      data.EndState,
		Compensating:  data.Compensating,
	})

	if result.Error != nil || result.RowsAffected != 1 {
		return errors.Wrapf(
			result.Error,
			"Couldn't update SagaInstance type:%s, id:%s",
			data.SagaType, data.SagaID,
		)
	}

	return nil
}
