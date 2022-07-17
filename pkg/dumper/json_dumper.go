package dumper

import (
	"encoding/json"
	"os"
)

type JsonDumper struct {
	File *os.File
}

func (jd *JsonDumper) Dump(d any) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}

	_, err = jd.File.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func (jd *JsonDumper) FromDump(d any) error {
	var b []byte
	_, err := jd.File.Read(b)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, d)
	if err != nil {
		return err
	}

	return nil
}
