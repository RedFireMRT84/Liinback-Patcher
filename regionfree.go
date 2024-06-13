package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"os"
)

const (
	offsetRegion = 0x19d
	offsetTMD    = 0xd00
	offsetBF     = 0x1c1
)

var regionStr = []string{"JAP/NTSC-J", "NTSC-U", "PAL", "*FREE*", "UNKNOWN"}

func hexdump(data []byte) {
	for i, b := range data {
		fmt.Printf("%02x ", b)
		if (i+1)%8 == 0 {
			fmt.Print(" ")
		}
		if (i+1)%16 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}

func truchaSignTicket(ticket []byte) bool {
	sha1Hash := sha1.New()
	var sha1Bytes [20]byte

	for i := 0; i < 65535; i++ {
		binary.LittleEndian.PutUint16(ticket[offsetBF:], uint16(i))
		sha1Hash.Reset()
		sha1Hash.Write(ticket[0x140:])
		sha1Hash.Sum(sha1Bytes[:0])

		if sha1Bytes[0] == 0x00 {
			return true
		}
	}

	return false
}

func patchTimelimit(f *os.File, tmdLen uint32, p bool) error {
	regionStr := regionStr
	var region byte
	input := bufio.NewReader(os.Stdin)

	_, err := f.Seek(int64(offsetTMD), 0)
	if err != nil {
		return err
	}

	tmd := make([]byte, tmdLen)
	_, err = f.Read(tmd)
	if err != nil {
		return err
	}

	region = tmd[offsetRegion]
	if p {
		for {
			fmt.Printf("Region is set to %s\n", regionStr[region])
			fmt.Println("New region:")
			for i := 0; i < 4; i++ {
				fmt.Printf("%d- %s\n", i, regionStr[i])
			}
			fmt.Print("Enter your new choice: (oh man you gotta save him!): ")
			newRegion, _, _ := input.ReadLine()
			region = newRegion[0] - '0'
			if region > 4 {
				region = 0x04
			}

			fmt.Printf("New region is set to %s, do you agree? (y/n) ", regionStr[region])
			confirm, _, _ := input.ReadLine()
			if confirm[0] == 'y' {
				break
			}
		}
	} else {
		region = 0x03
	}

	tmd[offsetRegion] = region

	fmt.Println("\tSigning...")
	if !truchaSignTicket(tmd) {
		return fmt.Errorf("error signing TMD")
	}
	fmt.Println("\tdone.\n")

	_, err = f.Seek(int64(offsetTMD), 0)
	if err != nil {
		return err
	}

	_, err = f.Write(tmd)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	fmt.Println("Free the WADs *testing version* by Superken7\n")

	if len(os.Args) < 2 {
		fmt.Printf("usage:\t%s <WAD file>\n", os.Args[0])
		return
	}

	filename := os.Args[1]
	p := len(os.Args) >= 3

	f, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("cannot open file %s\n", filename)
		return
	}
	defer f.Close()

	var tmdLen uint32
	if _, err := f.Seek(0x14, 0); err != nil {
		fmt.Println("error seeking file:", err)
		return
	}
	if err := binary.Read(f, binary.BigEndian, &tmdLen); err != nil {
		fmt.Println("error reading file:", err)
		return
	}
	tmdLen = tmdLen // You might need to adjust this for endianess

	fmt.Printf("tmd_len %04x\n", tmdLen)

	fmt.Println("Patching... ")
	if err := patchTimelimit(f, tmdLen, p); err != nil {
		fmt.Println("error patching file:", err)
		return
	}
	fmt.Println("done.\n\nFOR FREEEDOOOOM!!")
}
