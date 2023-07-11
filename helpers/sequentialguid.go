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

	guidBytes := make([]byte, 16)

	copy(guidBytes, timestampBytes[2:8])
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

func getEndian() (ret bool) {
	var i int = 0x1
	bs := (*[INT_SIZE]byte)(unsafe.Pointer(&i))
	return bs[0] == 0
}

func isBigEndian() bool {
	return getEndian()
}

func reverseBytes(bytes []byte) []byte {
	for i := 0; i < len(bytes)/2; i++ {
		j := len(bytes) - i - 1
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	return bytes
}

func reverseBytesFromIndex(bytes []byte, startIndex, length int) []byte {
	iterations := (length - startIndex) / 2

	for i := 0; i < iterations; i++ {
		bytes[startIndex], bytes[length-1] = bytes[length-1], bytes[startIndex]
		startIndex++
		length--
	}
	return bytes
}
