package hash

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSHA256MatchesPythonHelperAlgorithmFixture(t *testing.T) {
	const text = "hello\n"
	const expected = "5891b5b522d5df086d0ff0b110fbd9d21bb4fc7163af34d08286a2e846f6be03"

	fromReader, err := ReaderSHA256(strings.NewReader(text))
	if err != nil {
		t.Fatalf("ReaderSHA256: %v", err)
	}
	if fromReader != expected {
		t.Fatalf("reader digest = %s, want %s", fromReader, expected)
	}

	path := filepath.Join(t.TempDir(), "fixture.txt")
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	fromFile, err := FileSHA256(path)
	if err != nil {
		t.Fatalf("FileSHA256: %v", err)
	}
	if fromFile != expected {
		t.Fatalf("file digest = %s, want %s", fromFile, expected)
	}
}

func TestFileSHA256MatchesPythonHelperFixtureOutputs(t *testing.T) {
	tests := []struct {
		relative string
		expected string
	}{
		{
			relative: "tests/fixtures/kit-lock.yaml",
			expected: "25c20277b11464ec8083bc8459f701f023229a57fc6ed17a77e1f148013ddfee",
		},
		{
			relative: "tests/fixtures/kit-config.yaml",
			expected: "454b734b25583549364ceb055fe70785b0116619f9277adee8fc4e19ad8cb00c",
		},
		{
			relative: "profiles/standard.yaml",
			expected: "3529eb5d051b5bdb5120c4acffdddc24a1c3c9459e23c6131d65902ea9207f79",
		},
	}
	for _, test := range tests {
		t.Run(test.relative, func(t *testing.T) {
			path := filepath.Join("..", "..", test.relative)
			actual, err := FileSHA256(path)
			if err != nil {
				t.Fatalf("FileSHA256(%s): %v", test.relative, err)
			}
			if actual != test.expected {
				t.Fatalf("digest = %s, want Python helper output %s", actual, test.expected)
			}
		})
	}
}
