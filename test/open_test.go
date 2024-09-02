/*
Copyright (c) 2018-2021 the Go Library for Virtual Disk Development Kit contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/yrudtlr/virtual-disks/pkg/disklib"
	"github.com/yrudtlr/virtual-disks/pkg/virtual_disks"
)

func TestOpen(t *testing.T) {
	fmt.Println("Test Open")
	var majorVersion uint32 = 7
	var minorVersion uint32 = 0
	path := os.Getenv("LIBPATH")
	if path == "" {
		t.Skip("Skipping testing if environment variables are not set.")
	}
	disklib.Init(majorVersion, minorVersion, path)
	serverName := os.Getenv("IP")
	thumPrint := os.Getenv("THUMBPRINT")
	userName := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	fcdId := os.Getenv("FCDID")
	ds := os.Getenv("DATASTORE")
	identity := os.Getenv("IDENTITY")
	params := disklib.NewConnectParams("", serverName, thumPrint, userName,
		password, fcdId, ds, "", "", identity, "", disklib.VIXDISKLIB_FLAG_OPEN_COMPRESSION_SKIPZ,
		false, disklib.NBD)
	diskReaderWriter, err := virtual_disks.Open(params)
	if err != nil {
		disklib.EndAccess(params)
		t.Errorf("Open failed, got error code: %d, error message: %s.", err.VixErrorCode(), err.Error())
	}
	// QAB (assume at least 1GiB volume and 1MiB block size)
	abInitial, err := diskReaderWriter.QueryAllocatedBlocks(0, 2048*1024, 2048)
	if err != nil {
		t.Errorf("QueryAllocatedBlocks failed: %d, error message: %s", err.VixErrorCode(), err.Error())
	} else {
		fmt.Printf("Number of blocks: %d\n", len(abInitial))
		fmt.Printf("Offset      Length\n")
		for _, ab := range abInitial {
			fmt.Printf("0x%012x  0x%012x\n", ab.Offset(), ab.Length())
		}
	}
	// ReadAt
	fmt.Printf("ReadAt test\n")
	buffer := make([]byte, disklib.VIXDISKLIB_SECTOR_SIZE)
	n, err4 := diskReaderWriter.Read(buffer)
	fmt.Printf("Read byte n = %d\n", n)
	fmt.Println(buffer)
	fmt.Println(err4)

	// WriteAt
	fmt.Println("WriteAt start")
	buf1 := make([]byte, disklib.VIXDISKLIB_SECTOR_SIZE)
	for i := range buf1 {
		buf1[i] = 'E'
	}
	n2, err2 := diskReaderWriter.WriteAt(buf1, 0)
	fmt.Printf("Write byte n = %d\n", n2)
	fmt.Println(err2)

	buffer2 := make([]byte, disklib.VIXDISKLIB_SECTOR_SIZE)
	n2, err5 := diskReaderWriter.ReadAt(buffer2, 0)
	fmt.Printf("Read byte n = %d\n", n2)
	fmt.Println(buffer2)
	fmt.Println(err5)

	// QAB (assume at least 1GiB volume and 1MiB block size)
	abFinal, err := diskReaderWriter.QueryAllocatedBlocks(0, 2048*1024, 2048)
	if err != nil {
		t.Errorf("QueryAllocatedBlocks failed: %d, error message: %s", err.VixErrorCode(), err.Error())
	} else {
		fmt.Printf("Number of blocks: %d\n", len(abInitial))
		fmt.Printf("Offset      Length\n")
		for _, ab := range abFinal {
			fmt.Printf("0x%012x  0x%012x\n", ab.Offset(), ab.Length())
		}
	}

	diskReaderWriter.Close()
}
