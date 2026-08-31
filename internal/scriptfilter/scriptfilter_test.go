package scriptfilter

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestResponse_Write(t *testing.T) {
	resp := Response{
		Items: []Item{
			{
				Title:    "Check",
				Subtitle: "Inspect for invisible Unicode characters",
				Arg:      "check",
				Valid:    BoolPtr(true),
				Variables: map[string]string{
					"source": "clipboard",
				},
				Mods: map[string]Mod{
					"cmd": {
						Subtitle: "Copy report",
						Arg:      "copy-report",
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := resp.Write(&buf); err != nil {
		t.Fatalf("Write: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, buf.String())
	}

	items, ok := decoded["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("items = %v, want a 1-element array", decoded["items"])
	}
	item := items[0].(map[string]any)
	if item["title"] != "Check" {
		t.Errorf("title = %v, want %q", item["title"], "Check")
	}
	if item["valid"] != true {
		t.Errorf("valid = %v, want true", item["valid"])
	}
	mods, ok := item["mods"].(map[string]any)
	if !ok {
		t.Fatalf("mods missing or wrong type: %v", item["mods"])
	}
	cmd, ok := mods["cmd"].(map[string]any)
	if !ok || cmd["arg"] != "copy-report" {
		t.Errorf("mods.cmd = %v, want arg=copy-report", mods["cmd"])
	}
}

func TestResponse_Write_OmitsEmptyFields(t *testing.T) {
	resp := Response{Items: []Item{{Title: "No findings"}}}

	var buf bytes.Buffer
	if err := resp.Write(&buf); err != nil {
		t.Fatalf("Write: %v", err)
	}

	for _, absent := range []string{"\"arg\"", "\"mods\"", "\"variables\"", "\"icon\"", "\"valid\""} {
		if bytes.Contains(buf.Bytes(), []byte(absent)) {
			t.Errorf("output unexpectedly contains %s: %s", absent, buf.String())
		}
	}
}
