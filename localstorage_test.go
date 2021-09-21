package localstorage

import (
	"os"
	"strconv"
	"testing"
)

func TestLocalStorage(t *testing.T) {
	path := "test_storage"
	lc, err := NewLocalStorage(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		lc.Close()

		// clear
		os.RemoveAll(path)
	}()

	want := 111
	lc.SetItem("a", 111)

	ab, err := lc.GetItem("a")
	if err != nil {
		t.Fatal(err)
	}

	a, err := strconv.ParseInt(string(ab), 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if a != int64(want) {
		t.Fatalf("LocalStorage.GetItem() = %d, want: %d", a, want)
	}

	if !lc.Exists("a") {
		t.Fatal("LocalStorage.Exists() = false, want: true")
	}

	lc.RemoveItem("a")

	if lc.Exists("a") {
		t.Fatal("LocalStorage.Exists() = true, want: false")
	}
}
