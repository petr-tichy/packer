package common

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os"
	"testing"
)

func TestDownloadClient_VerifyChecksum(t *testing.T) {
	tf, err := ioutil.TempFile("", "packer")
	if err != nil {
		t.Fatalf("tempfile error: %s", err)
	}
	defer os.Remove(tf.Name())

	// "foo"
	checksum, err := hex.DecodeString("acbd18db4cc2f85cedef654fccc4a4d8")
	if err != nil {
		t.Fatalf("decode err: %s", err)
	}

	// Write the file
	tf.Write([]byte("foo"))
	tf.Close()

	config := &DownloadConfig{
		Hash:     md5.New(),
		Checksum: checksum,
	}

	d := NewDownloadClient(config)
	result, err := d.VerifyChecksum(tf.Name())
	if err != nil {
		t.Fatalf("Verify err: %s", err)
	}

	if !result {
		t.Fatal("didn't verify")
	}
}

func TestHashForType(t *testing.T) {
	if h := HashForType("md5"); h == nil {
		t.Fatalf("md5 hash is nil")
	} else {
		h.Write([]byte("foo"))
		result := h.Sum(nil)

		expected := "acbd18db4cc2f85cedef654fccc4a4d8"
		actual := hex.EncodeToString(result)
		if actual != expected {
			t.Fatalf("bad hash: %s", actual)
		}
	}

	if HashForType("fake") != nil {
		t.Fatalf("fake hash is not nil")
	}
}
