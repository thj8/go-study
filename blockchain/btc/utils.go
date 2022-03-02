package btc

import (
	"bytes"
	"encoding/binary"
	"log"
)

func IntToHex(val int64) []byte {
	buff := new(bytes.Buffer)

	err := binary.Write(buff, binary.BigEndian, val)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
