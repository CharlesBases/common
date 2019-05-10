package encode

/*
 AES
	非对称加密，加解密使用不同的密钥

 AES:   高级加密标准（Advanced Encryption Standard），又称 Rijndael 加密法，这个标准用来替代原先的 DES。AES 加密数据块分组长度必须为 128bit(byte[16])，密钥长度可以是 128bit(byte[16])、192bit(byte[24])、256bit(byte[32])中的任意一个。
 块: 对明文进行加密的时候，先要将明文按照 128bit 进行划分。
 填充方式:  因为明文的长度不一定总是 128 的整数倍，所以要进行补位，我们这里采用的是 PKCS7 填充方式。

 AES实现的方式多样，其中包括 ECB，CBC，CTR，CFB，OFB 等。
	1.ECB 电话本模式(Electronic Codebook Book)	将明文分组加密之后的结果直接称为分组
	2.CBC 密码分组链接模式（Cipher Block Chaining） 将明文分组与前一个密文分组进行 XOR 运算，然后再进行加密。每个分组的加解密都依赖于前一个分组。而第一个分组没有前一个分组，因此需要一个初始化向量
	3.CTR 计算器模式（Counter）
	4.CFB 密码反馈模式（Cipher FeedBack） 前一个密文分组会被送回到密码算法的输入端。 在 CBC 和 EBC 模式中，明文分组都是通过密码算法进行加密的。而在 CFB 模式中，明文分组并没有通过加密算法直接进行加密，明文分组和密文分组之间只有一个 XOR。
	5.OFB 输出反馈模式（Output FeedBack）
*/

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func pKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// GoAES 加密
type GoAES struct {
	Key []byte
}

// NewGoAES 返回GoAES
func NewGoAES(key []byte) *GoAES {
	return &GoAES{Key: key}
}

// Encrypt 加密数据
func (a *GoAES) Encrypt(origData []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = pKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, a.Key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// Decrypt 解密数据
func (a *GoAES) Decrypt(crypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, a.Key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = pKCS7UnPadding(origData)
	return origData, nil
}
