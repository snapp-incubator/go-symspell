package normalization

import "github.com/snapp-incubator/seperno/internal"

func NewNormalize(ops ...internal.Options) *internal.Normalize {
	opts := internal.DefaultOptions
	for _, config := range ops {
		config.Apply(&opts)
	}
	return internal.NewNormalizer(opts)
}

func WithConvertHalfSpaceToSpace() internal.Options {
	return internal.NewFuncWireOption(func(options *internal.NormalizerOptions) {
		options.ConvertHalfSpaceToSpace = true
	})
}

func WithSpaceCombiner() internal.Options {
	return internal.NewFuncWireOption(func(options *internal.NormalizerOptions) {
		options.SpaceCombiner = true
	})
}

func WithOuterSpaceRemover() internal.Options {
	return internal.NewFuncWireOption(func(options *internal.NormalizerOptions) {
		options.OuterSpaceRemover = true
	})
}

func WithURLRemover() internal.Options {
	return internal.NewFuncWireOption(func(options *internal.NormalizerOptions) {
		options.URLRemover = true
	})
}

func WithNormalizePunctuations() internal.Options {
	return internal.NewFuncWireOption(func(options *internal.NormalizerOptions) {
		options.NormalizePunctuations = true
	})
}

func WithEndsWithEndOfLineChar() internal.Options {
	return internal.NewFuncWireOption(func(options *internal.NormalizerOptions) {
		options.EndsWithEndOfLineChar = true
	})
}
