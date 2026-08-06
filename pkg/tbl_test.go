package pkg

import "testing"

func TestUnmarshalRejectsMalformedData(t *testing.T) {
	badVersion := make([]byte, headerBytes)
	badVersion[8] = 1
	for _, data := range [][]byte{{}, badVersion} {
		if _, err := Unmarshal(data); err == nil {
			t.Fatal("expected malformed TBL error")
		}
	}
}

func FuzzUnmarshal(f *testing.F) {
	f.Add([]byte{})
	f.Add(make([]byte, headerBytes))
	f.Fuzz(func(t *testing.T, data []byte) { _, _ = Unmarshal(data) })
}
