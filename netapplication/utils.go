package netapplication

import (
	"bytes"
	"encoding/binary"
	"strings"
	"syscall"
)

const (
	batchSize = 1024
)

// Application interface
type Application interface {
	Close() error
	Run() error
}

func htons(num uint16) uint16 {
	return (num<<8)&0xff00 | num>>8
}

func tcpRecieveMsg(conn int) ([]byte, error) {
	recieveBytes := func(n int) ([]byte, error) {
		data := make([]byte, n)
		recievedCnt := 0
		for recievedCnt < n {
			len, _, err := syscall.Recvfrom(conn, data[recievedCnt:min(n, recievedCnt+batchSize)], 0)

			// conection err
			if err != nil {
				return nil, err
			}
			recievedCnt += len
		}

		return data, nil
	}

	lengthBytes, err := recieveBytes(4)
	if err != nil {
		return nil, err
	}

	var length int32
	binary.Read(bytes.NewReader(lengthBytes), binary.BigEndian, &length)

	msg, err := recieveBytes(int(length))
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func tcpSendMsg(conn int, msg []byte) error {
	wr := &strings.Builder{}

	msglen := int32(len(msg))
	binary.Write(wr, binary.BigEndian, msglen)
	binary.Write(wr, binary.BigEndian, msg)

	length := len(wr.String())

	for i := 0; i < length; i += batchSize {
		err := syscall.Sendto(conn, []byte(wr.String())[i:min(length, i+batchSize)], 0, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
