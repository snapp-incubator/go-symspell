package editdistance

import (
	"testing"
)

func TestDamerauLevenshteinDistance(t *testing.T) {
	type args struct {
		a string
		b string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test DamerauLevenshteinDistance",
			args: args{
				a: "kitten",
				b: "sitting",
			},
			want: 3,
		},
		{
			name: "Test DamerauLevenshteinDistance",
			args: args{
				a: "kitten",
				b: "kitten",
			},
			want: 0,
		},
		{
			name: "Test DamerauLevenshteinDistance",
			args: args{
				a: "تحریش",
				b: "تجریش",
			},
			want: 1,
		},
		{
			name: "Test DamerauLevenshteinDistance",
			args: args{
				a: "میدان",
				b: "میذات",
			},
			want: 3,
		},
		{
			name: "Test DamerauLevenshteinDistance",
			args: args{
				a: "کناب",
				b: "کتاب",
			},
			want: 2,
		},
		{
			name: "Test DamerauLevenshteinDistance",
			args: args{
				a: "pipe",
				b: "pips",
			},
			want: 1,
		},
		{
			name: "Test DamerauLevenshteinDistance",
			args: args{
				a: "stream",
				b: "steamb",
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEditDistance(DamerauLevenshtein).Distance(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Distance() = %v, want %v", got, tt.want)
			}
		})
	}
}
