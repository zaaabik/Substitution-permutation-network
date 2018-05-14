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
	KeyPath    string `long:"key" short:"k"`
	Mode       string `long:"mode" short:"m"`
	File       string `long:"file" short:"f"`
	OutFile    string `long:"out" short:"o"`
	Len        int    `long:"len" short:"l"`
	Count      int    `long:"count" short:"c"`
	BlocksPath string `long:"blocks" short:"b"`
}

const (
	encrypt         = "e"
	decrypt         = "d"
	historgramm     = "h"
	generateKey     = "g"
	test            = "t"
	generateSblocks = "b"
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
			os.Exit(1)
		}

		s := sp_net.SPNet{}
		sBlocks, err := s.ReadBlock1(opts.BlocksPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		k, err := reader.ReadKey("key")
		if err != nil {
			fmt.Println("cant open key!")
			os.Exit(1)
		}

		encryptedData, err := s.Encrypt(data, k, roundCount, sBlocks)
		if err != nil {
			fmt.Println("cant encrypt data")
			os.Exit(1)
		}

		reader.WriteBmp(opts.OutFile, header, encryptedData)
	} else if opts.Mode == decrypt {
		reader := b.Reader{}
		header, data, err := reader.ReadBmp(opts.File)
		if err != nil {
			fmt.Println("file error!")
		}

		s := sp_net.SPNet{}
		sBlocks, err := s.ReadBlock1(opts.BlocksPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

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
	} else if opts.Mode == generateSblocks {
		s := sp_net.SPNet{}
		s.GenerateBlock(opts.OutFile, opts.Count)
	}
}
