package main

import (
	"os"
	"github.com/jessevdk/go-flags"
	b "github.com/FeistelCipher/binary_reader"
	"fmt"
	"github.com/FeistelCipher/sp_net"
)

var opts struct {
	PBlocksPath string `long:"pBlocks" short:"p"`
	Mode        string `long:"mode" short:"m"`
	File        string `long:"file" short:"f"`
	OutFile     string `long:"out" short:"o"`
	Len         int    `long:"len" short:"l"`
	Count       int    `long:"count" short:"c"`
	SBlocksPath string `long:"sBlocks" short:"s"`
}

const (
	encrypt         = "e"
	decrypt         = "d"
	histogram       = "h"
	generatePBlocks = "p"
	test            = "t"
	generateSBlocks = "s"
	correlation     = "c"
)

const (
	roundCount = 1
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
		sBlocks, err := s.ReadSBlocks(opts.SBlocksPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		pBlocks, err := reader.ReadPBlocks(opts.PBlocksPath)
		if err != nil {
			fmt.Println("cant open key!")
			os.Exit(1)
		}

		encryptedData, err := s.Encrypt(data, pBlocks, roundCount, sBlocks)
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
			os.Exit(1)
		}

		s := sp_net.SPNet{}
		sBlocks, err := s.ReadSBlocks(opts.SBlocksPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		k, err := reader.ReadPBlocks(opts.PBlocksPath)
		if err != nil {
			fmt.Println("cant open key!")
			os.Exit(1)
		}

		encryptedData, err := s.Decrypt(data, k, roundCount, sBlocks)
		if err != nil {
			fmt.Println("cant encrypt data")
			os.Exit(1)
		}

		reader.WriteBmp(opts.OutFile, header, encryptedData)
	} else if opts.Mode == histogram {
		reader := b.Reader{}
		_, data, err := reader.ReadBmp(opts.File)
		if err != nil {
			fmt.Println("cant read file!")
			fmt.Println(err)
			os.Exit(1)
		}
		sp_net.MakeHist(opts.OutFile, data)
	} else if opts.Mode == generatePBlocks {
		s := sp_net.SPNet{}
		s.GeneratePBlocks(roundCount, opts.OutFile)
	} else if opts.Mode == test {
		reader := b.Reader{}
		_, data, err := reader.ReadBmp(opts.File)
		if err != nil {
			fmt.Println("cant read file!")
			fmt.Println(err)
			os.Exit(1)
		}
		sp_net.Test(data, opts.OutFile)
	} else if opts.Mode == generateSBlocks {
		s := sp_net.SPNet{}
		s.GenerateSBlock(opts.OutFile, roundCount)
	} else if opts.Mode == correlation {
		s := sp_net.SPNet{}
		r := b.Reader{}
		_, data1, _ := r.ReadBmp(opts.File)
		_, data2, _ := r.ReadBmp(opts.OutFile)
		len := opts.Len
		fmt.Println("corelation f1 -> f2 ", s.Correlation(data1, data2))
		fmt.Print("######################\n")
		autoCor1 := s.AutoCorrelation(data1, len)
		autoCor2 := s.AutoCorrelation(data2, len)
		for i := 0; i < len; i++ {
			fmt.Printf("autocorelation %s Δx = %d | %.5f\n", opts.File, i, autoCor1[i])
		}
		fmt.Print("######################\n")
		for i := 0; i < len; i++ {
			fmt.Printf("autocorelation %s Δx = %d | = %.5f\n", opts.OutFile, i, autoCor2[i])
		}
	}
}
