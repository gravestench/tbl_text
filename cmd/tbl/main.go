package main

import (
	"encoding/json"
	"fmt"
	"github.com/gravestench/tbl_text"
	"os"
	"sort"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		usage()
	}
	cmd, path := os.Args[1], os.Args[2]
	data, e := os.ReadFile(path)
	if e != nil {
		fail(e)
	}
	t, e := tbl_text.Unmarshal(data)
	if e != nil {
		fail(e)
	}
	switch cmd {
	case "get":
		if len(os.Args) != 4 {
			usage()
		}
		v, ok := t[os.Args[3]]
		if !ok {
			os.Exit(1)
		}
		fmt.Println(v)
	case "search":
		if len(os.Args) != 4 {
			usage()
		}
		q := strings.ToLower(os.Args[3])
		keys := make([]string, 0)
		for k, v := range t {
			if strings.Contains(strings.ToLower(k), q) || strings.Contains(strings.ToLower(v), q) {
				keys = append(keys, k)
			}
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("%s\t%s\n", k, t[k])
		}
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if e := enc.Encode(t); e != nil {
			fail(e)
		}
	case "diff":
		if len(os.Args) != 4 {
			usage()
		}
		b, e := os.ReadFile(os.Args[3])
		if e != nil {
			fail(e)
		}
		other, e := tbl_text.Unmarshal(b)
		if e != nil {
			fail(e)
		}
		diff := map[string][2]string{}
		for k, v := range t {
			if other[k] != v {
				diff[k] = [2]string{v, other[k]}
			}
		}
		for k, v := range other {
			if _, ok := t[k]; !ok {
				diff[k] = [2]string{"", v}
			}
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if e := enc.Encode(diff); e != nil {
			fail(e)
		}
	default:
		usage()
	}
}
func usage()       { fmt.Fprintln(os.Stderr, "usage: tbl <get|search|json|diff> FILE [ARG]"); os.Exit(2) }
func fail(e error) { fmt.Fprintf(os.Stderr, "tbl: %v\n", e); os.Exit(1) }
