package main

import (
	"strings"
	"testing"
)

const sample = `# Заголовок

## Группа

## T1: ` + "`foo`" + `: первая

priority: P1
status: done
cases: 1.1
vikunja:

### Где

тело один

## Не задача

текст

## T2: вторая

priority: P2
status: todo
vikunja: https://v/tasks/5

тело два
`

func TestParseTasks(t *testing.T) {
	tasks := parseTasks(strings.Split(sample, "\n"))
	if len(tasks) != 2 {
		t.Fatalf("want 2 tasks, got %d", len(tasks))
	}
	a, b := tasks[0], tasks[1]
	if a.Num != 1 || a.Title != "foo: первая" || a.Priority != "P1" || a.Status != "done" || a.Cases != "1.1" || a.Vikunja != "" {
		t.Errorf("task1 fields: %+v", a)
	}
	if a.Body != "### Где\n\nтело один" {
		t.Errorf("task1 body: %q", a.Body)
	}
	if a.VikunjaLine != 9 {
		t.Errorf("task1 vikunja line: %d", a.VikunjaLine)
	}
	if b.Num != 2 || b.Vikunja != "https://v/tasks/5" || b.Cases != "" || b.Body != "тело два" {
		t.Errorf("task2 fields: %+v", b)
	}
}
