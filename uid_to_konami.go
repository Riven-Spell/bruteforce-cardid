package main

import (
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/Riven-Spell/generic/enumerable"
)

func checksum(buf []byte) byte {
	chk := 0
	for i := range 15 {
		chk += int(buf[i]) * (i%3 + 1)
	}

	for chk > 31 {
		chk = (chk >> 5) + (chk & 31)
	}

	return byte(chk)
}

// ported from https://gitea.tendokyu.moe/eamuse/eaapi/src/branch/master/eaapi/cardconv.py#L30
func UIDToKonami(in string, decrypter cipher.BlockMode) (string, error) {
	if len(in) != 16 {
		return "", fmt.Errorf("invalid card ID length %d expected 16", len(in))
	}

	// grab the card type from the prefix
	cardType := byte(0)
	switch {
	case strings.HasPrefix(in, "E004"):
		cardType = 1
	case strings.HasPrefix(in, "0"):
		cardType = 2
	default:
		return "", fmt.Errorf("invalid UID prefix, expected 0 (FeliCa) or E004")
	}

	// decode the hexadecimal to bytes
	uidBytes, err := hex.DecodeString(in)
	if err != nil {
		return "", err
	}
	if len(uidBytes) != 8 {
		return "", errors.New("expected 8 bytes after hex decoding")
	}

	// reverse it first
	slices.Reverse(uidBytes)

	// encrypt it
	cryptDest := make([]byte, 8)
	decrypter.CryptBlocks(cryptDest, uidBytes)

	// proceed to unpack it
	buf, err := Unpack(cryptDest, 5)
	if err != nil {
		return "", fmt.Errorf("failed to unpack: %w", err)
	}

	buf = buf[:13]
	buf = append(buf, 0, 0, 0) // append up to 16 bytes

	buf[0] ^= cardType
	buf[13] = 1
	for i := range 13 {
		buf[i+1] ^= buf[i]
	}
	buf[14] = cardType
	buf[15] = checksum(buf)

	return enumerable.ExecChain1to1[byte, string](
		enumerable.FromList(buf, false),
		enumerable.ChainMap(func(i byte) byte {
			return CardAlphabet[i]
		}),
		enumerable.ChainCollect[byte](),
		enumerable.ChainFuncRaw(func(in []byte) string {
			return string(in)
		}),
	), nil
}
