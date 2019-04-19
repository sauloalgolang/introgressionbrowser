package save

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	// "io/ioutil"
	"log"
	"os"
)

// https://golang.org/pkg/encoding/binary/
// https://varunpant.com/posts/reading-and-writing-binary-files-in-go-lang

// Example payload
type payload struct {
	One   float32
	Two   float64
	Three uint32
}

//
// Read
//

func readFile() {
	file, err := os.Open("res/test.bin")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	m := payload{}
	for i := 0; i < 10; i++ {
		data := readNextBytes(file, 16)
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &m)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}

		fmt.Println(m)
	}
}

func readNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

//
// Write
//

func writeFile() {
	file, err := os.Create("res/test.bin")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		s := &payload{
			2.0,
			3.0,
			1,
		}
		var bin_buf bytes.Buffer
		binary.Write(&bin_buf, binary.BigEndian, s)
		//b :=bin_buf.Bytes()
		//l := len(b)
		//fmt.Println(l)
		writeNextBytes(file, bin_buf.Bytes())
	}
}

func writeNextBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}

func ArrayWriter(v interface{}) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, v)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	fmt.Printf("% x", buf.Bytes())
}

//
// MultiArrayFile
//

type MultiArrayFile struct {
	fileName   string
	endianness binary.ByteOrder
	buf        *bytes.Buffer
	serial     int64
	isFinished bool
	writeMode  bool
	bufReader  *bufio.Reader
	bufWriter  *bufio.Writer
	file       *os.File
}

func NewMultiArrayFile(fileName string, mode string) *MultiArrayFile {
	m := MultiArrayFile{
		fileName:   fileName,
		endianness: binary.LittleEndian,
		buf:        new(bytes.Buffer),
		serial:     0,
		isFinished: false,
	}

	if mode == "w" {
		log.Println("Saving binary matrix to", fileName)

		file, err := os.Create(fileName)
		if err != nil {
			log.Fatalln(err)
		}

		m.writeMode = true
		m.bufWriter = bufio.NewWriter(file)
		m.file = file

	} else if mode == "r" {
		log.Println("Loading binary matrix from", fileName)

		file, err := os.Open(fileName)
		if err != nil {
			log.Fatalln(err)
		}

		m.writeMode = false
		m.bufReader = bufio.NewReader(file)
		m.file = file

	} else {
		log.Fatalf("invalid mode '%s'. wither w or r\n", mode)
	}

	return &m
}

//
// MultiArrayFile :: Writer
//

func (m *MultiArrayFile) write(dataLen int64) {
	if !m.writeMode {
		log.Fatalln("Trying to write to a reader")
	}

	hasData := true
	err1 := binary.Write(m.bufWriter, m.endianness, &hasData)

	if err1 != nil {
		log.Fatalln("binary.Write failed to write hasData:", err1)
	}

	err2 := binary.Write(m.bufWriter, m.endianness, &m.serial)

	if err2 != nil {
		log.Fatalln("binary.Write failed to write serial:", err2)
	}

	err3 := binary.Write(m.bufWriter, m.endianness, &dataLen)

	if err3 != nil {
		log.Fatalln("binary.Write failed to write dataLen:", err3)
	}

	m.serial++
}

func (m *MultiArrayFile) Write16(data *[]uint16) {
	m.write(int64(len(*data)))

	err := binary.Write(m.bufWriter, m.endianness, *data)

	if err != nil {
		log.Fatalln("binary.Write failed to write data16:", err)
	}
}

func (m *MultiArrayFile) Write32(data *[]uint32) {
	m.write(int64(len(*data)))

	err := binary.Write(m.bufWriter, m.endianness, *data)

	if err != nil {
		log.Fatalln("binary.Write failed to write data32:", err)
	}
}

func (m *MultiArrayFile) Write64(data *[]uint64) {
	m.write(int64(len(*data)))

	err := binary.Write(m.bufWriter, m.endianness, *data)

	if err != nil {
		log.Fatalln("binary.Write failed to write data64:", err)
	}
}

//
// MultiArrayFile :: Reader
//

func (m *MultiArrayFile) read() (hasData bool, dataLen int64) {
	if m.writeMode {
		log.Fatalln("Trying to read from a writer")
	}
	if m.isFinished {
		log.Fatalln("Trying to read a finished file")
	}

	dataLen = int64(0)
	hasData = false
	serial := int64(0)

	err1 := binary.Read(m.bufReader, m.endianness, &hasData)

	if err1 != nil {
		log.Fatalln("binary.Read failed reading hasData:", err1)
	}

	if !hasData {
		m.isFinished = true
		return hasData, dataLen
	}

	err2 := binary.Read(m.bufReader, m.endianness, &serial)

	if err2 != nil {
		log.Fatalln("binary.Read failed reading serial:", err2)
	}

	if serial < 0 {
		log.Fatalln("serial < 0", serial)
	}

	if serial != m.serial {
		log.Fatalln("serial out of order", serial, " != ", m.serial)
	}

	err3 := binary.Read(m.bufReader, m.endianness, &dataLen)

	if err3 != nil {
		log.Fatalln("binary.Read failed reading dataLen:", err3)
	}

	if dataLen <= 0 {
		log.Fatalln("Length <= 0", dataLen)
	}

	m.serial++

	return hasData, dataLen
}

func (m *MultiArrayFile) Read16(data *[]uint16) (hasData bool) {
	dataLen := int64(0)

	hasData, dataLen = m.read()

	*data = make([]uint16, dataLen, dataLen)

	err := binary.Read(m.bufReader, m.endianness, data)

	if err != nil {
		log.Fatalln("binary.Read failed reading data16:", err)
	}

	return hasData
}

func (m *MultiArrayFile) Read32(data *[]uint32) (hasData bool) {
	dataLen := int64(0)

	hasData, dataLen = m.read()

	*data = make([]uint32, dataLen, dataLen)

	err := binary.Read(m.bufReader, m.endianness, data)

	if err != nil {
		log.Fatalln("binary.Read failed reading data32:", err)
	}

	return hasData
}

func (m *MultiArrayFile) Read64(data *[]uint64) (hasData bool) {
	dataLen := int64(0)

	hasData, dataLen = m.read()

	*data = make([]uint64, dataLen, dataLen)

	err := binary.Read(m.bufReader, m.endianness, data)

	if err != nil {
		log.Fatalln("binary.Read failed reading data64:", err)
	}

	return hasData
}

//
// Close
//

func (m *MultiArrayFile) Close() {
	if m.writeMode {
		err := binary.Write(m.bufWriter, m.endianness, false)

		if err != nil {
			log.Fatalln("binary.Write failed:", err)
		}

		m.bufWriter.Flush()
	}

	m.file.Close()
}
