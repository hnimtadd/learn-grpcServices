package serializer

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func ProtoBufToJSON(message proto.Message) (string, error) {
	marshaller := protojson.MarshalOptions{
		EmitUnpopulated: true,
		Indent:          " ",
		UseProtoNames:   true,
		AllowPartial:    true,
		UseEnumNumbers:  false,
	}
	return marshaller.Format(message), nil
}
