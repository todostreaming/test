package main

import (
       "fmt"
       "crypto/des"
       "crypto/cipher"
       "os"
       "bytes"
)


func PKCS5Padding(src []byte, blockSize int) []byte {
     padding := blockSize - len(src)%blockSize
     padtext := bytes.Repeat([]byte{byte(padding)}, padding)
     return append(src, padtext...)
}


func PKCS5UnPadding(src []byte) []byte {
     length := len(src)
     unpadding := int(src[length-1])
     return src[:(length - unpadding)]
}


func main() {
    // because we are going to use TripleDES... therefore we Triple it!
    triplekey := "12345678" + "12345678" + "12345678" 
    // you can use append as well if you want


    // plaintext will cause panic: crypto/cipher: input not full blocks
    // IF it is not the correct BlockSize. ( des.BlockSize = 8 bytes )
    // to fix this issue, plaintext will be padded to full blocks
    // and unpadded upon decryption

    plaintext := []byte("Hello World!") // Hello World! = 12 bytes.



    block,err := des.NewTripleDESCipher([]byte(triplekey))

    if err != nil {
        fmt.Printf("%s \n", err.Error())
        os.Exit(1)
    }

    fmt.Printf("%d bytes NewTripleDESCipher key with block size of %d bytes\n", len(triplekey), block.BlockSize)


    ciphertext := []byte("abcdef1234567890")
    iv := ciphertext[:des.BlockSize] // const BlockSize = 8

    // encrypt

    mode := cipher.NewCBCEncrypter(block, iv)

    plaintext = PKCS5Padding(plaintext, block.BlockSize())

    encrypted := make([]byte, len(plaintext))
    mode.CryptBlocks(encrypted, plaintext)
    fmt.Printf("%s encrypt to %x \n", plaintext, encrypted)


    //decrypt

    decrypter := cipher.NewCBCDecrypter(block, iv)
    decrypted := make([]byte, len(plaintext))
    decrypter.CryptBlocks(decrypted, encrypted)

    decrypted = PKCS5UnPadding(decrypted)

    fmt.Printf("%x decrypt to %s\n", encrypted, decrypted)

}
 