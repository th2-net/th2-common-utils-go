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

type Table struct {
	Type    string        `json:"type"`
	Rows    []interface{} `json:"rows"`
	Headers []string      `json:"headers"`
}

func GetNewTable(headers ...string) *Table {
	return &Table{
		Type:    "table",
		Rows:    nil,
		Headers: headers,
	}
}
func (table *Table) AddRow(args ...string) {
	row := make(map[string]string)
	for i, arg := range args {
		row[table.Headers[i]] = arg
	}
	table.Rows = append(table.Rows, row)
}
