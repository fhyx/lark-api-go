package lark

import (
	"testing"
)

func TestCrypto(t *testing.T) {
	text := "hello world"
	cd := NewCrypto("test key")

	plant, err := cd.DecryptString("P37w+VZImNgPEO1RBhJ6RtKl7n6zymIbEG1pReEzghk=")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("plan: %q", plant)

	if plant != text {
		t.Errorf("mismatch result: %q === %q", plant, text)
	}
}
