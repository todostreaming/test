package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"bytes"
	"fmt"
)

/**
 *	PKCS7补码
 *	这里可以参考下http://blog.studygolang.com/167.html
 */
func PKCS7Padding(data []byte) []byte {
	blockSize := 16
	padding := blockSize - len(data)%blockSize
    	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    	return append(data, padtext...)

}

/**
 *	去除PKCS7的补码
 */
func UnPKCS7Padding(data []byte) []byte {
	length := len(data)
     	// 去掉最后一个字节 unpadding 次
     	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

func main() {

	key := "mofeimofeimofeimofeimofeimofeimo"
	ckey, err := aes.NewCipher([]byte(key))
	if nil != err {
		fmt.Println("钥匙创建错误:", err)
	}

	str := []byte(`polaris@studygo`)
	str = append(str, 0x02)
	iv := []byte("1234567890123456")
	fmt.Println("加密的字符串", string(str), "\n加密钥匙", key, "\n向量IV", string(iv))

	fmt.Println("加密前的字节：", str, "\n")

	encrypter := cipher.NewCBCEncrypter(ckey, iv)

	// PKCS7补码
	str = PKCS7Padding(str)
	out := make([]byte, len(str))

	encrypter.CryptBlocks(out, str)
	fmt.Println("加密后字节：", out)

	base64Out := base64.URLEncoding.EncodeToString(out)
	fmt.Println("Base64后：", base64Out)

	fmt.Println("\n开始解码")
	decrypter := cipher.NewCBCDecrypter(ckey, iv)
	base64In, _ := base64.URLEncoding.DecodeString(base64Out)
	in := make([]byte, len(base64In))
	decrypter.CryptBlocks(in, base64In)

	fmt.Println("解密后的字节：", in)

	// 去除PKCS7补码
	in = UnPKCS7Padding(in)

	fmt.Println("去PKCS7补码：", in)
	fmt.Println("解密：", string(in))
}
