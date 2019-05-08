package save

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Saver holds the informatio for saving to file
type Saver struct {
	Prefix              string
	Format              string
	Compressor          string
	FormatExtension     string
	CompressorExtension string
	Extension           string
}

// NewSaver creates new instance of Saver
func NewSaver(prefix string, format string) *Saver {
	return newSaver(prefix, format, "none", GetFormatExtension(format), GetCompressExtension("none"))
}

// NewSaverCompressed creates new instance of Saver with compression
func NewSaverCompressed(prefix string, format string, compressor string) *Saver {
	return newSaver(prefix, format, compressor, GetFormatExtension(format), GetCompressExtension(compressor))
}

func newSaver(prefix string, format string, compressor string, formatExtension string, compressorExtension string) *Saver {
	s := Saver{
		Prefix:              prefix,
		Format:              format,
		Compressor:          compressor,
		FormatExtension:     formatExtension,
		CompressorExtension: compressorExtension,
		Extension:           "",
	}

	return &s
}

//
// Setters
//

// SetFormat sets the current format
func (s *Saver) SetFormat(format string) {
	s.Format = format
	s.FormatExtension = GetFormatExtension(format)
}

// SetCompressor sets the compressor
func (s *Saver) SetCompressor(compressor string) {
	s.Compressor = compressor
	s.CompressorExtension = GetFormatExtension(compressor)
}

// SetFormatExtension override format extension
func (s *Saver) SetFormatExtension(extension string) {
	s.FormatExtension = extension
}

// SetCompressorExtension override compression extension
func (s *Saver) SetCompressorExtension(extension string) {
	s.CompressorExtension = extension
}

// SetExtension override the whole extension
func (s *Saver) SetExtension(extension string) {
	s.Extension = extension
}

//
// Getters
//

// GenFilename generates the final filename
func (s *Saver) GenFilename() string {
	outname := s.Prefix

	if s.Extension != "" {
		outname += "." + s.Extension

	} else {
		if s.FormatExtension != "" {
			outname += "." + s.FormatExtension
		}

		if s.CompressorExtension != "" {
			outname += "." + s.CompressorExtension
		}
	}

	return outname
}

// Exists checks whether output file exists
func (s *Saver) Exists() (bool, error) {
	fileName := s.GenFilename()

	_, err := os.Stat(fileName)

	if err == nil {
		// path/to/whatever exists
		return true, err

	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
		return false, err
	}

	return false, err
	// Schrodinger: file may or may not exist. See err for details.
	// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
}

//
// Save
//

// Save saves to file
func (s *Saver) Save(val interface{}) {
	format := s.Format
	compress := s.Compressor

	// fmt.Println("format       ", format)
	// fmt.Println("compress     ", compress)

	hasStreamer := GetFormatHasStreamer(format)
	hasMarshal := GetFormatHasMarshal(format)
	isCompressed := GetCompressIsCompressed(compress)

	// fmt.Println("hasStreamer  ", hasStreamer)
	// fmt.Println("hasMarshal   ", hasMarshal)
	// fmt.Println("isCompressed ", isCompressed)

	outfile := s.GenFilename()

	// fmt.Println("outfile      ", outfile)

	if hasStreamer {
		if isCompressed {
			marshaler := GetFormatMarshalerStreamerWriter(format)
			compressor := GetCompressInterfaceWriter(compress)
			saveDataStreamCompressed(outfile, marshaler, compressor, val)
		} else {
			marshaler := GetFormatMarshalerStreamer(format)
			saveDataStream(outfile, marshaler, val)
		}

	} else if hasMarshal {
		marshaler := GetFormatMarshaler(format)
		saveData(outfile, marshaler, val)
	}
}

func saveData(outfile string, marshaler Marshaler, val interface{}) {
	d, err := marshaler(val)

	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	fmt.Println("saving data to ", outfile)

	err = ioutil.WriteFile(outfile, d, 0644)
	fmt.Println("  done")
}

func saveDataStream(outfile string, marshaler MarshalerStreamer, val interface{}) {
	// fmt.Println("saving stream to ", outfile)

	marshaler(outfile, val)
}

func saveDataStreamCompressed(outfile string, marshaler MarshalerStreamerWriter, compressor GenericNewWriter, val interface{}) error {
	// fmt.Println("saveDataStreamCompressed :: outfile    ", outfile)
	// fmt.Println("saveDataStreamCompressed :: marshaler  ", marshaler)
	// fmt.Println("saveDataStreamCompressed :: compressor ", compressor)

	file, err := os.OpenFile(outfile, os.O_CREATE|os.O_WRONLY, 0660) //|os.O_APPEND
	defer file.Close()

	if err == nil {
		comp := compressor(file)

		marshaler(comp, val)

		comp.Flush()
		comp.Close()

	} else {
		fmt.Println(err)
		os.Exit(1)
	}

	return err
}

//
//
// Load
//
//

// Load loads from file
func (s *Saver) Load(val interface{}) {
	format := s.Format

	outfile := s.GenFilename()

	hasStreamer := GetFormatHasStreamer(format)
	hasMarshal := GetFormatHasMarshal(format)

	if hasStreamer {
		unmarshaler := GetFormatUnMarshalerStreamer(format)
		loadDataStream(outfile, unmarshaler, val)

	} else if hasMarshal {
		unmarshaler := GetFormatUnMarshaler(format)
		loadData(outfile, unmarshaler, val)
	}
}

func loadData(outfile string, unmarshaler UnMarshaler, val interface{}) {
	data, err := ioutil.ReadFile(outfile)

	if err != nil {
		fmt.Printf("dump file. Get err   #%v ", err)
	}

	err = unmarshaler(data, val)

	if err != nil {
		fmt.Printf("cannot unmarshal data: %v", err)
	}
}

func loadDataStream(outfile string, unmarshaler UnMarshalerStreamer, val interface{}) {
	fmt.Println("loading from ", outfile)
	unmarshaler(outfile, val)
}

//
// Guess file format
//

// GuessFormat guesses the file format by its extension
func GuessFormat(filename string) (found bool, format string, compression string, prefix string) {
	found = false
	format = ""
	compression = ""
	prefix = ""

	for fn, fv := range Formats {
		fext := fv.Extension
		for cn, cv := range Compressors {
			cext := cv.Extension
			extension := "." + fext

			if cext != "" {
				extension += "." + cext
			}

			if strings.HasSuffix(filename, extension) {
				found = true
				format = fn
				compression = cn
				prefix = strings.TrimSuffix(filename, extension)
				return found, format, compression, prefix
			}
		}
	}
	return found, format, compression, prefix
}

// GuessPrefixFormat guesses the file format given a file prefix
func GuessPrefixFormat(prefix string) (found bool, format string, compression string, filename string) {
	found = false
	format = ""
	compression = ""
	filename = ""

	for fn, fv := range Formats {
		fext := fv.Extension
		for cn, cv := range Compressors {
			cext := cv.Extension
			filename := prefix + "." + fext

			if cext != "" {
				filename += "." + cext
			}

			_, err := os.Stat(filename)

			if err != nil {
				continue
			} else {
				found = true
				format = fn
				compression = cn
				return found, format, compression, filename
			}
		}
	}

	return found, format, compression, filename
}
