package util

import (
	"io"
)

// Search finds substring in reader, returns
// false if error happens or it can't find
// the substring in reader.
//
// Make sure, that your Reader is ReadSeeker
// or you can start reading from it's origin
// (you'll get nothing if you will try to
// read after Search call).
//
// Also, try to use buffers for readers of
// files or net connections: Search reads
// by single byte so bufferising can make
// significant performance speed up.
func Search(r io.Reader, text string) bool {
	// It's nothing to do with empty text.
	if len(text) == 0 {
		return false
	}

	// This construction is faster than make([]byte, 1, 1) or []byte{0}.
	// The single byte is used to read from reader.
	ba := [1]byte{0} //byte array
	bs := ba[:]      //byte slice

	// match keeps result of comparison in inner loop.
	match := false
	for {
		// This inner loop does two things:
		// - it gets new byte from reader,
		//   and thus goes through all bytes
		//   in reader.
		// - tries to find byte sequence
		//   which starts as text
		for i := 0; i < len(text); i++ {
			_, err := r.Read(bs)
			if err != nil {
				return false
			}

			match = true
			if text[i] != bs[0] {
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
