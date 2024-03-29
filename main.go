package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/sys/unix"
)

// https://terinstock.com/post/2018/10/memfd_create-Temporary-in-memory-files-with-Go-and-Linux/

// memfile takes a file name used, and the byte slice
// containing data the file should contain.
//
// name does not need to be unique, as it's used only
// for debugging purposes.
//
// It is up to the caller to close the returned descriptor.
func memfile(name string, b []byte) (int, error) {
	fd, err := unix.MemfdCreate(name, 0)
	if err != nil {
		return 0, fmt.Errorf("MemfdCreate: %v", err)
	}

	err = unix.Ftruncate(fd, int64(len(b)))
	if err != nil {
		return 0, fmt.Errorf("Ftruncate: %v", err)
	}

	data, err := unix.Mmap(fd, 0, len(b), unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		return 0, fmt.Errorf("Mmap: %v", err)
	}

	copy(data, b)

	err = unix.Munmap(data)
	if err != nil {
		return 0, fmt.Errorf("Munmap: %v", err)
	}

	return fd, nil
}

func main() {
	fd, err := memfile("hello", []byte("hello world!"))
	if err != nil {
		log.Fatalf("memfile: %+v", err)
	}

	// filepath to our newly created in-memory file descriptor
	fp := fmt.Sprintf("/proc/self/fd/%d", fd)

	// create an *os.File, should you need it
	// alternatively, pass fd or fp as input to a library.
	f := os.NewFile(uintptr(fd), fp)
	defer f.Close()

  b, err := io.ReadAll(f)
  if err != nil {
		log.Fatalf("memfile read: %+v", err)
	}
  fmt.Printf("File name: %s\n", f.Name())
  fmt.Printf("File content:\n%s\nEnd File Content\n", b)
}