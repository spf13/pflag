package pflag

import (
	"strconv"
	"strings"
	"testing"
)

func BenchmarkFlagUsagesShort(b *testing.B) {
	for _, n := range []int{50, 200} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			f := NewFlagSet("x", ContinueOnError)
			for i := 0; i < n; i++ {
				f.String("flag"+strconv.Itoa(i), "", "a short usage")
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if usages := f.FlagUsages(); usages == "" {
					b.Fatal("empty usages")
				}
			}
		})
	}
}

func BenchmarkFlagUsagesWrappedLong(b *testing.B) {
	for _, words := range []int{2000, 8000} {
		b.Run(strconv.Itoa(words), func(b *testing.B) {
			f := NewFlagSet("x", ContinueOnError)
			f.String("mount", "", strings.Repeat("word ", words))
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if usages := f.FlagUsagesWrapped(80); usages == "" {
					b.Fatal("empty usages")
				}
			}
		})
	}
}


