package binary_reader

import (
	"os"
	"errors"
	s "github.com/FeistelCipher/sp_net"
	"encoding/binary"
	"fmt"
)

const bmpHeaderSize = 54

type Reader struct {
}

func (f Reader) ReadKey(path string) (s.Key, error) {
	file, err := os.Open(path)

	if err != nil {
		return s.Key{}, err
	}

	key := s.Key{}

	err = binary.Read(file, binary.BigEndian, &key)
	if err != nil{
		fmt.Print("cant read key")
		fmt.Print(err)
	}

	return key, nil
}
func (r Reader) ReadBmp(path string) (header, rgb []byte, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	stats, err := f.Stat()
	size := stats.Size()
	if bmpHeaderSize >= size {
		return nil, nil, errors.New("wrong file")
	}
	header = make([]byte, bmpHeaderSize)
	_, err = f.Read(header)
	rgb = make([]byte, size-bmpHeaderSize)
	_, err = f.Read(rgb)

	return header, rgb, err
}

func (r Reader) WriteBmp(path string, header, rgb []byte) (err error) {
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	}
	f.Write(header)
	f.Write(rgb)
	return nil
}
