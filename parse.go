package main

import (
	"regexp"
	"strconv"
	"strings"
)

type Task struct {
	Num         int
	Title       string
	Priority    string
	Status      string
	Cases       string
	Vikunja     string
	Body        string
	VikunjaLine int // индекс строки "vikunja:" в файле, для обратной записи
}

var (
	taskHeader = regexp.MustCompile(`^## T(\d+): (.+)$`)
	fieldLine  = regexp.MustCompile(`^(priority|status|cases|vikunja):\s*(.*)$`)
)

// parseTasks находит секции "## T<N>: ..." и разбирает поля priority/status/cases/vikunja.
// Тело задачи - всё после строки "vikunja:" до следующего "## ".
func parseTasks(lines []string) []Task {
	var tasks []Task
	for i := 0; i < len(lines); i++ {
		m := taskHeader.FindStringSubmatch(lines[i])
		if m == nil {
			continue
		}
		t := Task{VikunjaLine: -1}
		t.Num, _ = strconv.Atoi(m[1])
		t.Title = strings.ReplaceAll(m[2], "`", "")

		end := len(lines)
		for j := i + 1; j < len(lines); j++ {
			if strings.HasPrefix(lines[j], "## ") {
				end = j
				break
			}
		}

		bodyStart := i + 1
		for j := i + 1; j < end; j++ {
			f := fieldLine.FindStringSubmatch(lines[j])
			if f == nil {
				continue
			}
			switch f[1] {
			case "priority":
				t.Priority = f[2]
			case "status":
				t.Status = f[2]
			case "cases":
				t.Cases = f[2]
			case "vikunja":
				t.Vikunja = f[2]
				t.VikunjaLine = j
				bodyStart = j + 1
			}
			if t.VikunjaLine >= 0 {
				break
			}
		}
		t.Body = strings.TrimSpace(strings.Join(lines[bodyStart:end], "\n"))
		tasks = append(tasks, t)
		i = end - 1
	}
	return tasks
}
