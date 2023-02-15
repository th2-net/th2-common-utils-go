/*
 * Copyright 2023 Exactpro (Exactpro Systems Limited)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
