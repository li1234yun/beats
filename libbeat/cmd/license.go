package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

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

	// TODO: 解密

	license := &License{}
	err = json.Unmarshal(data, license)
	if err != nil {
		return nil, fmt.Errorf("invalid license format")
	}

	return license, nil
}
