package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// AES 加密key，长度须为 16/24/32 位
const EncryptKey = ""

// license flag register
var licensePath = flag.String("license", "", "license file path")

type Time time.Time

func (c *Time) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`) //get rid of "
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.ParseInLocation(time.DateTime, value, time.Local) //parse time
	if err != nil {
		return err
	}
	*c = Time(t) //set result using the pointer
	return nil
}

func (c Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(c).Format(time.DateTime) + `"`), nil
}

type License struct {
	ID        string `json:"id"`         // 许可ID
	Issuer    string `json:"issuer"`     // 许可颁发人
	Type      string `json:"type"`       // 许可类型，试用（trial)，正式（formal)
	User      string `json:"user"`       // 许可用户
	ExpiredAt *Time  `json:"expired_at"` // 过期时间
	IssuedAt  Time   `json:"issued_at"`  // 许可颁发时间
}

// 许可是否有效
func (l License) Validate() error {
	issuedAt := time.Time(l.IssuedAt)
	if issuedAt.After(time.Now()) {
		return fmt.Errorf("license not yet in effect")
	}

	if l.ExpiredAt != nil {
		// 非永不过期
		expiredAt := time.Time(*l.ExpiredAt)
		if expiredAt.Before(time.Now()) {
			return fmt.Errorf("license expired")
		}
	}

	return nil
}

// 解析许可文件
func ParseLicenseFile(file string) (*License, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("open license file error: %v", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("read license file error: %v", err)
	}

	// 解密
	content, err := AesDecryptByCBC([]byte(EncryptKey), data)
	if err != nil {
		return nil, fmt.Errorf("decrypt license file error: %v", err)
	}

	license := &License{}
	err = json.Unmarshal(content, license)
	if err != nil {
		return nil, fmt.Errorf("invalid license format")
	}

	return license, nil
}

// 解密
func AesDecryptByCBC(key, data []byte) ([]byte, error) {
	keyLen := len(key)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return nil, fmt.Errorf("encrypt key length must be 16/24/32 bit")
	}

	// encrypted密文反解base64
	decodeData, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, fmt.Errorf("decode data error: %v", err)
	}

	// 创建一个cipher.Block接口。参数key为密钥，长度只能是16、24、32字节
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("new aes key error: %v", err)
	}

	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 选择加密模式
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	// 创建数组，存储解密结果
	decodeResult := make([]byte, len(decodeData))
	// 解密
	blockMode.CryptBlocks(decodeResult, decodeData)
	// 解码
	padding := PKCS7UnPadding(decodeResult)
	return padding, nil
}

// 解码
func PKCS7UnPadding(data []byte) []byte {
	length := len(data)
	unPaddingLength := int(data[length-1])
	return data[:(length - unPaddingLength)]
}
