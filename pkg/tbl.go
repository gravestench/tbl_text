package pkg

import (
	"log"
	"strconv"

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
	crcByteCount = 2
)

// Unmarshal the text dictionary from the given data
func Unmarshal(fileData []byte) (TextTable, error) {
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
	if _, err = stream.Next(1).Bytes().AsByte(); err != nil {
		log.Fatal("Error reading Version record")
	}

	stream.Next(4).Bytes() // StringOffset
	stream.Next(4).Bytes() // When the number of times you have missed a match with a hash key equals this value, you give up because it is not there.
	stream.Next(4).Bytes() // FileSize

	elementIndex := make([]uint16, numberOfElements)
	for i := 0; i < int(numberOfElements); i++ {
		elementIndex[i], err = stream.Next(2).Bytes().AsUInt16()
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

		stream.SetPosition(int(hashEntry.NameString))
		nameVal, err := stream.Next(int(hashEntry.NameLength - 1)).Bytes().AsBytes()
		if err != nil {
			return nil, err
		}
		value := string(nameVal)

		stream.SetPosition(int(hashEntry.IndexString))

		key := ""

		for {
			b, err := stream.Next(1).Bytes().AsByte()
			if b == 0 {
				break
			}

			if err != nil {
				return nil, err
			}

			key += string(b)
		}

		if key == "x" || key == "X" {
			key = "#" + strconv.Itoa(idx)
		}

		_, exists := lookupTable[key]
		if !exists {
			lookupTable[key] = value
		}
	}

	return lookupTable, nil
}
