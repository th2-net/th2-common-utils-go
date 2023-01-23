package main

import (
	/*"encoding/json"
	"fmt"
	"log"*/
	"math/rand"
	act "th2-grpc/th2_grpc_act_template"
	common_proto "th2-grpc/th2_grpc_common"
	"time"
	//"github.com/google/uuid"
)

func convertValue(value interface{}) interface{} {
	switch value.(type) {
	case common_proto.Value_SimpleValue:
		return value
	case common_proto.Value_ListValue:
		var values []common_proto.Value
		for val := range values {
			values = append(values, convertValue(val).(common_proto.Value))
		}
		return values
	case common_proto.Value_MessageValue:
		fields := make(map[string]interface{})
		for field, field_val := range value.(common_proto.Value_MessageValue).MessageValue.Fields {
			fields[field] = field_val
		}
		return fields
	default:
		return "error"
	} //I am searching for the method value.WhichOneOf("Kind") in go with that it will be better
}

func convertMessageToMap(message common_proto.Message) map[string]interface{} {
	metadata := message.Metadata

	fieldValues := make(map[string]interface{})

	for field, field_value := range message.Fields {
		fieldValues[field] = convertValue(*field_value)
	}

	m := map[string]interface{}{
		"Metadata": map[string]interface{}{
			"SessionAlias": metadata.Id.ConnectionId.SessionAlias,
			"SessionGroup": metadata.Id.ConnectionId.SessionGroup,
			"Direction":    metadata.Id.Direction,
			"Sequence":     metadata.Id.Sequence,
			"Subsequence":  metadata.Id.Subsequence,
			"Timestamp":    metadata.Timestamp,
			"MessageType":  metadata.MessageType,
			"Properties":   metadata.Properties,
			"Protocol":     metadata.Protocol,
		},
		"Fields": fieldValues,
	}

	if message.ParentEventId != nil {
		m["ParentEventId"] = message.ParentEventId.Id
	} else {
		m["ParentEventId"] = ""
	}

	return m

}

func convertRequestToMap(request act.PlaceMessageRequest) map[string]interface{} {
	m := map[string]interface{}{
		"Message":      convertMessageToMap(*request.Message),
		"Description":  request.Description,
		"ConnectionId": request.ConnectionId,
	}

	if request.ParentEventId != nil {
		m["ParentEventId"] = request.ParentEventId.Id
	} else {
		m["ParentEventId"] = nil
	}

	return m
}

func convertMapToRequest(m map[string]interface{}) act.PlaceMessageRequest {
	metadata := m["Message"].(map[string]interface{})["Metadata"]
	message := common_proto.Message{}

	if m["ParentEventId"] != nil {
		message.ParentEventId = &common_proto.EventID{Id: m["ParentEventId"].(string)}
	}

	message.Metadata = &common_proto.MessageMetadata{
		Id: &common_proto.MessageID{
			Sequence:    metadata.(map[string]interface{})["Sequence"].(int64),
			Subsequence: metadata.(map[string]interface{})["Subsequence"].([]uint32),
			Direction:   metadata.(map[string]interface{})["Direction"].(common_proto.Direction),
			ConnectionId: &common_proto.ConnectionID{
				SessionAlias: metadata.(map[string]interface{})["SessionAlias"].(string),
				SessionGroup: metadata.(map[string]interface{})["SessionGroup"].(string),
			},
		},
	}

	/*
		fieldsConverted := make(map[string]*common_proto.Value)

		fields := m["Fields"]

		for field, field_value := range fields.(map[string]interface{}) {
			fieldsConverted[field] = field_value.(*common_proto.Value)
		}

		message.Fields = fieldsConverted
	*/

	request := act.PlaceMessageRequest{
		Message:       &message,
		ParentEventId: &common_proto.EventID{Id: m["ParentEventId"].(string)},
		Description:   m["Description"].(string),
		ConnectionId:  m["ConnectionId"].(*common_proto.ConnectionID),
	}
	return request
}
