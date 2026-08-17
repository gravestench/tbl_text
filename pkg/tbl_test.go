package pkg

import (
	"bytes"
	"encoding/binary"
	"sync"
	"testing"
)

func TestUnmarshalRejectsMalformedData(t *testing.T) {
	badVersion := make([]byte, headerBytes)
	badVersion[8] = 2
	for _, data := range [][]byte{{}, badVersion} {
		if _, err := Unmarshal(data); err == nil {
			t.Fatal("expected malformed TBL error")
		}
	}
}

func validTableBytes() []byte {
	data := make([]byte, headerBytes+2+hashEntryBytes)
	binary.LittleEndian.PutUint16(data[2:4], 1)
	binary.LittleEndian.PutUint32(data[4:8], 1)
	data[8] = 1
	entry := data[headerBytes+2:]
	entry[0] = 1
	keyOffset := len(data)
	valueOffset := keyOffset + len("key") + 1
	binary.LittleEndian.PutUint32(entry[7:11], uint32(keyOffset))
	binary.LittleEndian.PutUint32(entry[11:15], uint32(valueOffset))
	binary.LittleEndian.PutUint16(entry[15:17], uint16(len("value")+1))
	return append(data, []byte("key\x00value\x00")...)
}

func TestUnmarshalReaderAtReadsOffsetStrings(t *testing.T) {
	data := validTableBytes()
	table, err := UnmarshalReaderAt(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatal(err)
	}
	if got := table["key"]; got != "value" {
		t.Fatalf("value = %q", got)
	}
}

func TestUnmarshalReaderAtSupportsConcurrentReaders(t *testing.T) {
	data := validTableBytes()
	source := bytes.NewReader(data)
	var group sync.WaitGroup
	for i := 0; i < 8; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			if table, err := UnmarshalReaderAt(source, int64(len(data))); err != nil || table["key"] != "value" {
				t.Errorf("table=%v err=%v", table, err)
			}
		}()
	}
	group.Wait()
}

func FuzzUnmarshal(f *testing.F) {
	f.Add([]byte{})
	f.Add(make([]byte, headerBytes))
	f.Fuzz(func(t *testing.T, data []byte) { _, _ = Unmarshal(data) })
}
