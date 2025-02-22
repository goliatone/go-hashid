package version

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/andreyvit/diff"
)

var (
	create = false
	update = false
)

func Test_GetVersion(t *testing.T) {
	Tag = "0.0.0"
	Time = "Time"
	User = "username"

	expected := "0.0.0-Time:username"

	result := GetVersion()
	if result == "" {
		t.Errorf("GetVersion failed, expected %v got empty value", expected)
	}

	if result != expected {
		t.Errorf("GetVersion failed, expected %v got %v", expected, result)
	}
}

func Test_Print(t *testing.T) {

	stdout := new(bytes.Buffer)

	Tag = "0.0.1"
	Time = "2024-10-19T03:29:45Z"
	User = "goliatone"
	Commit = "9ae92b384895797a5b291349eb64434d74a96b81"

	actual := stdout.Bytes()
	expected := LoadFixture(t, "version.txt", actual)

	if err := Print(stdout); err != nil {
		t.Fatalf("unwanted error: %v", err)
	}

	if res := string(expected) == string(stdout.String()); !res {
		d := diff.LineDiff(string(expected), string(stdout.String()))

		fmt.Printf("%s", expected)
		fmt.Printf("\n=========\n\n")
		fmt.Printf("%s", stdout.String())

		t.Errorf("output not as expected:\n%v", d)
	}
}

func LoadFixture(t *testing.T, name string, actual []byte) []byte {
	if update {
		UpdateFixture(t, name, actual)
	}
	g, err := _loadFixtureFile(name)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok || !create {
			t.Fatalf("error reading golden file %s: %v", name, err)
		}
		CreateFixture(t, name)
		LoadFixture(t, name, actual)
	}
	return g
}

func CreateFixture(t *testing.T, name string) {
	path := filepath.Join("../../testdata", name)
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("error creating golden file %s: %v", name, err)
	}
	if err = f.Close(); err != nil {
		t.Fatalf("error closing golden file %s: %v", name, err)
	}
}

func UpdateFixture(t *testing.T, name string, actual []byte) {
	path := filepath.Join("../../testdata", name)
	if err := os.WriteFile(path, actual, 0644); err != nil {
		t.Fatalf("error updating golden file %s: %v", name, err)
	}
}

func _loadFixtureFile(name string) ([]byte, error) {
	path := filepath.Join("../../testdata", name)
	return os.ReadFile(path)
}
