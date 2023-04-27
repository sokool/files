package files_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/sokool/files"
)

func TestNewStorage(t *testing.T) {
	if _, err := files.NewStorage("s3://key:secret@hostname:9999/bucket_name"); err != nil {
		t.Fatal(err)
	}
	if _, err := files.NewStorage("memory://"); err != nil {
		t.Fatal(err)
	}
	if _, err := files.NewStorage("http://wp.pl"); err == nil {
		t.Fatal()
	}
}

func TestService_ReadWrite(t *testing.T) {
	s, err := files.NewService("memory://")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(s.Files("?!"))
	b := "hello word"
	b1 := bytes.NewBufferString(b)
	if err = s.Write("/test/file", b1); err != nil {
		t.Fatal(err)
	}
	b2 := bytes.NewBuffer(nil)
	if err = s.Read("/test/file", b2); err != nil {
		t.Fatalf("%v and %d", err, b2.Len())
	}
	if b2.String() != b {
		t.Fatal()
	}

}
