package main

import (
	"fmt"
	"os"
	"reflect"
	"unsafe"

	"github.com/edsrzf/mmap-go"
)

// Strings interface for a list of strings
type Strings interface {
	Get(i int) string
	Len() int
	Close() error
}

type mapfile struct {
	file *os.File
	mmap mmap.MMap
}

func openRead(filename string) (m *mapfile, err error) {
	m = &mapfile{}
	m.file, err = os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	m.mmap, err = mmap.Map(m.file, mmap.RDONLY, 0)
	if err != nil {
		m.file.Close()
		return nil, err
	}
	return m, nil
}

func (m *mapfile) Data() []byte { return []byte(m.mmap) }

func (m *mapfile) Close() error {
	_ = m.mmap.Unmap()
	return m.file.Close()
}

type strings struct {
	data  []byte
	index []int64

	datafile  *mapfile
	indexfile *mapfile
}

func (strs *strings) Len() int {
	return len(strs.index)
}

func (strs *strings) Get(i int) string {
	start := int(strs.index[i])
	end := len(strs.data)

	if i+1 < len(strs.index) {
		end = int(strs.index[i+1])
	}

	return string(strs.data[start:end])
}

func (strs *strings) Close() error {
	a := strs.datafile.Close()
	b := strs.indexfile.Close()

	if a != nil {
		return a
	}

	if b != nil {
		return b
	}

	return nil
}

func unsafeInt64ToBytes(xs []int64) []byte {
	var v int64
	count := len(xs) * int(unsafe.Sizeof(v))
	slice := reflect.SliceHeader{uintptr(unsafe.Pointer(&xs[0])), count, count}
	return *(*[]byte)(unsafe.Pointer(&slice))
}

func unsafeBytesToInt64(xs []byte) []int64 {
	var v int64
	count := len(xs) / int(unsafe.Sizeof(v))
	slice := reflect.SliceHeader{uintptr(unsafe.Pointer(&xs[0])), count, count}
	return *(*[]int64)(unsafe.Pointer(&slice))
}

func (strs *strings) save(prefix string) error {
	dat, err := os.OpenFile(prefix+".dat", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	dat.Write(strs.data)
	dat.Close()

	idx, err := os.OpenFile(prefix+".idx", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	idx.Write(unsafeInt64ToBytes(strs.index))
	idx.Close()

	return nil
}

func (strs *strings) load(prefix string) error {
	var err error

	strs.datafile, err = openRead(prefix + ".dat")
	if err != nil {
		return err
	}
	strs.indexfile, err = openRead(prefix + ".idx")
	if err != nil {
		strs.datafile.Close()
		strs.datafile = nil
		return err
	}

	strs.data = []byte(strs.datafile.Data())
	strs.index = unsafeBytesToInt64(strs.indexfile.Data())

	return nil
}

// SaveTo saves string list to file
func SaveTo(prefix string, xs []string) error {
	t := 0
	for _, x := range xs {
		t += len(x)
	}

	strs := &strings{
		data:  make([]byte, 0, t),
		index: make([]int64, 0, len(xs)),
	}

	for _, x := range xs {
		strs.index = append(strs.index, int64(len(strs.data)))
		strs.data = append(strs.data, []byte(x)...)
	}

	return strs.save(prefix)
}

// Load loads string list from file
func Load(prefix string) (Strings, error) {
	strs := &strings{}
	return strs, strs.load(prefix)
}

func main() {
	err := SaveTo("example", []string{"alpha", "beta", "gamma", "delta"})
	if err != nil {
		panic(err)
	}
	strs, err := Load("example")
	if err != nil {
		panic(err)
	}
	defer strs.Close()

	for i := 0; i < strs.Len(); i++ {
		fmt.Println(strs.Get(i))
	}
}
