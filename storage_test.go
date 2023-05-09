package files_test

import (
	"bytes"
	"testing"

	"git2.gamingtec.com/sportsbook-digitain/kit/files"
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

func TestMeta_Filter(t *testing.T) {
	m := files.Meta{"One": "1", "Two": "2", "Three": "3"}
	n := files.Meta{"One": "jeden", "Two": "dwa"}
	o := m.Map(n)
	b := bytes.Buffer{}
	if o["jeden"] != "1" || o["dwa"] != "2" || len(o) != 2 {
		t.Fatal()
	}
	if _, err := o.WriteTo(&b); err != nil {
		t.Fatal(err)
	}
	if b.String() != `{"dwa":"2","jeden":"1"}` {
		t.Fatal()
	}
}
