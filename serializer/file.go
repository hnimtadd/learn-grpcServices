package serializer

import (
	"fmt"
	"os"

	"google.golang.org/protobuf/proto"
)

func WriteProtbufToJSONFile(message proto.Message, filename string) error {
	data, err := ProtoBufToJSON(message)
	if err != nil {
		return fmt.Errorf("Cannot marshal proto message to JSON: %v", err)
	}
	err = os.WriteFile(filename, []byte(data), 0666)
	if err != nil {
		return fmt.Errorf("Cannot write binary data to file: %v", err)
	}
	return nil
}

func WriteProtobufToBinaryFile(message proto.Message, filename string) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, data, 0666)
	if err != nil {
		return fmt.Errorf("Cannot write binary data to file: %v", err)
	}
	return nil
}

func ReadProtobufFromBinaryFile(filename string, message proto.Message) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(data, message)
	if err != nil {
		return fmt.Errorf("Cannot unmashal biary to proto message: %v", err)
	}
	return nil
}
