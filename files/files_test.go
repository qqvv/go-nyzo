package files

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestReadWriteDelete(t *testing.T) {
	fn := "testfile"
	content := "#thestandard\n\nfigisfidis\r"
	if err := Write(fn, []byte(content)); err != nil {
		t.Error(err)
	}

	bytesRead, err := ReadBytes(fn)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal([]byte(content), bytesRead) {
		t.Error(errors.New("written and read content do not match"))
	}

	stringRead, err := ReadString(fn)
	if err != nil {
		t.Error(err)
	}
	if strings.Compare(content, stringRead) != 0 {
		t.Error(errors.New("written and read content do not match"))
	}

	if err = Delete(fn); err != nil {
		t.Error(err)
	}
}
