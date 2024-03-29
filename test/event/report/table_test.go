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

package report

import (
	"encoding/json"
	"fmt"
	utils "github.com/th2-net/th2-common-utils-go/pkg/event/report"
	"math/rand"
	"reflect"
	"regexp"
	"testing"
	"time"
)

func TestGetNewTable(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	amount := 5 + rand.Intn(10)
	var headers []string
	for i := 0; i < amount; i++ {
		headers = append(headers, fmt.Sprintf("header %v", i))
	}
	table := utils.GetNewTable(headers...)

	if table.Type != "table" {
		t.Errorf("Incorrect type for table: expected 'table' got '%v'", table.Type)
	}
	if amountHeaders := len(table.Headers); amountHeaders != amount {
		t.Errorf("Incorrect amount of headers: expected %v got %v", amount, amountHeaders)
	}
	if !reflect.DeepEqual(table.Headers, headers) {
		t.Errorf("Headers are not correct: expected %v got %v", headers, table.Headers)
	}
	if _, err := json.Marshal(&table); err != nil {
		t.Error("Error occurred during json encoding of the table")
	}
}

func TestAddRow(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	amount := 5 + rand.Intn(10)
	var headers []string
	for i := 0; i < amount; i++ {
		headers = append(headers, fmt.Sprintf("header %v", i))
	}
	table := utils.GetNewTable(headers...)

	expectedRow := make(map[string]string)
	for _, arg := range table.Headers {
		expectedRow[arg] = arg
	}
	for i := 0; i < amount; i++ {
		table.AddRow(headers...)
	}
	if amountRows := len(table.Rows); amountRows != amount {
		t.Errorf("Incorrect amount of rows: expected %v got %v", amount, amountRows)
	}
	for _, row := range table.Rows {
		if !reflect.DeepEqual(row, expectedRow) {
			t.Errorf("One of the rows is incorrectly added to the table; expected %v got %v", expectedRow, row)
		}
	}
	encoded, err := json.Marshal(&table)
	if err != nil {
		t.Error("Error occurred during json encoding of the table")
	}
	pattern := "^{\"type\":\"table\",\"rows\":\\[({(\"header \\d+\":\"header \\d+\",{0,1})*},{0,1})*\\],\"headers\":\\[(\"header \\d+\",{0,1})*\\]}$"
	match, _ := regexp.MatchString(pattern, string(encoded))
	if !match {
		t.Errorf("Incorrect string format: got %v", string(encoded))
	}
}
