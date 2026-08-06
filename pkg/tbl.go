package pkg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gravestench/bitstream"
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
	if len(fileData) < headerBytes {
		return nil, fmt.Errorf("TBL header is truncated: got %d bytes", len(fileData))
	}
	lookupTable := make(TextTable)

	stream := bitstream.NewReader().FromBytes(fileData...)

	// skip past the CRC
	_, err := stream.Next(crcByteCount).Bytes().AsBytes()
	if err != nil {
		return nil, err
	}

	numberOfElements, err := stream.Next(2).Bytes().AsUInt16()
	if err != nil {
		return nil, err
	}

	hashTableSize, err := stream.Next(4).Bytes().AsUInt32()
	if err != nil {
		return nil, err
	}

	// Version (always 0)
	version, err := stream.Next(1).Bytes().AsByte()
	if err != nil {
		return nil, fmt.Errorf("reading TBL version: %w", err)
	}
	if version != 0 {
		return nil, fmt.Errorf("unsupported TBL version %d", version)
	}

	stream.Next(4).Bytes() // StringOffset
	stream.Next(4).Bytes() // When the number of times you have missed a match with a hash key equals this value, you give up because it is not there.
	stream.Next(4).Bytes() // FileSize

	remaining := len(fileData) - headerBytes
	if uint64(numberOfElements)*2+uint64(hashTableSize)*hashEntryBytes > uint64(remaining) {
		return nil, fmt.Errorf("TBL index tables exceed payload")
	}
	elementIndex := make([]uint16, numberOfElements)
	for i := 0; i < int(numberOfElements); i++ {
		elementIndex[i], err = stream.Next(2).Bytes().AsUInt16()
		if err != nil {
			return nil, fmt.Errorf("reading element index %d: %w", i, err)
		}
	}

	hashEntries := make([]hashEntry, hashTableSize)
	for i := 0; i < int(hashTableSize); i++ {
		td := hashEntry{}

		td.IsActive, _ = stream.Next(1).Bytes().AsBool()
		td.Index, _ = stream.Next(2).Bytes().AsUInt16()
		td.HashValue, _ = stream.Next(4).Bytes().AsUInt32()
		td.IndexString, _ = stream.Next(4).Bytes().AsUInt32()
		td.NameString, _ = stream.Next(4).Bytes().AsUInt32()
		td.NameLength, err = stream.Next(2).Bytes().AsUInt16()

		if err != nil {
			return nil, err
		}

		hashEntries[i] = td
	}

	for idx, hashEntry := range hashEntries {
		if !hashEntry.IsActive {
			continue
		}
		if hashEntry.NameLength == 0 || uint64(hashEntry.NameString)+uint64(hashEntry.NameLength) > uint64(len(fileData)) ||
			uint64(hashEntry.IndexString) >= uint64(len(fileData)) {
			return nil, fmt.Errorf("hash entry %d points outside TBL payload", idx)
		}

		stream.SetPosition(int(hashEntry.NameString))
		nameVal, err := stream.Next(int(hashEntry.NameLength - 1)).Bytes().AsBytes()
		if err != nil {
			return nil, err
		}
		value := string(nameVal)

		stream.SetPosition(int(hashEntry.IndexString))

		var key strings.Builder

		for {
			b, err := stream.Next(1).Bytes().AsByte()
			if err != nil {
				return nil, err
			}
			if b == 0 {
				break
			}

			key.WriteByte(b)
		}
		keyString := key.String()

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
