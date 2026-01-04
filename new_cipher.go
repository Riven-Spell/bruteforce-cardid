package main

import (
	"crypto/cipher"
	"crypto/des"
)

func NewEncrypter() cipher.BlockMode {
	//fmt.Println(hex.EncodeToString(CardConvKey))

	// create a block cipher
	c, err := des.NewTripleDESCipher(CardConvKey)
	if err != nil {
		panic(err)
	}

	// create the encrypter
	bm := cipher.NewCBCEncrypter(c, make([]byte, 8))
	return bm
}
