package main

import (
	"strings"
	"testing"
)

func TestToHTML(t *testing.T) {
	html, err := toHTML("### Где\n\n`a` и **b**\n\n- [ ] один\n- [x] два\n")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"<h3>Где</h3>", "<code>a</code>", "<strong>b</strong>", `type="checkbox"`, "checked"} {
		if !strings.Contains(html, want) {
			t.Errorf("нет %q в:\n%s", want, html)
		}
	}
}
