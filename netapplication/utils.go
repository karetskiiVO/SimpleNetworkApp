package netapplication

import (
	"bytes"
	"encoding/binary"
	"syscall"
)

func recieveExact(conn int, n int) ([]byte, error) {
	data := make([]byte, n)
	recievedCnt := 0
	for recievedCnt < n {
		len, _, err := syscall.Recvfrom(conn, data[recievedCnt:min(n, batchSize)], 0)

		// conection err
		if err != nil {
			return nil, err
		}
		recievedCnt += len
	}

	return data, nil
}

func recieveFullMsg(conn int) ([]byte, error) {
	lengthBytes, err := recieveExact(conn, 4)
	if err != nil {
		return nil, err
	}
	var length int
	// могут быть проблемы
	binary.Read(bytes.NewReader(lengthBytes), binary.BigEndian, &length)

	msg, err := recieveExact(conn, length)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func sendMsg(conn int, msg []byte) error {
	buf := bytes.NewBuffer([]byte{})
	length := len(msg)
	binary.Write(buf, binary.BigEndian, length)

	err := syscall.Sendto(conn, buf.AvailableBuffer(), 0, nil)
	if err != nil {
		return err
	}

	buf.Reset()
	binary.Write(buf, binary.BigEndian, msg)

	for i := 0; i < length; i += batchSize {
		err := syscall.Sendto(conn, buf.AvailableBuffer()[i:min(length, i+batchSize)], 0, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
