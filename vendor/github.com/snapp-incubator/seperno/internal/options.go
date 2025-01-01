package internal

var DefaultOptions = NormalizerOptions{
	ConvertHalfSpaceToSpace: false,
	URLRemover:              false,
	OuterSpaceRemover:       false,
	SpaceCombiner:           false,
	NormalizePunctuations:   false,
	EndsWithEndOfLineChar:   false,
}

type NormalizerOptions struct {
	ConvertHalfSpaceToSpace bool
	URLRemover              bool
	OuterSpaceRemover       bool
	SpaceCombiner           bool
	NormalizePunctuations   bool
	EndsWithEndOfLineChar   bool
}

type Options interface {
	Apply(options *NormalizerOptions)
}

type FuncConfig struct {
	ops func(options *NormalizerOptions)
}

func (w FuncConfig) Apply(conf *NormalizerOptions) {
	w.ops(conf)
}

func NewFuncWireOption(f func(options *NormalizerOptions)) *FuncConfig {
	return &FuncConfig{ops: f}
}
