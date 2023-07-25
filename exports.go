package tbl_text

import (
	"github.com/gravestench/tbl_text/pkg"
)

type TextTable = pkg.TextTable

func Unmarshal(fileData []byte) (TextTable, error) {
	return pkg.Unmarshal(fileData)
}
