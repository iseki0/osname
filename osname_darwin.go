//go:build darwin

package osname

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

func osname() (string, error) {
	f, e := os.Open("/System/Library/CoreServices/SystemVersion.plist")
	if e != nil {
		return "", _ReadPlistFailed.Cause(e)
	}
	defer f.Close()
	d, e := io.ReadAll(io.LimitReader(f, 1024*1024*4))
	if e != nil {
		return "", _ReadPlistFailed.Cause(e)
	}
	var p plist
	if e := xml.Unmarshal(d, &p); e != nil {
		return "", _ReadPlistFailed.Cause(e)
	}
	var id = 0
	for i, s := range p.Dict.Key {
		if s == "ProductVersion" {
			id = i
		}
	}
	if id < len(p.Dict.String) {
		return fmt.Sprintf("macOS %s", p.Dict.String[id]), nil
	}
	return "", _ReadPlistFailed
}

type plist struct {
	XMLName xml.Name `xml:"plist"`
	Dict    struct {
		Key    []string `xml:"key"`
		String []string `xml:"string"`
	} `xml:"dict"`
}
