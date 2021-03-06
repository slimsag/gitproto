// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitproto

import (
	"errors"
	"fmt"
	"strconv"
)

// Note: A unofficial specification for the Git Smart/Dumb HTTP protocols is
// available at:
//
//  https://gist.github.com/schacon/6092633
//

// A pkt-line is a variable length binary string implemented in the Git Smart
// HTTP protocol.
type PktLine []byte

// Bytes returns the binary form of this pkt-line.
func (pl PktLine) Bytes() []byte {
	// The actual hex string (length 4) is included in the hex string.
	hexLen := fmt.Sprintf("%04x", len(pl)+4)
	return append([]byte(hexLen), pl...)
}

var (
	errPktLineNeedMore = errors.New("need more data")
)

// NextPktLine parses the next pkt-line from the given binary data.
//
// If the data provided is not enough then err=errPktLineNeedMore is
// returned.
//
// A special line prefixed with "0000" returns lineBreak=true directly.
//
// The returned integer is the number of bytes of consumed data.
func NextPktLine(data []byte) (pl PktLine, lineBreak bool, n int, err error) {
	// Newlines exist in encoded pkt-lines but they do not serve any real-world
	// purpose (aside from viewing the binary blob using a text editor). The data
	// in the line itself is binary and may include newlines etc inside of it.
	//
	// The only valid way to split the line data is to way operate based on the
	// first four (length) bytes.

	// Need at least four bytes.
	if len(data) < 4 {
		err = errPktLineNeedMore
		return
	}

	// The first four bytes of the line is the total length of the line, in
	// hexadecimal.
	var length uint64
	length, err = strconv.ParseUint(string(data[:4]), 16, 16)
	if err != nil {
		return
	}
	if length == 0 {
		// Special case: line break.
		n = 4
		lineBreak = true
		return
	}
	if int(length) > len(data) {
		err = errPktLineNeedMore
		return
	}
	pl = PktLine(data[4:length])
	n = int(length)
	return
}
