package main

import (
	"bytes"
	"io/ioutil"
)

func readDictionary(filename string) map[string]string {
	dict := make(map[string]string)
	if buf, e := ioutil.ReadFile(filename); e == nil {
		lines := bytes.Split(buf, []byte{'\n'})
		for _, line := range lines {
			kv := bytes.Split(line, []byte{'='})
			if len(kv) != 2 {
				continue
			}
			k := string(bytes.TrimSpace(kv[0]))
			dict[k] = string(bytes.TrimSpace(kv[1]))
		}
	}
	return dict
}
