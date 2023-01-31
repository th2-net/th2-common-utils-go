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
