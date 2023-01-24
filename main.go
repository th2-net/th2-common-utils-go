package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	act "th2-grpc/th2_grpc_act_template"
	common_proto "th2-grpc/th2_grpc_common"
	"time"

	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
)

func MessageToMapConvertValue(value common_proto.Value) interface{} {

	/*
		It will be easier to simplify the real meaning of values
		to avoid many nested field names. That's what this method is doing,
		just converting values into the simplest shape.
	*/

	kind := value.GetKind()
	switch kind.(type) {

	case *common_proto.Value_SimpleValue:
		return value.GetSimpleValue()

	case *common_proto.Value_ListValue:
		var listValues []interface{}
		for _, val := range value.GetListValue().Values {
			listValues = append(listValues, MessageToMapConvertValue(*val))
		}
		return listValues

	case *common_proto.Value_MessageValue:
		fields := make(map[string]interface{})
		for field, field_val := range value.GetMessageValue().GetFields() {
			fields[field] = MessageToMapConvertValue(*field_val)
		}
		return fields

	default:
		return "error"

	}
}

func convertMessageToMap(message common_proto.Message) map[string]interface{} {

	/*
		The main field of request is message, so it will be quite easier to firstly
		write a converting function of message to map and then use it in the main function.
		Similarly to the above method, in some unnecessary fields the nested field names are neglected.
	*/

	metadata := message.Metadata

	fieldValues := make(map[string]interface{})

	for field, field_value := range message.Fields {
		fieldValues[field] = MessageToMapConvertValue(*field_value)
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

	/*
		With the aforementioned approach, we can easily convert request to map.
		So with this method we are done the first part.
	*/

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

func mapToMessageConvertValue(entity interface{}) *common_proto.Value {

	/*
		Now everything goes the same way but backwards. These method is inverse of the MessageToMapConvertValue method.
	*/

	switch entity.(type) {
	case string, int, float32:
		return &common_proto.Value{Kind: &common_proto.Value_SimpleValue{SimpleValue: entity.(string)}}

	case []interface{}:
		var listValue []*common_proto.Value
		for _, val := range entity.([]interface{}) {
			listValue = append(listValue, mapToMessageConvertValue(val))
		}
		return &common_proto.Value{Kind: &common_proto.Value_ListValue{ListValue: &common_proto.ListValue{Values: listValue}}}

	case map[string]interface{}:
		messageFields := make(map[string]*common_proto.Value)
		for k, v := range entity.(map[string]interface{}) {
			messageFields[k] = mapToMessageConvertValue(v)
		}
		return &common_proto.Value{Kind: &common_proto.Value_MessageValue{MessageValue: &common_proto.Message{Fields: messageFields}}}

	default:
		return &common_proto.Value{}
	}

}

func convertMapToRequest(m map[string]interface{}) act.PlaceMessageRequest {

	/*
		Simply decoding the map taking into consideration the way we have filled that.
	*/

	metadata := m["Message"].(map[string]interface{})["Metadata"]
	message := common_proto.Message{}

	if m["ParentEventId"] != nil {
		message.ParentEventId = &common_proto.EventID{Id: m["Message"].(map[string]interface{})["ParentEventId"].(string)}
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
		Timestamp:   metadata.(map[string]interface{})["Timestamp"].(*timestamp.Timestamp),
		MessageType: metadata.(map[string]interface{})["MessageType"].(string),
		//Properties:  metadata.(map[string]interface{})["MessageType"].(map[string]string),
		Protocol: metadata.(map[string]interface{})["Protocol"].(string),
	}

	fieldsConverted := make(map[string]*common_proto.Value)

	fields := m["Message"].(map[string]interface{})["Fields"]

	for field, field_value := range fields.(map[string]interface{}) {
		fieldsConverted[field] = mapToMessageConvertValue(field_value)
	}

	message.Fields = fieldsConverted

	request := act.PlaceMessageRequest{
		Message:       &message,
		ParentEventId: &common_proto.EventID{Id: m["ParentEventId"].(string)},
		Description:   m["Description"].(string),
		ConnectionId:  m["ConnectionId"].(*common_proto.ConnectionID),
	}
	return request
}

// random functions for certain fields in request object

func StringWithCharset(length int, charset string) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func generateRandomClordID(length int) string {
	const charset = "0123456789"
	return StringWithCharset(length, charset)
}

func genrateSecondaryRandomClordID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return StringWithCharset(length, charset)
}

func main() {

	// Creating request object to test the methods.

	eventID := common_proto.EventID{Id: uuid.New().String()}

	clordid := generateRandomClordID(7)
	secondaryClordid := genrateSecondaryRandomClordID(7)

	tradingPartyFields := map[string]*common_proto.Value{
		"NoPartyIDs": {Kind: &common_proto.Value_ListValue{ListValue: &common_proto.ListValue{
			Values: []*common_proto.Value{{Kind: &common_proto.Value_MessageValue{MessageValue: &common_proto.Message{
				Metadata: &common_proto.MessageMetadata{MessageType: "TradingParty_NoPartyIDs"},
				Fields: map[string]*common_proto.Value{
					"PartyID":       {Kind: &common_proto.Value_SimpleValue{SimpleValue: "Trader1"}},
					"PartyIDSource": {Kind: &common_proto.Value_SimpleValue{SimpleValue: "D"}},
					"PartyRole":     {Kind: &common_proto.Value_SimpleValue{SimpleValue: "76"}},
				},
			}}},
				{Kind: &common_proto.Value_MessageValue{MessageValue: &common_proto.Message{
					Metadata: &common_proto.MessageMetadata{MessageType: "TradingParty_NoPartyIDs"},
					Fields: map[string]*common_proto.Value{
						"PartyID":       {Kind: &common_proto.Value_SimpleValue{SimpleValue: "0"}},
						"PartyIDSource": {Kind: &common_proto.Value_SimpleValue{SimpleValue: "D"}},
						"PartyRole":     {Kind: &common_proto.Value_SimpleValue{SimpleValue: "3"}},
					},
				},
				}},
			}}},
		},
	}

	fields := map[string]*common_proto.Value{
		"Side":             {Kind: &common_proto.Value_SimpleValue{SimpleValue: "1"}},
		"SecurityID":       {Kind: &common_proto.Value_SimpleValue{SimpleValue: "INSTR1"}},
		"SecurityIDSource": {Kind: &common_proto.Value_SimpleValue{SimpleValue: "8"}},
		"OrdType":          {Kind: &common_proto.Value_SimpleValue{SimpleValue: "2"}},
		"AccountType":      {Kind: &common_proto.Value_SimpleValue{SimpleValue: "1"}},
		"OrderCapacity":    {Kind: &common_proto.Value_SimpleValue{SimpleValue: "A"}},
		"OrderQty":         {Kind: &common_proto.Value_SimpleValue{SimpleValue: "100"}},
		"Price":            {Kind: &common_proto.Value_SimpleValue{SimpleValue: "10"}},
		"ClOrdID":          {Kind: &common_proto.Value_SimpleValue{SimpleValue: clordid}},
		"SecondaryClOrdID": {Kind: &common_proto.Value_SimpleValue{SimpleValue: secondaryClordid}},
		"TransactTime":     {Kind: &common_proto.Value_SimpleValue{SimpleValue: time.Now().Format(time.RFC3339)}},
		"TradingParty":     {Kind: &common_proto.Value_MessageValue{MessageValue: &common_proto.Message{Fields: tradingPartyFields}}},
	}

	msg := common_proto.Message{Metadata: &common_proto.MessageMetadata{
		MessageType: "NewOrderSingle",
		Id:          &common_proto.MessageID{ConnectionId: &common_proto.ConnectionID{SessionAlias: "demo-conn1"}},
	},
		Fields:        fields,
		ParentEventId: &common_proto.EventID{Id: "randomString"},
	}

	conid := common_proto.ConnectionID{SessionAlias: "1", SessionGroup: "2"}

	request := act.PlaceMessageRequest{
		Message:       &msg,
		ParentEventId: &eventID,
		Description:   "User places an order.",
		ConnectionId:  &conid,
	}

	reqToMap := convertRequestToMap(request)
	//fmt.Println(reqToMap)

	requestAsJson, err := json.MarshalIndent(reqToMap, "", "   ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(requestAsJson)) // this thing and the thing which is printed on 298th line have to be the same, in this case, they are.

	fmt.Println()
	fmt.Println("First printing is done")
	fmt.Println()
	mapToReq := convertMapToRequest(reqToMap)

	mapToReqToMap := convertRequestToMap(mapToReq)

	requestAsJson1, err1 := json.MarshalIndent(mapToReqToMap, "", "   ")
	if err1 != nil {
		log.Fatal(err1)
	}

	fmt.Println(string(requestAsJson1))

}
