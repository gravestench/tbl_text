package pkg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// TextTable is a string map
type TextTable map[string]string

type hashEntry struct {
	IsActive    bool
	Index       uint16
	HashValue   uint32
	IndexString uint32
	NameString  uint32
	NameLength  uint16
}

const (
	crcByteCount   = 2
	headerBytes    = 21
	hashEntryBytes = 17
)

// Unmarshal the text dictionary from the given data
func Unmarshal(fileData []byte) (TextTable, error) {
	return UnmarshalReaderAt(bytes.NewReader(fileData), int64(len(fileData)))
}

// UnmarshalReaderAt decodes a TBL without copying the complete file. TBL hash
// entries contain absolute string offsets, so ReaderAt is the honest contract.
func UnmarshalReaderAt(source io.ReaderAt, size int64) (TextTable, error) {
	if source == nil {
		return nil, fmt.Errorf("TBL reader is nil")
	}
	if size < headerBytes {
		return nil, fmt.Errorf("TBL header is truncated: got %d bytes", size)
	}
	lookupTable := make(TextTable)
	header := make([]byte, headerBytes)
	if _, err := source.ReadAt(header, 0); err != nil {
		return nil, fmt.Errorf("reading TBL header: %w", err)
	}
	numberOfElements := binary.LittleEndian.Uint16(header[crcByteCount:])
	hashTableSize := binary.LittleEndian.Uint32(header[crcByteCount+2:])

	// Diablo II string tables use format version 1. Font TBL files have a
	// different header and are decoded by a separate codec.
	version := header[crcByteCount+2+4]
	if version != 1 {
		return nil, fmt.Errorf("unsupported TBL version %d", version)
	}

	remaining := size - headerBytes
	if uint64(numberOfElements)*2+uint64(hashTableSize)*hashEntryBytes > uint64(remaining) {
		return nil, fmt.Errorf("TBL index tables exceed payload")
	}
	offset := int64(headerBytes)
	elementBytes := make([]byte, int(numberOfElements)*2)
	if _, err := source.ReadAt(elementBytes, offset); err != nil {
		return nil, fmt.Errorf("reading element index: %w", err)
	}
	offset += int64(len(elementBytes))

	hashEntries := make([]hashEntry, hashTableSize)
	entryBytes := make([]byte, hashEntryBytes)
	for i := 0; i < int(hashTableSize); i++ {
		if _, err := source.ReadAt(entryBytes, offset+int64(i*hashEntryBytes)); err != nil {
			return nil, fmt.Errorf("reading hash entry %d: %w", i, err)
		}
		hashEntries[i] = hashEntry{
			IsActive:    entryBytes[0] != 0,
			Index:       binary.LittleEndian.Uint16(entryBytes[1:3]),
			HashValue:   binary.LittleEndian.Uint32(entryBytes[3:7]),
			IndexString: binary.LittleEndian.Uint32(entryBytes[7:11]),
			NameString:  binary.LittleEndian.Uint32(entryBytes[11:15]),
			NameLength:  binary.LittleEndian.Uint16(entryBytes[15:17]),
		}
	}

	for idx, hashEntry := range hashEntries {
		if !hashEntry.IsActive {
			continue
		}
		if hashEntry.NameLength == 0 || uint64(hashEntry.NameString)+uint64(hashEntry.NameLength) > uint64(size) ||
			uint64(hashEntry.IndexString) >= uint64(size) {
			return nil, fmt.Errorf("hash entry %d points outside TBL payload", idx)
		}

		nameVal := make([]byte, int(hashEntry.NameLength)-1)
		if _, err := source.ReadAt(nameVal, int64(hashEntry.NameString)); err != nil {
			return nil, fmt.Errorf("reading value for hash entry %d: %w", idx, err)
		}
		value := string(nameVal)
		key, err := readCStringAt(source, int64(hashEntry.IndexString), size)
		if err != nil {
			return nil, fmt.Errorf("reading key for hash entry %d: %w", idx, err)
		}
		keyString := key

		if keyString == "x" || keyString == "X" {
			keyString = "#" + strconv.Itoa(idx)
		}

		_, exists := lookupTable[keyString]
		if !exists {
			lookupTable[keyString] = value
		}
	}

	return lookupTable, nil
}

func readCStringAt(source io.ReaderAt, offset, size int64) (string, error) {
	var result strings.Builder
	var buffer [1]byte
	for offset < size {
		if _, err := source.ReadAt(buffer[:], offset); err != nil {
			return "", err
		}
		offset++
		if buffer[0] == 0 {
			return result.String(), nil
		}
		result.WriteByte(buffer[0])
	}
	return "", io.ErrUnexpectedEOF
}
