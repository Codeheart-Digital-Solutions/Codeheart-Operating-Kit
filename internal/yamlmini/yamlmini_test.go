package yamlmini

import (
	"reflect"
	"testing"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/kitfs"
)

func TestParseAndRoundTripStandardProfile(t *testing.T) {
	text, err := kitfs.ReadText("profiles/standard.yaml")
	if err != nil {
		t.Fatalf("read profile: %v", err)
	}
	parsed, err := Parse(text)
	if err != nil {
		t.Fatalf("parse profile: %v", err)
	}
	roundTripped, err := Parse(Dump(parsed))
	if err != nil {
		t.Fatalf("parse dumped profile: %v", err)
	}
	if !reflect.DeepEqual(parsed, roundTripped) {
		t.Fatalf("round trip changed parsed profile\nbefore: %#v\nafter: %#v", parsed, roundTripped)
	}
}

func TestParseCommentsAndQuotedScalars(t *testing.T) {
	parsed, err := MustMap(`
name: "Codeheart # literal"
enabled: true
empty_list: []
count: 7
items:
  -
    path: docs/repo/README.md
    ownership: scaffold
`)
	if err != nil {
		t.Fatalf("parse fixture: %v", err)
	}
	if parsed["name"] != "Codeheart # literal" {
		t.Fatalf("quoted comment scalar changed: %#v", parsed["name"])
	}
	if parsed["enabled"] != true {
		t.Fatalf("bool scalar changed: %#v", parsed["enabled"])
	}
	if parsed["count"] != 7 {
		t.Fatalf("int scalar changed: %#v", parsed["count"])
	}
}

func TestParseCRLFInput(t *testing.T) {
	parsed, err := MustMap("items:\r\n  -\r\n    path: docs/repo/README.md\r\n    ownership: scaffold\r\n")
	if err != nil {
		t.Fatalf("parse CRLF fixture: %v", err)
	}
	items, ok := parsed["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("items = %#v, want one list item", parsed["items"])
	}
	item, ok := items[0].(map[string]any)
	if !ok || item["ownership"] != "scaffold" {
		t.Fatalf("item = %#v, want parsed mapping", items[0])
	}
}
