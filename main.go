package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var priorities = map[string]int{"P1": 3, "P2": 2, "P3": 1}

func main() {
	configPath := flag.String("config", "config.yaml", "путь к конфигу")
	filePath := flag.String("file", "tasks.md", "путь к файлу с задачами")
	dryRun := flag.Bool("dry-run", false, "только показать, что будет создано")
	flag.Parse()
	log.SetFlags(0)

	var cfg Config
	var err error
	if !*dryRun {
		if cfg, err = loadConfig(*configPath); err != nil {
			log.Fatal(err)
		}
	}

	data, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(data), "\n")
	tasks := parseTasks(lines)

	for _, t := range tasks {
		title := fmt.Sprintf("T%d: %s", t.Num, t.Title)
		switch {
		case t.Status == "done":
			log.Printf("skip %s: done", title)
			continue
		case t.Vikunja != "":
			log.Printf("skip %s: уже создана %s", title, t.Vikunja)
			continue
		case t.VikunjaLine < 0:
			log.Printf("skip %s: нет строки vikunja:", title)
			continue
		}

		description := t.Body
		if t.Cases != "" {
			description = "cases: " + t.Cases + "\n\n" + description
		}
		priority := priorities[t.Priority]

		if *dryRun {
			log.Printf("[dry-run] %s (priority %d, %d символов описания)", title, priority, len(description))
			continue
		}

		link, err := createTask(cfg, title, description, priority)
		if err != nil {
			log.Fatalf("%s: %v", title, err)
		}

		lines[t.VikunjaLine] = "vikunja: " + link
		if err := os.WriteFile(*filePath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
			log.Fatalf("%s создана (%s), но файл не записан: %v", title, link, err)
		}
		log.Printf("created %s -> %s", title, link)
	}
}
