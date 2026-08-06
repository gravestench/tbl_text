// Command tbl-to-json decodes a Diablo II text table to JSON.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gravestench/tbl_text"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s FILE.tbl\n", os.Args[0])
		os.Exit(2)
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fail(err)
	}
	table, err := tbl_text.Unmarshal(data)
	if err != nil {
		fail(err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(table); err != nil {
		fail(err)
	}
}

func fail(err error) { fmt.Fprintf(os.Stderr, "tbl-to-json: %v\n", err); os.Exit(1) }
