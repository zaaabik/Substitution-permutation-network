package main

import (
	"os"
	"github.com/jessevdk/go-flags"
	"crypto/rand"
	b "github.com/FeistelCipher/binary_reader"
	"fmt"
	"github.com/FeistelCipher/sp_net"
)

var opts struct {
	KeyPath string `long:"key" short:"k"`
	Mode    string `long:"mode" short:"m"`
	File    string `long:"file" short:"f"`
	OutFile string `long:"out" short:"o"`
	Len     int    `long:"len" short:"l"`
}

const (
	encrypt     = "e"
	decrypt     = "d"
	historgramm = "h"
	generateKey = "g"
	test        = "t"
)

const (
	roundCount = 40
)

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}

	if opts.Mode == encrypt {
		reader := b.Reader{}
		header, data, err := reader.ReadBmp(opts.File)
		if err != nil {
			fmt.Println("file error!")
		}

		s := sp_net.SPNet{}
		sBlocks := make([][]byte, 6)
		sBlocks[0], sBlocks[1], err = s.ReadBlock1("p1")
		sBlocks[2], sBlocks[3], err = s.ReadBlock1("p2")
		sBlocks[4], sBlocks[5], err = s.ReadBlock1("p3")

		k, err := reader.ReadKey("key")
		if err != nil {
			fmt.Println("cant open key!")
		}

		encryptedData, err := s.Encrypt(data, k, roundCount, sBlocks)
		if err != nil {
			fmt.Println("cant encrypt data")
		}

		reader.WriteBmp(opts.OutFile, header, encryptedData)
	} else if opts.Mode == decrypt {
		reader := b.Reader{}
		header, data, err := reader.ReadBmp(opts.File)
		if err != nil {
			fmt.Println("file error!")
		}

		s := sp_net.SPNet{}
		sBlocks := make([][]byte, 6)
		sBlocks[0], sBlocks[1], err = s.ReadBlock1("p1")
		sBlocks[2], sBlocks[3], err = s.ReadBlock1("p2")
		sBlocks[4], sBlocks[5], err = s.ReadBlock1("p3")

		k, err := reader.ReadKey("key")
		if err != nil {
			fmt.Println("cant open key!")
		}

		encryptedData, err := s.Decrypt(data, k, roundCount, sBlocks)
		if err != nil {
			fmt.Println("cant encrypt data")
		}

		reader.WriteBmp(opts.OutFile, header, encryptedData)
	} else if opts.Mode == historgramm {
		reader := b.Reader{}
		_, data, err := reader.ReadBmp(opts.File)
		if err != nil {
			fmt.Println("cant read file!")
			fmt.Println(err)
		}
		sp_net.MakeHist(opts.OutFile, data)
	} else if opts.Mode == generateKey {
		key := make([]byte, opts.Len)
		_, err := rand.Read(key)
		if err != nil {
			fmt.Println(err)
		}
		f, err := os.Create(opts.OutFile)
		defer f.Close()
		if err != nil {
			fmt.Println(err)
		}
		f.Write(key)
		f.Close()
	} else if opts.Mode == test {
		reader := b.Reader{}
		_, data, err := reader.ReadBmp(opts.File)
		if err != nil {
			fmt.Println("cant read file!")
			fmt.Println(err)
		}
		sp_net.Test(data, opts.OutFile)
	}
}
