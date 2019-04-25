package interfaces

import (
	"fmt"
)

type CallBackParameters struct {
	ContinueOnError bool
	NumBits         int
	NumThreads      int
}

type Parameters struct {
	BlockSize              uint64
	Chromosomes            string
	Compression            string
	ContinueOnError        bool
	CounterBits            int
	DebugFirstOnly         bool
	DebugMaxRegisterThread int64
	DebugMaxRegisterChrom  int64
	Description            string
	Format                 string
	KeepEmptyBlock         bool
	MaxSnpPerBlock         uint64
	MinSnpPerBlock         uint64
	SourceFile             string
}

func (p Parameters) String() (res string) {
	res += fmt.Sprintf("Parameters:\n")
	res += fmt.Sprintf(" BlockSize              : %d\n", p.BlockSize)
	res += fmt.Sprintf(" Chromosomes            : %#v\n", p.Chromosomes)
	res += fmt.Sprintf(" Compression            : %#v\n", p.Compression)
	res += fmt.Sprintf(" ContinueOnError        : %#v\n", p.ContinueOnError)
	res += fmt.Sprintf(" CounterBits            : %#v\n", p.CounterBits)
	res += fmt.Sprintf(" DebugFirstOnly         : %#v\n", p.DebugFirstOnly)
	res += fmt.Sprintf(" DebugMaxRegisterThread : %#v\n", p.DebugMaxRegisterThread)
	res += fmt.Sprintf(" DebugMaxRegisterChrom  : %#v\n", p.DebugMaxRegisterChrom)
	res += fmt.Sprintf(" Description            : %#v\n", p.Description)
	res += fmt.Sprintf(" Format                 : %#v\n", p.Format)
	res += fmt.Sprintf(" KeepEmptyBlock         : %#v\n", p.KeepEmptyBlock)
	res += fmt.Sprintf(" MaxSnpPerBlock         : %d\n", p.MaxSnpPerBlock)
	res += fmt.Sprintf(" MinSnpPerBlock         : %d\n", p.MinSnpPerBlock)
	res += fmt.Sprintf(" SourceFile             : %#v\n", p.SourceFile)
	return res
}
