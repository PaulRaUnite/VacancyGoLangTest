package util

import (
	"io"
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	strs := []string{
		"substring texttext",
		"text substring text",
		"texttext substring",
	}
	subs := []struct {
		str string
		res bool
	}{
		{"", false},
		{"substring", true},
		{"text", true},
		{"blablablablablablablablablabla", false},
	}
	for _, s := range strs {
		r := strings.NewReader(s)
		for _, sub := range subs {
			r.Seek(0, io.SeekStart)
			ok := Search(r, sub.str)
			if ok != sub.res {
				t.Fatal("can't find", sub, "in", s)
			}
		}
	}
}

func BenchmarkSearch(b *testing.B) {
	const text = "こんいちは！私はポールです。わたしは１８せいです。よろしくね。ありがとございます。"
	const searchText = "ありがと"

	var rs io.ReadSeeker = strings.NewReader(text)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		found := Search(rs, searchText)
		if found != true {
			b.Fatal("search func doesn't works")
		}
		_, err := rs.Seek(0, io.SeekStart)
		if err != nil {
			b.Fatal(err)
		}
	}
}
