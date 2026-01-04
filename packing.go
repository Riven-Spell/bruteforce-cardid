package main

import (
	"errors"
)

// masks is a generated table of congruous masks. 0 is included to 1-index the table.
var masks = func() [9]byte {
	out := [9]byte{}

	n := byte(0)
	for idx := range 8 {
		n <<= 1
		n |= 1
		out[idx+1] = n << (8 - idx - 1) // invert for left hand
	}

	return out
}()

func Unpack(input []byte, width uint8) ([]byte, error) {
	if width == 0 || width > 8 {
		return nil, errors.New("width must be between 1 and 7 bits")
	}

	out := make([]byte, 1)

	// data about where we are with the current byte
	iByteIdx := 0
	iByteUsed, oByteWrote := uint8(0), uint8(0)
	// data about where we are with the overall data
	consumedBits, availableBits := 0, len(input)*8

	for consumedBits < availableBits {
		// first, if we have used all of our input byte, we need to pull a fresh one.
		if iByteUsed == 8 {
			iByteIdx++
			iByteUsed = 0
		}
		// maybe we also need a new output byte...
		if oByteWrote == width {
			out = append(out, 0)
			oByteWrote = 0
		}

		// We need what's left of our current width,
		remainingNeeded := width - oByteWrote
		// But, we can only take what's left of the read byte.
		remainingNeeded = min(8-iByteUsed, remainingNeeded)

		// Grab our mask for the remaining # of bits. this shouldn't be > 8, ever.
		// offset it by the number of bits we've consumed in the current byte.
		mask := masks[remainingNeeded] >> iByteUsed

		// Copy the data over our mask,
		mask = mask & input[iByteIdx]
		// Make our mask align with the write head.
		shr := 8 - int8(iByteUsed+remainingNeeded)      // first, right align
		shr -= int8(width)                              // push to behind width
		shr += int8(remainingNeeded) + int8(oByteWrote) // then pull to the target point

		if shr > 0 {
			mask >>= shr
		} else if shr < 0 {
			mask <<= -shr
		}

		// then, we write to our output.
		out[len(out)-1] |= mask

		// then, increment the consumed bytes by how much we're reading.
		iByteUsed += remainingNeeded
		oByteWrote += remainingNeeded
		consumedBits += int(remainingNeeded)
	}

	return out, nil
}
