package sequentialguid

import (
	"bytes"
	"encoding/binary"
	"log"
	"time"
	"unsafe"

	"github.com/google/uuid"
)

const INT_SIZE int = int(unsafe.Sizeof(0))

var Endianess binary.ByteOrder

// New creates a sequential guid that can be used as the primary key to solve the randomness of Globally Unique Identifier (GUID or UUID).
//
// By RFC 4122, the standard GUID is 16 bytes, the first 6 bytes is generated from UTC timestamp and the next 10 bytes are generated randomly,
// the Endianess are the computer architecture are also considered
//
// Follow up article link: https://www.codeproject.com/Articles/388157/GUIDs-as-fast-primary-keys-under-multiple-database.
// Golang implementation of code written in C#
func New() uuid.UUID {
	unixNano := time.Now().UTC().UnixMilli() / 10000
	randomBytes := []byte(uuid.New().String())

	buf := new(bytes.Buffer)

	if !isBigEndian() {
		Endianess = binary.LittleEndian
	} else {
		Endianess = binary.BigEndian
	}

	err := binary.Write(buf, Endianess, unixNano)
	if err != nil {
		log.Fatalln("binary.Write failed:", err)
	}

	timestampBytes := buf.Bytes()

	if !isBigEndian() {
		reverseBytes(timestampBytes)
	}

	//make 16 bytes slice
	guidBytes := make([]byte, 16)

	//copy from the third bytes from the timestamp to the first 6 bytes of the guid bytes
	copy(guidBytes, timestampBytes[2:8])

	//copy from the first 10 bytes from the random bytes to the last 10 bytes of the guid bytes
	copy(guidBytes[6:], randomBytes[:10])

	if !isBigEndian() {
		reverseBytesFromIndex(guidBytes, 0, 4)
		reverseBytesFromIndex(guidBytes, 4, 6)
	}

	newuuid, err := uuid.FromBytes(guidBytes)
	if err != nil {
		log.Fatalln(err)
	}

	return newuuid
}

// Get the endianess of the computer architecture
func getEndian() (ret bool) {
	var i int = 0x1
	bs := (*[INT_SIZE]byte)(unsafe.Pointer(&i))
	return bs[0] == 0
}

// Return true or false for the check for the endianess
func isBigEndian() bool {
	return getEndian()
}

// Reverse the whole byte from the end to the beginning
func reverseBytes(bytes []byte) []byte {
	for i := 0; i < len(bytes)/2; i++ {
		j := len(bytes) - i - 1
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	return bytes
}

// Reverse from the specific start index and the length
func reverseBytesFromIndex(bytes []byte, startIndex, length int) []byte {
	iterations := (length - startIndex) / 2

	for i := 0; i < iterations; i++ {
		bytes[startIndex], bytes[length-1] = bytes[length-1], bytes[startIndex]
		startIndex++
		length--
	}
	return bytes
}
