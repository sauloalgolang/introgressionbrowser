package main

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"syscall"
	"unsafe"
)

import (
	"github.com/edsrzf/mmap-go"
)

// https://medium.com/@arpith/adventures-with-mmap-463b33405223

// RegisterType - array type
type RegisterType = uint64

// FileMode - File mode
type FileMode int

const (
	// RO - Read only
	RO FileMode = 0
	// RW - Read and Write
	RW FileMode = 1
)

func (mode FileMode) String() string {
	names := [...]string{
		"RO",
		"RW",
	}

	if mode < RO || mode > RW {
		return "Unknown"
	}

	return names[mode]
}

// Db - mmap struct
type Db struct {
	Filename       string
	Mode           FileMode
	Data           *[]RegisterType
	array          *[]byte
	file           *os.File
	mmapFile       *mmap.MMap
	fd             int
	size           int64
	length         int64
	registerLength int64
}

// NewDb - creates a new instance of db
func NewDb(filename string, mode FileMode) (dbi *Db, err error) {
	var registerType RegisterType

	dbi = &Db{
		Filename:       filename,
		Mode:           mode,
		registerLength: int64(unsafe.Sizeof(registerType)),
	}

	if err = dbi.open(); err != nil {
		return nil, err
	}

	if err = dbi.mmap(); err != nil {
		return nil, err
	}

	return dbi, nil
}

func (db Db) String() string {
	return fmt.Sprintf("%s - Mode %s Size %d length %d register length %d",
		db.Filename,
		db.Mode,
		db.size,
		db.length,
		db.registerLength,
	)
}

// GetLength - returns the current register length
func (db *Db) Len() int64 {
	return db.length
}

// Append - extends the size of the mmap by 1
func (db *Db) Append() (err error) {
	return db.Extend(1)
}

// Extend - extends the size of the mmap
func (db *Db) Extend(size int64) (err error) {
	if db.Mode == RO {
		return errors.New("Trying to resize in a read only file")
	}

	if err = db.file.Close(); err != nil {
		return err
	}

	if err = db.open(); err != nil {
		return err
	}

	if err = db.resizeRegisters(size); err != nil {
		return err
	}

	if err = db.mmap(); err != nil {
		return err
	}

	return nil
}

// Close - close file
func (db *Db) Close() (err error) {
	if db.mmapFile != nil {
		if err = db.mmapFile.Flush(); err != nil {
			return err
		}

		if err = db.mmapFile.Unmap(); err != nil {
			return err
		}
	}

	if err = db.file.Close(); err != nil {
		return err
	}

	db.Filename = ""
	db.Mode = -1
	db.file = nil
	db.mmapFile = nil
	db.Data = nil
	db.array = nil
	db.fd = 0
	db.size = 0
	db.length = 0
	db.registerLength = 0

	return nil
}

func (db *Db) resizeRegisters(length int64) (err error) { // in uint64
	fmt.Println("Resizing registers: ", length)

	size := length * db.registerLength

	err = db.resizeBytes(size)

	return err
}

func (db *Db) resizeBytes(size int64) (err error) { // in bytes
	fmt.Println("Resizing bytes: ", size)

	err = syscall.Ftruncate(db.fd, size)

	if err != nil {
		fmt.Println("Error resizing: ", err)
		return err
	}

	db.size = size
	db.length = size / db.registerLength

	return nil
}

func (db *Db) open() (err error) {
	fmt.Println("Getting file descriptor")

	var f *os.File

	if db.Mode == RW {
		f, err = os.OpenFile(db.Filename, os.O_CREATE|os.O_RDWR, 0664)
	} else {
		f, err = os.OpenFile(db.Filename, os.O_RDONLY, 0664)
	}

	if err != nil {
		fmt.Println("Could not open file: ", err)
		return err
	}

	db.fd = int(f.Fd())
	db.file = f

	fi, err := f.Stat()
	if err != nil {
		// Could not obtain stat, handle error
		fmt.Println("Could not obtain stat: ", err)
		return err
	}

	db.size = fi.Size()
	db.length = db.size / db.registerLength

	return nil
}

func (db *Db) mmap() (err error) {
	fmt.Println("mmapping: ", db.size, "bytes - ", db.length, "registers")

	var mmapFile mmap.MMap

	if db.Mode == RW {
		mmapFile, err = mmap.Map(db.file, mmap.RDWR, 0)
	} else {
		mmapFile, err = mmap.Map(db.file, mmap.RDONLY, 0)
	}

	if err != nil {
		db.Close()
		return err
	}

	db.mmapFile = &mmapFile
	bytem := []byte(mmapFile)
	db.array = &bytem
	db.unsafeBytesToRegister()

	fmt.Println(" mmapped: ", db.size, "bytes - ", db.length, "registers", "array length", len(*db.Data))

	return nil
}

func (db *Db) mmapRaw() (err error) {
	fmt.Println("mmapping: ", db.size, "bytes - ", db.length, "registers")

	if db.size == 0 {
		return nil
	}

	var array []byte

	if db.Mode == RW {
		array, err = syscall.Mmap(db.fd, 0, int(db.size), syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	} else {
		array, err = syscall.Mmap(db.fd, 0, int(db.size), syscall.PROT_READ, syscall.MAP_SHARED)
	}

	if err != nil {
		fmt.Println("Error mmapping: ", err)
		db.Close()
		return err
	}

	db.array = &array
	db.unsafeBytesToRegister()

	fmt.Println(" mmapped: ", db.size, "bytes - ", db.length, "registers", "array length", len(*db.Data))

	return nil
}

// func (db *Db) flush(addr, len uintptr) error {
// 	_, _, errno := syscall.Syscall(_SYS_MSYNC, addr, len, _MS_SYNC)
// 	if errno != 0 {
// 		return syscall.Errno(errno)
// 	}
// 	return nil
// }

func (db *Db) unsafeRegisterToBytes() {
	count := len(*db.Data) * int(db.registerLength)

	slice := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&(*db.Data)[0])),
		Len:  count,
		Cap:  count,
	}

	db.array = (*[]byte)(unsafe.Pointer(&slice))
}

func (db *Db) unsafeBytesToRegister() {
	count := len(*db.array) / int(db.registerLength)

	slice := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&(*db.array)[0])),
		Len:  count,
		Cap:  count,
	}

	db.Data = (*[]RegisterType)(unsafe.Pointer(&slice))
}

func main() {
	println("Creating mmap2.bin")

	loopSize := int64(5000)

	db, err := NewDb("mmap2.bin", RW)

	if err != nil {
		panic(err)
	}

	println(" Created")

	fmt.Println(db)

	println("Extending")

	db.Extend(loopSize)

	println(" Extended")

	fmt.Println(db)

	println("Looping")

	for i := int64(0); i < loopSize; i++ {
		// fmt.Println((*db.Data)[i])
	}

	println("Changing")

	for i := int64(0); i < loopSize; i++ {
		(*db.Data)[i] = uint64(i)
	}

	println("Looping")

	for i := int64(0); i < loopSize; i++ {
		// fmt.Println((*db.Data)[i])
		if (*db.Data)[i] != uint64(i) {
			panic(fmt.Sprintf("value mismatch %d != %d", i, (*db.Data)[i]))
		}
	}

	println("Closing")

	fmt.Println(db)

	db.Close()

	println(" Closed")

	println("Opening again")

	db2, err2 := NewDb("mmap2.bin", RO)

	if err2 != nil {
		panic(err2)
	}

	println("Looping", db2.Len())

	for i := int64(0); i < db2.Len(); i++ {
		// fmt.Println((*db2.Data)[i])
	}

	println("Closing")

	fmt.Println(db2)

	db2.Close()

	println(" Closed")
}
