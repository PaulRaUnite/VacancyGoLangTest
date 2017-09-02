package util

import (
	"bufio"
	"io"
)

// Search finds substring in reader, returns false if error happens
// or it can't find the substring in reader.
// Be sure, that your Reader is ReadSeeker
// or you can start reading from it's origin
// (you'll get nothing if you will try to
// read after Search call)
func Search(r io.Reader, text string) bool {
	// It's nothing to do with empty text.
	if len(text) == 0 {
		return false
	}

	// Wrap reader with bufio.Reader to access ReadRune method.
	br := bufio.NewReader(r)
	// match keeps result of comparison in inner loop.
	match := false
	for {
		// The inner loop does two things:
		// - it gets new runes from reader,
		//   and thus goes through all symbols
		//   in reader.
		// - tries to find symbol where
		//   substring starts and when
		//   make sure that it really is
		//   there.
		for _, searchedSymbol := range text {
			readerSymbol, _, err := br.ReadRune()
			if err != nil {
				return false
			}

			match = true
			if searchedSymbol != readerSymbol {
				match = false
				break
			}
		}
		// Match == true means that
		// substring was found.
		if match {
			return true
		}
	}
	return false
}
