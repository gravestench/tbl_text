package tbl_text

import (
	"io"

	"github.com/gravestench/tbl_text/pkg"
)

type TextTable = pkg.TextTable

func Unmarshal(fileData []byte) (TextTable, error) {
	return pkg.Unmarshal(fileData)
}

func UnmarshalReaderAt(source io.ReaderAt, size int64) (TextTable, error) {
	return pkg.UnmarshalReaderAt(source, size)
}
