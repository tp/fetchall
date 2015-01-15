package fetchall

import (
	"io"
	"testing"

	"appengine/aetest"
)

const bigFileUrl = "http://ipv4.download.thinkbroadband.com/50MB.zip"

func TestLoadingHugefile(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}

	client := Client(c)

	r, err := client.Get(bigFileUrl)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()

	buf := make([]byte, 3*1024*1024) // 3 MB
	totalReadBytes := 0

	for {
		readBytes, err := r.Body.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}

		totalReadBytes += readBytes
	}

	t.Logf("read %d bytes in total", totalReadBytes)

	if totalReadBytes < 32*1024*1024 {
		t.Errorf("expected to read more than 32 MB but just read ~ %d MB", totalReadBytes/1014/1014)
	}
}
