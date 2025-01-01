package symspell

import (
	"fmt"
	"github.com/snapp-incubator/symspell/internal/verbosity"
	"testing"
)

func TestSymspellLookup(t *testing.T) {
	type args struct {
		a               string
		maxEditDistance int
		verbosity       verbosity.Verbosity
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test 1",
			args: args{
				a:               "تحریش",
				maxEditDistance: 3,
				verbosity:       verbosity.Top,
			},
			want: 1,
		},
	}
	symSpell := NewSymSpellWithLoadDictionary("/Users/sepehr/Downloads/symspell/vocab.txt", 0, 1,
		WithCountThreshold(0),
		WithMaxDictionaryEditDistance(3),
		WithPrefixLength(5),
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggests, err :=
				symSpell.Lookup(tt.args.a, tt.args.verbosity, tt.args.maxEditDistance)
			if err != nil {
				t.Errorf("err = %v, want %v", err, nil)
			}
			fmt.Println(suggests)
			//if len(suggests) != tt.want {
			//	t.Errorf("len(suggests) = %v, want %v", len(suggests), tt.want)
			//}
		})
	}
}

func TestSymspellLookupCompound(t *testing.T) {
	type args struct {
		a               string
		maxEditDistance int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test 1",
			args: args{
				a:               "whereis th elove hehad dated forImuch of thepast who ",
				maxEditDistance: 2,
			},
			want: "whereas the glove head dated for much of the past who",
		},
		{
			name: "Test 2",
			args: args{
				a:               "whereis th elove heHAd dated forImuch of thEPast who",
				maxEditDistance: 2,
			},
			want: "whereas the glove head dated for much of the last who",
		},
	}
	symSpell := NewSymSpellWithLoadBigramDictionary("internal/tests/vocab.txt", "internal/tests/vocab_bigram.txt",
		0, 1,
		WithCountThreshold(1),
		WithMaxDictionaryEditDistance(2),
		WithPrefixLength(7),
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggests :=
				symSpell.LookupCompound(tt.args.a, tt.args.maxEditDistance)
			if suggests[0].Term != tt.want {
				fmt.Println(suggests[0].Term)
				t.Errorf("want = %v, got %v", tt.want, suggests[0].Term)
			}
		})
	}
}
