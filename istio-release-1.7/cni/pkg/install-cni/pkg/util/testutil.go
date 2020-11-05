package util

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func CopyExistingConfFiles(t *testing.T, targetDir string, confFiles ...string) {
	t.Helper()
	for _, f := range confFiles {
		data, err := ioutil.ReadFile(filepath.Join("testdata", f))
		if err != nil {
			t.Fatal(err)
		}
		err = ioutil.WriteFile(filepath.Join(targetDir, f), data, 0644)
		if err != nil {
			t.Fatal(err)
		}
	}
}
