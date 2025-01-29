package symspell

import (
	"fmt"
	"testing"

	"github.com/snapp-incubator/go-symspell/pkg/options"
	"github.com/snapp-incubator/go-symspell/pkg/verbosity"
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
		want string
	}{
		{
			name: "خیابان",
			args: args{
				a:               "حیابان",
				maxEditDistance: 3,
				verbosity:       verbosity.Top,
			},
			want: "خیابان",
		},
		{
			name: "میدان",
			args: args{
				a:               "میذان",
				maxEditDistance: 3,
				verbosity:       verbosity.Top,
			},
			want: "میدان",
		},
		{
			name: "ملاصدرا",
			args: args{
				a:               "ملاصدزا",
				maxEditDistance: 3,
				verbosity:       verbosity.Top,
			},
			want: "ملاصدرا",
		},
		{
			name: "چهاردانگه",
			args: args{
				a:               "چهاردنگه",
				maxEditDistance: 3,
				verbosity:       verbosity.Top,
			},
			want: "چهاردانگه",
		},
	}
	symSpell := NewSymSpellWithLoadDictionary("internal/tests/vocab_fa.txt", 0, 1,
		options.WithCountThreshold(10),
		options.WithMaxDictionaryEditDistance(3),
		options.WithPrefixLength(5),
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggests, err :=
				symSpell.Lookup(tt.args.a, tt.args.verbosity, tt.args.maxEditDistance)
			if err != nil {
				t.Errorf("err = %v, want %v", err, nil)
			}
			if suggests[0].Term != tt.want {
				t.Errorf("got = %v, want %v", suggests[0].Term, tt.want)
			}
		})
	}
}

func TestLookupCompound(t *testing.T) {
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
			want: "where is the love he had dated for much of the past who",
		},
		{
			name: "Test 2",
			args: args{
				a:               "Can yu readthis",
				maxEditDistance: 2,
			},
			want: "can you read this",
		},
		{
			name: "Test 3",
			args: args{
				a:               "sekretplan",
				maxEditDistance: 1,
			},
			want: "secret plan",
		},
	}
	symSpell := NewSymSpellWithLoadBigramDictionary("internal/tests/vocab.txt", "internal/tests/vocab_bigram.txt",
		0, 1,
		options.WithCountThreshold(1),
		options.WithMaxDictionaryEditDistance(3),
		options.WithPrefixLength(7),
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggests :=
				symSpell.LookupCompound(tt.args.a, tt.args.maxEditDistance)
			if suggests.Term != tt.want {
				fmt.Println(suggests.Term)
				t.Errorf("want = %v, got %v", tt.want, suggests.Term)
			}
		})
	}
}

func TestSymspellLookupCompoundUnigram(t *testing.T) {
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
				a:               "میذان ملاصدزا",
				maxEditDistance: 3,
			},
			want: "میدان ملاصدرا",
		}, {
			name: "Test 3",
			args: args{
				a:               "حیابان کارکر",
				maxEditDistance: 3,
			},
			want: "خیابان کارگر",
		}, {
			name: "Test 4",
			args: args{
				a:               "حیابانکارکر",
				maxEditDistance: 3,
			},
			want: "خیابان کارگر",
		},
		{
			name: "Test 5",
			args: args{
				a:               "حیابانملاصدزا",
				maxEditDistance: 3,
			},
			want: "خیابان ملاصدرا",
		},
	}
	symSpell := NewSymSpellWithLoadBigramDictionary("internal/tests/vocab_fa.txt",
		"",
		0,
		1,
		options.WithCountThreshold(10),
		options.WithPrefixLength(5),
		options.WithMaxDictionaryEditDistance(3),
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggest := symSpell.LookupCompound(tt.args.a, tt.args.maxEditDistance)
			if suggest.Term != tt.want {
				t.Errorf("got = %v, want %v", suggest.Term, tt.want)
			}
		})
	}
}
