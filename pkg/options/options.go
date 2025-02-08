package options

var DefaultOptions = SymspellOptions{
	MaxDictionaryEditDistance: 2,
	PrefixLength:              7,
	CountThreshold:            1,
	SplitItemThreshold:        1,
	PreserveCase:              false,
	SplitWordBySpace:          false,
	SplitWordAndNumber:        false,
	MinimumCharacterToChange:  1,
}

type SymspellOptions struct {
	MaxDictionaryEditDistance int
	PrefixLength              int
	CountThreshold            int
	SplitItemThreshold        int
	PreserveCase              bool
	SplitWordBySpace          bool
	SplitWordAndNumber        bool
	MinimumCharacterToChange  int
}

type Options interface {
	Apply(options *SymspellOptions)
}

type FuncConfig struct {
	ops func(options *SymspellOptions)
}

func (w FuncConfig) Apply(conf *SymspellOptions) {
	w.ops(conf)
}

func NewFuncOption(f func(options *SymspellOptions)) *FuncConfig {
	return &FuncConfig{ops: f}
}

func WithMaxDictionaryEditDistance(maxDictionaryEditDistance int) Options {
	return NewFuncOption(func(options *SymspellOptions) {
		options.MaxDictionaryEditDistance = maxDictionaryEditDistance
	})
}

func WithPrefixLength(prefixLength int) Options {
	return NewFuncOption(func(options *SymspellOptions) {
		options.PrefixLength = prefixLength
	})
}

func WithCountThreshold(countThreshold int) Options {
	return NewFuncOption(func(options *SymspellOptions) {
		options.CountThreshold = countThreshold
	})
}

func WithSplitItemThreshold(splitThreshold int) Options {
	return NewFuncOption(func(options *SymspellOptions) {
		options.SplitItemThreshold = splitThreshold
	})
}

func WithPreserveCase() Options {
	return NewFuncOption(func(options *SymspellOptions) {
		options.PreserveCase = true
	})
}

func WithSplitWordBySpace() Options {
	return NewFuncOption(func(options *SymspellOptions) {
		options.SplitWordBySpace = true
	})
}

func WithMinimumCharacterToChange(charLength int) Options {
	return NewFuncOption(func(options *SymspellOptions) {
		options.MinimumCharacterToChange = charLength
	})
}

func WithSplitWordAndNumbers() Options {
	return NewFuncOption(func(options *SymspellOptions) {
		options.SplitWordAndNumber = true
	})
}
