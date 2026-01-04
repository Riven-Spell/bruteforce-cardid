package main

import (
	"github.com/Riven-Spell/generic/enumerable"
)

var (
	CardAlphabet = "0123456789ABCDEFGHJKLMNPRSTUWXYZ"
	CardConvKey  = func() []byte {
		out := []byte("?I'llB2c.YouXXXeMeHaYpy!")

		return enumerable.ExecChain1to1[byte, []byte](
			enumerable.FromList(out, false),
			enumerable.ChainMap[byte, byte](func(b byte) byte {
				return b * 2
			}),
			enumerable.ChainCollect[byte]())
	}()
)
