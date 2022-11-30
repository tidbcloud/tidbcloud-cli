// Copyright 2022 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

const (
	JsonFormat  string = "json"
	HumanFormat string = "human"
)

type Column string
type Row []string

func PrintJson(out io.Writer, items interface{}) error {
	v, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(out, string(v))
	return nil
}

func PrintHumanTable(out io.Writer, columns []Column, rows []Row) error {
	headerFmt := color.New(color.FgCyan, color.Underline).SprintfFunc()

	c := make([]interface{}, len(columns))
	for i, col := range columns {
		c[i] = col
	}
	tbl := table.New(c...)
	tbl.WithHeaderFormatter(headerFmt)

	for _, row := range rows {
		r := make([]interface{}, len(row))
		for i, col := range row {
			r[i] = col
		}
		tbl.AddRow(r...)
	}

	fmt.Fprintln(out)
	tbl.Print()

	// for human format, we print the table with brief information.
	return nil
}
