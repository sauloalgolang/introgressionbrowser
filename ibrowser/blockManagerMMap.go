// +build darwin dragonfly freebsd linux openbsd solaris netbsd

package ibrowser

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

// BlockManagerMMap - mmap struct
type BlockManagerMMap struct {
	Filename       string
	Mode           FileMode
	CounterBits    uint64
	RegisterLength uint64
	BlockLength    uint64
	NumBlocks      uint64
	NumRegisters   uint64
	NumBytes       uint64
	Data16         *DistanceRow16
	Data32         *DistanceRow32
	Data64         *DistanceRow64
	fd             int
	file           *os.File
	array          *[]byte
	mmapFile       *mmap.MMap
}

func (bmm BlockManagerMMap) String() (res string) {
	res = fmt.Sprintf("%s - Mode %s CounterBits %d BlockLength %d RegisterLength %d NumBlocks %d NumRegisters %d NumBytes %d",
		bmm.Filename,
		bmm.Mode,
		bmm.CounterBits,
		bmm.BlockLength,
		bmm.RegisterLength,
		bmm.NumBlocks,
		bmm.NumRegisters,
		bmm.NumBytes,
	)

	switch bmm.CounterBits {
	case 16:
		res += fmt.Sprintf("array length %d", len(*bmm.Data16))
	case 32:
		res += fmt.Sprintf("array length %d", len(*bmm.Data32))
	case 64:
		res += fmt.Sprintf("array length %d", len(*bmm.Data64))
	default:
		panic("invalid counterBits")
	}

	return
}

// NewBlockManagerMMap - creates a new instance of db
func NewBlockManagerMMap(filename string, counterBits uint64, blockLength uint64, mode FileMode) (bmm *BlockManagerMMap, err error) {
	fmt.Println("NewBlockManagerMMap")

	registerLength := uint64(0)

	switch counterBits {
	case 16:
		registerLength = uint64(unsafe.Sizeof(DistanceType16(0)))
	case 32:
		registerLength = uint64(unsafe.Sizeof(DistanceType32(0)))
	case 64:
		registerLength = uint64(unsafe.Sizeof(DistanceType64(0)))
	default:
		panic("invalid counterBits")
	}

	bmm = &BlockManagerMMap{
		Filename:       filename,
		Mode:           mode,
		CounterBits:    counterBits,
		RegisterLength: registerLength,
		BlockLength:    blockLength,
		NumBlocks:      0,
		NumRegisters:   0,
		NumBytes:       0,
		Data16:         nil,
		Data32:         nil,
		Data64:         nil,
		fd:             0,
		file:           nil,
		array:          nil,
		mmapFile:       nil,
	}

	if err = bmm.Init(); err != nil {
		return nil, err
	}

	return bmm, nil
}

// Init - Initialize files
func (bmm *BlockManagerMMap) Init() (err error) {
	fmt.Println("Initializing MMAP")

	if err = bmm.open(); err != nil {
		return err
	}

	if err = bmm.mmap(); err != nil {
		return err
	}

	return nil
}

// Close - close file
func (bmm *BlockManagerMMap) Close() (err error) {
	fmt.Println("Closing MMAP")

	if err = bmm.close(); err != nil {
		return
	}

	bmm.Filename = ""
	bmm.Mode = -1
	bmm.CounterBits = 0
	bmm.RegisterLength = 0
	bmm.BlockLength = 0
	bmm.NumBlocks = 0
	bmm.NumRegisters = 0
	bmm.NumBytes = 0
	bmm.Data16 = nil
	bmm.Data32 = nil
	bmm.Data64 = nil
	bmm.fd = 0
	bmm.file = nil
	bmm.array = nil
	bmm.mmapFile = nil

	return nil
}

// LenBlocks - returns the current number of blcks
func (bmm *BlockManagerMMap) LenBlocks() uint64 {
	return bmm.NumBlocks
}

// LenRegisters - returns the current number of registers
func (bmm *BlockManagerMMap) LenRegisters() uint64 {
	return bmm.NumRegisters
}

// LenBytes - returns the current number of bytes
func (bmm *BlockManagerMMap) LenBytes() uint64 {
	return bmm.NumBytes
}

// GetNewBlock - returns a newly created matrix
func (bmm *BlockManagerMMap) GetNewBlock() (interface{}, error) {
	nb := bmm.NumBlocks
	err := bmm.CreateNewBlock()

	if err != nil {
		return nil, err
	}

	data, err2 := bmm.GetBlockPos(nb)

	if err2 != nil {
		return nil, err2
	}

	return data, nil
}

// GetBlockPos - returns a matrix
func (bmm *BlockManagerMMap) GetBlockPos(pos uint64) (interface{}, error) {
	bl := bmm.BlockLength
	startPos := bl * pos
	endPos := bl * (pos + 1)

	if pos > bmm.NumBlocks {
		return nil, errors.New("Position > length")
	}

	switch bmm.CounterBits {
	case 16:
		frag := (*bmm.Data16)[startPos:endPos]
		return &frag, nil
	case 32:
		frag := (*bmm.Data32)[startPos:endPos]
		return &frag, nil
	case 64:
		frag := (*bmm.Data64)[startPos:endPos]
		return &frag, nil
	default:
		panic("invalid counterBits")
	}

	return nil, nil
}

// CreateNewBlock - add a new block
func (bmm *BlockManagerMMap) CreateNewBlock() error {
	fmt.Println("Creating New Block")

	return bmm.CreateNewBlocks(1)
}

// CreateNewBlocks - extends number of blocks
func (bmm *BlockManagerMMap) CreateNewBlocks(numNewBlocks uint64) (err error) {
	fmt.Println("Creating New Blocks: ", numNewBlocks)

	numNewRegisters := bmm.BlockLength * numNewBlocks

	err = bmm.CreateNewRegisters(numNewRegisters)

	return
}

// CreateNewRegisters - extends the NumBytes of the mmap
func (bmm *BlockManagerMMap) CreateNewRegisters(numNewRegisters uint64) (err error) {
	fmt.Println("Creating New Registers: ", numNewRegisters)

	numNewBytes := bmm.RegisterLength * numNewRegisters

	err = bmm.CreateNewBytes(numNewBytes)

	return err
}

// CreateNewBytes - extends the NumBytes of the mmap
func (bmm *BlockManagerMMap) CreateNewBytes(numNewBytes uint64) (err error) {
	fmt.Println("Creating New Bytes: ", numNewBytes)

	numBytes := bmm.NumBytes + numNewBytes

	err = bmm.ResizeToBytes(numBytes)

	return err
}

// ResizeToBlocks - Resize to a specific number of blocks
func (bmm *BlockManagerMMap) ResizeToBlocks(numBlocks uint64) (err error) {
	fmt.Println("Resizing To Blocks: ", numBlocks)

	numRegisters := bmm.BlockLength * numBlocks

	err = bmm.ResizeToRegisters(numRegisters)

	return err
}

// ResizeToRegisters - Resize to a specific number of registers
func (bmm *BlockManagerMMap) ResizeToRegisters(numRegisters uint64) (err error) { // in uint64
	fmt.Println("Resizing To Registers: ", numRegisters)

	numBytes := bmm.RegisterLength * numRegisters

	err = bmm.ResizeToBytes(numBytes)

	return err
}

// ResizeToBytes - Resize to a specific number of bytes
func (bmm *BlockManagerMMap) ResizeToBytes(numBytes uint64) (err error) { // in bytes
	fmt.Println("Resizing To Bytes: ", numBytes)

	if bmm.Mode == RO {
		return errors.New("Trying to resize in a read only file")
	}

	if err = bmm.close(); err != nil {
		return err
	}

	if err = bmm.open(); err != nil {
		return err
	}

	err = syscall.Ftruncate(bmm.fd, int64(numBytes))

	if err != nil {
		fmt.Println("Error resizing: ", err)
		return err
	}

	if err = bmm.updateSize(); err != nil {
		return err
	}

	if err = bmm.mmap(); err != nil {
		return err
	}

	return nil
}

// GetMatrixMaker - Return default matrix maker
func (bmm *BlockManagerMMap) GetMatrixMaker() MatrixMaker {
	return bmm.GetNewBlock
}

// GetFallbackMatrixMaker - Return fallback matrix maker
func (bmm *BlockManagerMMap) GetFallbackMatrixMaker() MatrixMaker {
	return bmm.GetNewBlock
}

func (bmm *BlockManagerMMap) open() (err error) {
	fmt.Println("Opening file")

	var f *os.File

	if bmm.Mode == RW {
		f, err = os.OpenFile(bmm.Filename, os.O_CREATE|os.O_RDWR, 0664)
	} else {
		f, err = os.OpenFile(bmm.Filename, os.O_RDONLY, 0664)
	}

	if err != nil {
		fmt.Println("Could not open file: ", err)
		return err
	}

	bmm.fd = int(f.Fd())
	bmm.file = f

	if err = bmm.updateSize(); err != nil {
		return err
	}

	return nil
}

func (bmm *BlockManagerMMap) close() (err error) {
	fmt.Println("Closing file")

	if bmm.mmapFile != nil {
		fmt.Println(" Flushing mmap")
		if err = bmm.mmapFile.Flush(); err != nil {
			return err
		}

		fmt.Println(" Unmapping mmap file")
		if err = bmm.mmapFile.Unmap(); err != nil {
			return err
		}
	}

	if err = bmm.file.Close(); err != nil {
		fmt.Println(" Closing file handler")
		return err
	}

	return nil
}

func (bmm *BlockManagerMMap) updateSize() error {
	fmt.Println("Updating size")

	fi, err := bmm.file.Stat()

	if err != nil {
		// Could not obtain stat, handle error
		fmt.Println("Could not obtain stat: ", err)
		return err
	}

	bmm.NumBytes = uint64(fi.Size())
	bmm.NumRegisters = bmm.NumBytes / bmm.RegisterLength
	bmm.NumBlocks = bmm.NumRegisters / bmm.BlockLength

	return nil
}

func (bmm *BlockManagerMMap) mmap() (err error) {
	fmt.Println("mmapping: ", bmm.NumBytes, "bytes - ", bmm.NumRegisters, "registers")

	var mmapFile mmap.MMap

	if bmm.Mode == RW {
		mmapFile, err = mmap.Map(bmm.file, mmap.RDWR, 0)
	} else {
		mmapFile, err = mmap.Map(bmm.file, mmap.RDONLY, 0)
	}

	if err != nil {
		bmm.Close()
		return err
	}

	bytem := []byte(mmapFile)
	bmm.mmapFile = &mmapFile
	bmm.array = &bytem
	bmm.unsafeBytesToRegister()

	fmt.Print(" mmapped: ", bmm.NumBytes, "bytes - ", bmm.NumRegisters, "registers")

	switch bmm.CounterBits {
	case 16:
		fmt.Println("array length", len(*bmm.Data16))
	case 32:
		fmt.Println("array length", len(*bmm.Data32))
	case 64:
		fmt.Println("array length", len(*bmm.Data64))
	default:
		panic("invalid counterBits")
	}

	return nil
}

func (bmm *BlockManagerMMap) unsafeRegisterToBytes() {
	var ptr uintptr
	var count int

	if bmm.CounterBits == 16 {
		count = len(*bmm.Data16) * int(bmm.RegisterLength)
		ptr = uintptr(unsafe.Pointer(&(*bmm.Data16)[0]))
	} else if bmm.CounterBits == 32 {
		count = len(*bmm.Data32) * int(bmm.RegisterLength)
		ptr = uintptr(unsafe.Pointer(&(*bmm.Data32)[0]))
	} else if bmm.CounterBits == 64 {
		count = len(*bmm.Data64) * int(bmm.RegisterLength)
		ptr = uintptr(unsafe.Pointer(&(*bmm.Data64)[0]))
	}

	slice := reflect.SliceHeader{
		Data: ptr,
		Len:  count,
		Cap:  count,
	}

	bmm.array = (*[]byte)(unsafe.Pointer(&slice))
}

func (bmm *BlockManagerMMap) unsafeBytesToRegister() {
	count := len(*bmm.array) / int(bmm.RegisterLength)

	slice := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&(*bmm.array)[0])),
		Len:  count,
		Cap:  count,
	}

	ptr := unsafe.Pointer(&slice)

	switch bmm.CounterBits {
	case 16:
		bmm.Data16 = (*DistanceRow16)(ptr)
	case 32:
		bmm.Data32 = (*DistanceRow32)(ptr)
	case 64:
		bmm.Data64 = (*DistanceRow64)(ptr)
	default:
		panic("invalid counterBits")
	}
}

// func (bmm *BlockManagerMMap) mmapRaw() (err error) {
// 	fmt.Println("mmapping: ", bmm.NumBytes, "bytes - ", bmm.NumRegisters, "registers")

// 	if bmm.NumBytes == 0 {
// 		return nil
// 	}

// 	var array []byte

// 	if bmm.Mode == RW {
// 		array, err = syscall.Mmap(bmm.fd, 0, int(bmm.NumBytes), syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
// 	} else {
// 		array, err = syscall.Mmap(bmm.fd, 0, int(bmm.NumBytes), syscall.PROT_READ, syscall.MAP_SHARED)
// 	}

// 	if err != nil {
// 		fmt.Println("Error mmapping: ", err)
// 		bmm.Close()
// 		return err
// 	}

// 	bmm.array = &array
// 	bmm.unsafeBytesToRegister()

// 	fmt.Print(" mmapped: ", bmm.NumBytes, "bytes - ", bmm.NumRegisters, "registers")

// 	switch bmm.CounterBits {
// 	case 16:
// 		fmt.Println("array length", len(*bmm.Data16))
// 	case 32:
// 		fmt.Println("array length", len(*bmm.Data32))
// 	case 64:
// 		fmt.Println("array length", len(*bmm.Data64))
// 	default:
// 		panic("invalid counterBits")
// 	}

// 	return nil
// }

// func (bmm *BlockManagerMMap) flush(addr, len uintptr) error {
// 	_, _, errno := syscall.Syscall(_SYS_MSYNC, addr, len, _MS_SYNC)
// 	if errno != 0 {
// 		return syscall.Errno(errno)
// 	}
// 	return nil
// }
