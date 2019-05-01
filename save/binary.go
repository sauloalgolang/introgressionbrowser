package save

import (
	"bufio"
	// "bytes"
	"encoding/binary"
	"log"
	"math"
	"os"
	// "io/ioutil"
	// "fmt"
)

// https://golang.org/pkg/encoding/binary/
// https://varunpant.com/posts/reading-and-writing-binary-files-in-go-lang

//
// MultiArrayFile
//

type MultiArrayFile struct {
	fileName   string
	endianness binary.ByteOrder
	// buf          *bytes.Buffer
	serial       uint64
	counterBits  uint64
	dataLen      uint64
	isFinished   bool
	isSoft       bool
	writeMode    bool
	bufReaderIdx *bufio.Reader
	bufWriterIdx *bufio.Writer
	bufReaderDta *bufio.Reader
	bufWriterDta *bufio.Writer
	fileIdx      *os.File
	fileDta      *os.File
}

type RegisterLocation struct {
	Size           uint64
	HeaderSize     uint64
	MatrixSize     uint64
	StartPosition  uint64
	MatrixPosition uint64
	EndPosition    uint64
}

type RegisterHeader struct {
	HasData     bool
	Serial      uint64
	CounterBits uint64
	DataLen     uint64
	SumData     uint64
}

//
// New
//

func NewMultiArrayFile(fileName string, isSave bool, isSoft bool) *MultiArrayFile {
	if isSave && isSoft {
		log.Fatal("Cant save in soft mode")
		os.Exit(1)
	}

	m := MultiArrayFile{
		fileName:   fileName,
		endianness: binary.LittleEndian,
		// buf:         new(bytes.Buffer),
		serial:      0,
		counterBits: 0,
		dataLen:     0,
		isFinished:  false,
		isSoft:      isSoft,
	}

	if isSave {
		log.Println("Saving binary matrix to", fileName)

		fileIdx, err1 := os.Create(fileName + ".idx")
		if err1 != nil {
			log.Fatalln(err1)
		}

		fileDta, err2 := os.Create(fileName)
		if err2 != nil {
			log.Fatalln(err2)
		}

		m.writeMode = true
		m.bufWriterIdx = bufio.NewWriter(fileIdx)
		m.bufWriterDta = bufio.NewWriter(fileDta)
		m.fileIdx = fileIdx
		m.fileDta = fileDta

	} else {
		log.Println("Loading binary matrix from", fileName)

		fileIdx, err1 := os.Open(fileName + ".idx")
		if err1 != nil {
			log.Fatalln(err1)
		}

		fileDta, err2 := os.Open(fileName)
		if err2 != nil {
			log.Fatalln(err2)
		}

		m.writeMode = false
		m.bufReaderIdx = bufio.NewReader(fileIdx)
		m.bufReaderDta = bufio.NewReader(fileDta)
		m.fileIdx = fileIdx
		m.fileDta = fileDta
	}

	return &m
}

//
// TODO: MMAP
//
// https://stackoverflow.com/questions/9203526/mapping-an-array-to-a-file-via-mmap-in-go
//
// package main
//
// import (
//     "fmt"
//     "os"
//     "syscall"
//     "unsafe"
// )
//
// func main() {
//     const n = 1e3
//     t := int(unsafe.Sizeof(0)) * n
//
//     map_file, err := os.Create("/tmp/test.dat")
//     if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//     }
//     _, err = map_file.Seek(int64(t-1), 0)
//     if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//     }
//     _, err = map_file.Write([]byte(" "))
//     if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//     }
//
// 	   // func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error)
//     mmap, err := syscall.Mmap(int(map_file.Fd()), 0, int(t), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
//     if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//     }
//     map_array := (*[n]int)(unsafe.Pointer(&mmap[0]))

func (m *MultiArrayFile) SetSerial(serial uint64) {
	m.serial = serial
}

func (m *MultiArrayFile) GetSerial() (serial uint64) {
	return m.serial
}

func (m *MultiArrayFile) CalculateRegisterHeaderSize(counterBits uint64, size uint64) (res uint64) {
	res += 1 // hasData     bool
	res += 8 // serial      int64
	res += 8 // counterBits uint64
	res += 8 // dataLen     int64
	res += 8 // sumData     uint64
	return
}

func (m *MultiArrayFile) CalculateRegisterMatrixSize(counterBits uint64, size uint64) (res uint64) {
	dbytes := uint64(0)
	switch counterBits {
	case 16:
		dbytes = 2
	case 32:
		dbytes = 4
	case 64:
		dbytes = 8
	}

	if dbytes == 0 {
		panic("wrong counterbits")
	}

	res += dbytes * size

	return
}

func (m *MultiArrayFile) CalculateRegisterSize(counterBits uint64, size uint64) (res uint64) {
	res += m.CalculateRegisterHeaderSize(counterBits, size)
	res += m.CalculateRegisterMatrixSize(counterBits, size)

	return
}

func (m *MultiArrayFile) CalculateRegisterLocation(counterBits uint64, size uint64, serial uint64) (res RegisterLocation) {
	res = RegisterLocation{
		HeaderSize: m.CalculateRegisterHeaderSize(counterBits, size),
		MatrixSize: m.CalculateRegisterMatrixSize(counterBits, size),
	}

	res.Size = m.CalculateRegisterSize(counterBits, size)
	res.StartPosition = res.Size * serial
	res.MatrixPosition = res.StartPosition + res.HeaderSize
	res.EndPosition = res.StartPosition + res.Size

	if res.EndPosition != (res.Size * (serial + 1)) {
		log.Panic("error calculation end position")
	}

	return
}

//
// MultiArrayFile :: Writer
//

// func (m *MultiArrayFile) write() (serial uint64) {
// 	if !m.writeMode {
// 		log.Fatalln("Trying to write to a reader")
// 	}

// 	hasData := true

// 	err := binary.Write(m.bufWriter, m.endianness, &hasData)
// 	if err != nil {
// 		log.Fatalln("binary.Write failed to write hasData:", err)
// 	}

// 	err = binary.Write(m.bufWriter, m.endianness, &m.serial)
// 	if err != nil {
// 		log.Fatalln("binary.Write failed to write serial:", err)
// 	}

// 	err = binary.Write(m.bufWriter, m.endianness, &m.counterBits)
// 	if err != nil {
// 		log.Fatalln("binary.Write failed to write counterBits:", err)
// 	}

// 	err = binary.Write(m.bufWriter, m.endianness, &m.dataLen)
// 	if err != nil {
// 		log.Fatalln("binary.Write failed to write dataLen:", err)
// 	}

// 	serial = m.serial

// 	m.serial++

// 	return serial
// }

func (m *MultiArrayFile) Write16(data *[]uint16) (serial uint64) {
	if m.counterBits == 0 {
		m.counterBits = 16
	} else if m.counterBits != 16 {
		log.Panicln("can't write different bits", m.counterBits, " != ", 16)
	}

	dataLen := uint64(len(*data))

	if m.dataLen == 0 {
		m.dataLen = dataLen
	} else if m.dataLen != dataLen {
		log.Panicln("can't write different sizes", m.dataLen, " != ", dataLen)
	}

	serial = m.serial

	ndata := make([]uint16, dataLen, dataLen)
	sumData := uint64(0)

	for i, v := range *data {
		if int64(v) > int64(math.MaxInt16) {
			log.Panicln("overflow")
		}

		ndata[i] = uint16(v)
		sumData += uint64(v)
	}

	header := RegisterHeader{
		HasData:     true,
		Serial:      serial,
		CounterBits: m.counterBits,
		DataLen:     m.dataLen,
		SumData:     sumData,
	}

	err1 := binary.Write(m.bufWriterIdx, m.endianness, &header)

	if err1 != nil {
		log.Fatalln("binary.Write failed to write data16:", err1)
	}

	err2 := binary.Write(m.bufWriterDta, m.endianness, &ndata)

	if err2 != nil {
		log.Fatalln("binary.Write failed to write data16:", err2)
	}

	m.serial++

	return
}

func (m *MultiArrayFile) Write32(data *[]uint32) (serial uint64) {
	if m.counterBits == 0 {
		m.counterBits = 32
	} else if m.counterBits != 32 {
		log.Panicln("can't write different bits", m.counterBits, " != ", 32)
	}

	dataLen := uint64(len(*data))

	if m.dataLen == 0 {
		m.dataLen = dataLen
	} else if m.dataLen != dataLen {
		log.Panicln("can't write different sizes", m.dataLen, " != ", dataLen)
	}

	serial = m.serial

	ndata := make([]uint32, dataLen, dataLen)
	sumData := uint64(0)

	for i, v := range *data {
		if int64(v) > int64(math.MaxInt32) {
			log.Panicln("overflow")
		}

		ndata[i] = uint32(v)
		sumData += uint64(v)
	}

	header := RegisterHeader{
		HasData:     true,
		Serial:      serial,
		CounterBits: m.counterBits,
		DataLen:     m.dataLen,
		SumData:     sumData,
	}

	err1 := binary.Write(m.bufWriterIdx, m.endianness, &header)

	if err1 != nil {
		log.Fatalln("binary.Write failed to write data32:", err1)
	}

	err2 := binary.Write(m.bufWriterDta, m.endianness, &ndata)

	if err2 != nil {
		log.Fatalln("binary.Write failed to write data32:", err2)
	}

	m.serial++

	return
}

func (m *MultiArrayFile) Write64(data *[]uint64) (serial uint64) {
	if m.counterBits == 0 {
		m.counterBits = 64
	} else if m.counterBits != 64 {
		log.Panicln("can't write different bits", m.counterBits, " != ", 64)
	}

	dataLen := uint64(len(*data))

	if m.dataLen == 0 {
		m.dataLen = dataLen
	} else if m.dataLen != dataLen {
		log.Panicln("can't write different sizes", m.dataLen, " != ", dataLen)
	}

	serial = m.serial

	sumData := uint64(0)

	for _, v := range *data {
		sumData += v
	}

	header := RegisterHeader{
		HasData:     true,
		Serial:      serial,
		CounterBits: m.counterBits,
		DataLen:     m.dataLen,
		SumData:     sumData,
	}

	err1 := binary.Write(m.bufWriterIdx, m.endianness, &header)

	if err1 != nil {
		log.Fatalln("binary.Write failed to write data64:", err1)
	}

	err2 := binary.Write(m.bufWriterDta, m.endianness, &data)

	if err2 != nil {
		log.Fatalln("binary.Write failed to write data64:", err2)
	}

	m.serial++

	return
}

//
// MultiArrayFile :: Reader
//

// func (m *MultiArrayFile) read() (hasData bool, serial uint64, counterBits uint64, dataLen uint64, sumData uint64) {
// 	if m.writeMode {
// 		log.Fatalln("Trying to read from a writer")
// 	}

// 	if m.isFinished {
// 		log.Fatalln("Trying to read a finished file")
// 	}

// 	hasData = false
// 	serial = uint64(0)
// 	counterBits = uint64(0)
// 	dataLen = uint64(0)
// 	sumData = uint64(0)

// 	//
// 	// HasData

// 	err := binary.Read(m.bufReader, m.endianness, &hasData)

// 	if err != nil {
// 		log.Fatalln("binary.Read failed reading hasData:", err)
// 	}

// 	if !hasData {
// 		m.isFinished = true
// 		return hasData, serial, counterBits, dataLen, sumData
// 	}

// 	//
// 	// Serial
// 	err = binary.Read(m.bufReader, m.endianness, &serial)

// 	if err != nil {
// 		log.Fatalln("binary.Read failed reading serial:", err)
// 	}

// 	if serial < 0 {
// 		log.Fatalln("serial < 0", serial)
// 	}

// 	if serial != m.serial {
// 		log.Fatalln("serial out of order", serial, " != ", m.serial)
// 	}

// 	//
// 	// counterBits
// 	err = binary.Read(m.bufReader, m.endianness, &counterBits)

// 	if err != nil {
// 		log.Fatalln("binary.Read failed reading counterBits:", err)
// 	}

// 	if counterBits <= 0 {
// 		log.Fatalln("Length <= 0", counterBits)
// 	} else if counterBits != 16 && counterBits != 32 && counterBits != 64 {
// 		log.Fatalln("Length not 16, 32 or 64", counterBits)
// 	}

// 	if m.counterBits == 0 {
// 		m.counterBits = counterBits
// 	} else if m.counterBits != counterBits {
// 		log.Fatalln("counterBits mismatch", m.counterBits != counterBits)
// 	}

// 	//
// 	// dataLen
// 	err = binary.Read(m.bufReader, m.endianness, &dataLen)

// 	if err != nil {
// 		log.Fatalln("binary.Read failed reading dataLen:", err)
// 	}

// 	if dataLen <= 0 {
// 		log.Fatalln("Length <= 0", dataLen)
// 	}

// 	//
// 	// sumData
// 	err = binary.Read(m.bufReader, m.endianness, &sumData)

// 	if err != nil {
// 		log.Fatalln("binary.Read failed reading sumData:", err)
// 	}

// 	m.serial++

// 	return hasData, serial, counterBits, dataLen, sumData
// }

func (m *MultiArrayFile) Read16(data *[]uint16) (hasData bool, serial uint64) {
	// if m.isSoft {
	// registerLocation := m.CalculateRegisterLocation(counterBits, dataLen, serial)
	// log.Println(registerLocation)

	// func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error)
	// mmap, err := syscall.Mmap(int(m.file.Fd()), 0, int(t), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// map_array := (*[n]int)(unsafe.Pointer(&mmap[0]))
	// } else {
	header := RegisterHeader{}

	err1 := binary.Read(m.bufReaderIdx, m.endianness, &header)

	if err1 != nil {
		log.Fatalln("binary.Read failed reading data16:", err1)
	}

	(*data) = make([]uint16, header.DataLen, header.DataLen)

	err2 := binary.Read(m.bufReaderDta, m.endianness, data)

	if err2 != nil {
		log.Fatalln("binary.Read failed reading data16:", err2)
	}

	sumDataV := uint64(0)
	for _, w := range *data {
		sumDataV += uint64(w)
	}

	if header.SumData != sumDataV {
		log.Fatalln("binary.Read failed reading data16: checksum error", header.SumData, sumDataV)
	}
	return header.HasData, header.Serial
	// }
	return false, 0
}

func (m *MultiArrayFile) Read32(data *[]uint32) (hasData bool, serial uint64) {
	// if m.isSoft {
	// registerLocation := m.CalculateRegisterLocation(counterBits, dataLen, serial)
	// log.Println(registerLocation)

	// func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error)
	// mmap, err := syscall.Mmap(int(m.file.Fd()), 0, int(t), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// map_array := (*[n]int)(unsafe.Pointer(&mmap[0]))
	// } else {
	header := RegisterHeader{}

	err1 := binary.Read(m.bufReaderIdx, m.endianness, &header)

	if err1 != nil {
		log.Fatalln("binary.Read failed reading data32:", err1)
	}

	log.Println(header)

	(*data) = make([]uint32, header.DataLen, header.DataLen)

	err2 := binary.Read(m.bufReaderDta, m.endianness, data)

	if err2 != nil {
		log.Fatalln("binary.Read failed reading data32:", err2)
	}

	sumDataV := uint64(0)
	for _, w := range *data {
		sumDataV += uint64(w)
	}

	if header.SumData != sumDataV {
		log.Fatalln("binary.Read failed reading data32: checksum error", header.SumData, sumDataV)
	}
	return header.HasData, header.Serial
	// }
	return false, 0
}

func (m *MultiArrayFile) Read64(data *[]uint64) (hasData bool, serial uint64) {
	// if m.isSoft {
	// registerLocation := m.CalculateRegisterLocation(counterBits, dataLen, serial)
	// log.Println(registerLocation)

	// func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error)
	// mmap, err := syscall.Mmap(int(m.file.Fd()), 0, int(t), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// map_array := (*[n]int)(unsafe.Pointer(&mmap[0]))
	// } else {
	header := RegisterHeader{}

	err1 := binary.Read(m.bufReaderIdx, m.endianness, &header)

	if err1 != nil {
		log.Fatalln("binary.Read failed reading data64:", err1)
	}

	(*data) = make([]uint64, header.DataLen, header.DataLen)

	err2 := binary.Read(m.bufReaderDta, m.endianness, data)

	if err2 != nil {
		log.Fatalln("binary.Read failed reading data64:", err2)
	}

	sumDataV := uint64(0)
	for _, w := range *data {
		sumDataV += w
	}

	if header.SumData != sumDataV {
		log.Fatalln("binary.Read failed reading data32: checksum error", header.SumData, sumDataV)
	}
	return header.HasData, header.Serial
	// }
	return false, 0
}

//
// Close
//

func (m *MultiArrayFile) Close() {
	defer m.fileIdx.Close()
	defer m.fileDta.Close()

	if m.writeMode {
		header := RegisterHeader{
			HasData:     false,
			Serial:      0,
			CounterBits: 0,
			DataLen:     0,
			SumData:     0,
		}

		err1 := binary.Write(m.bufWriterIdx, m.endianness, &header)

		if err1 != nil {
			log.Fatalln("binary.Read failed closing file:", err1)
		}

		if m.counterBits == 16 {
			data := make([]int16, m.dataLen, m.dataLen)
			err1 = binary.Write(m.bufWriterDta, m.endianness, data)
		} else if m.counterBits == 32 {
			data := make([]int32, m.dataLen, m.dataLen)
			err1 = binary.Write(m.bufWriterDta, m.endianness, data)
		} else if m.counterBits == 64 {
			data := make([]int64, m.dataLen, m.dataLen)
			err1 = binary.Write(m.bufWriterDta, m.endianness, data)
		}

		if err1 != nil {
			log.Fatalln("binary.Read failed closing file:", err1)
		}

		m.bufWriterIdx.Flush()
		m.bufWriterDta.Flush()
	}
}
