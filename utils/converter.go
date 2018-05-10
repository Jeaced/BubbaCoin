package utils

import (
	"bytes"
	"encoding/binary"
	"log"
	"encoding/gob"
)

// IntToHex converts an int64 to a byte array
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func SerializeMap(m map[int]int) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(m)
	if err != nil {
		panic(err)
	}

	return result.Bytes()
}

func DeserializeMap(bs []byte) map[int]int {
	var m map[int]int
	decoder := gob.NewDecoder(bytes.NewReader(bs))
	err := decoder.Decode(&m)
	if err != nil {
		panic(err)
	}

	return m
}

func SerializeInt(n int) []byte {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, uint64(n))

	return bs
}

func DeserializeInt(bs []byte) int {
	result := binary.LittleEndian.Uint64(bs)

	return int(result)
}
