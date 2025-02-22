/*
 * Copyright 2023-2025 Exactpro (Exactpro Systems Limited)
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

package event

import (
	p_buff "github.com/th2-net/th2-common-go/pkg/common/grpc/th2_grpc_common"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
)

func CreateEventID(book string, scope string) *p_buff.EventID {
	return &p_buff.EventID{
		BookName:       book,
		Scope:          scope,
		Id:             uuid.New().String(),
		StartTimestamp: timestamppb.Now(),
	}
}

func CreateEventBatch(ParentEventId *p_buff.EventID, Events ...*p_buff.Event) *p_buff.EventBatch {
	return &p_buff.EventBatch{
		ParentEventId: ParentEventId,
		Events:        append([]*p_buff.Event(nil), Events...),
	}
}
