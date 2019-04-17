package interfaces

type Parameters struct {
	BlockSize              uint64
	Chromosomes            string
	Compression            string
	ContinueOnError        bool
	CounterBits            int
	DebugFirstOnly         bool
	DebugMaxRegisterThread int64
	DebugMaxRegisterChrom  int64
	Format                 string
	KeepEmptyBlock         bool
	MaxSnpPerBlock         uint64
	MinSnpPerBlock         uint64
	SourceFile             string
}
