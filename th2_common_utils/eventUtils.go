package utils

import (
	p_buff "th2-grpc/th2_grpc_common"

	"github.com/google/uuid"
)

func CreateEventID() *p_buff.EventID {
	return &p_buff.EventID{Id: uuid.New().String()}
}

func CreateEventBatch(ParentEventId *p_buff.EventID, Events ...*p_buff.Event) *p_buff.EventBatch {
	EventBatch := p_buff.EventBatch{
		ParentEventId: ParentEventId,
	}
	EventBatch.Events = append(EventBatch.Events, Events...)
	return &EventBatch
}
