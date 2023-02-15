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

package test

import (
	"math/rand"
	"regexp"
	"testing"
	"time"

	p_buff "th2-grpc/th2_grpc_common"

	utils "github.com/th2-net/th2-common-utils-go/th2_common_utils"
)

func TestCreateEventID(t *testing.T) {
	eventID := utils.CreateEventID()
	if eventID.Id == "" {
		t.Error("eventID.Id is empty")
	}
	match, _ := regexp.MatchString("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", eventID.Id)
	if !match {
		t.Error("eventID.Id is not a universally unique identifier (UUID)")
	}
}

func TestCreateEventBatch(t *testing.T) {
	var events []*p_buff.Event
	rand.Seed(time.Now().UnixNano())
	amount := 5 + rand.Intn(10)
	for i := 0; i < amount; i++ {
		newEvent := &p_buff.Event{
			Id: utils.CreateEventID(),
		}
		events = append(events, newEvent)
	}
	eventBatch := utils.CreateEventBatch(
		nil,
		events...,
	)
	if batchAmount := len(eventBatch.Events); batchAmount != amount {
		t.Errorf("Length test failed: expected %v got %v", amount, batchAmount)
	}
	for _, event := range eventBatch.Events {
		if event.Id.Id == "" {
			t.Error("eventID.Id is empty for one of the events in Batch")
		}
	}
}
