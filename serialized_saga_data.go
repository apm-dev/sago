package sago

import (
	"encoding/json"

	"google.golang.org/protobuf/proto"
)

// type SerializedSagaData struct {
// 	DataType  proto.Message
// 	DataBytes []byte
// }

func serializeSagaData(data proto.Message) ([]byte, error) {
	ser, err := json.Marshal(data)
	if err != nil {
		// TODO log
		return nil, err
	}
	return ser, nil
}

func deserializeSagaData(data []byte, dest proto.Message) error {
	err := json.Unmarshal(data, dest)
	if err != nil {
		// TODO log
		return err
	}
	return nil
}
