package binary_reader

import (
	"os"
	"errors"
	s "github.com/FeistelCipher/sp_net"
	"fmt"
)

const bmpHeaderSize = 54

type Reader struct {
}

func (r Reader) ReadPBlocks(path string) ([][]byte, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	stats, _ := file.Stat()
	size := stats.Size()
	buffer := make([]byte, size)
	file.Read(buffer)
	count := int(size / s.BlockSize)
	res := make([][]byte, count)
	for i := range res {
		res[i] = buffer[s.BlockSize*i:s.BlockSize*i+s.BlockSize]
	}

	if err != nil {
		fmt.Print("cant read pblock")
		fmt.Print(err)
	}

	return res, nil
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
