package messaging

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MessageProducerPostgreImpl struct {
	db *gorm.DB
}

type postgreSchema struct {
	ID          uint `gorm:"primaryKey"`
	Payload     []byte
	Destination string
	Headers     map[string]string
	CreatedAt   int64 `gorm:"autoCreateTime"`
}

func (postgreSchema) TableName() string {
	return "messages"
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

func NewMessageProducerPostgreImpl(conf PostgreConfig) *MessageProducerPostgreImpl {
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

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&postgreSchema{})
	if err != nil {
		panic(err)
	}
	return &MessageProducerPostgreImpl{db}
}

func (p *MessageProducerPostgreImpl) Send(destination string, msg Message) error {
	// prepare message headers
	msg.SetHeader(DESTINATION, destination)
	msg.SetHeader(DATE, time.Now().UTC().Format(time.RFC1123))

	data := postgreSchema{
		Payload:     msg.Payload(),
		Destination: destination,
		Headers:     msg.Headers(),
	}

	result := p.db.Create(&data)
	if result.Error != nil {
		errors.Wrap(result.Error, "Couldn't insert message to db")
	}
	return nil
}