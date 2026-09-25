package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	URL       string `yaml:"url"`
	Token     string `yaml:"token"`
	ProjectID int    `yaml:"project_id"`
}

func loadConfig(path string) (Config, error) {
	var c Config
	data, err := os.ReadFile(path)
	if err != nil {
		return c, err
	}
	if err := yaml.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("%s: %w", path, err)
	}
	c.URL = strings.TrimRight(c.URL, "/")
	if c.URL == "" || c.Token == "" || c.ProjectID == 0 {
		return c, fmt.Errorf("%s: url, token и project_id обязательны", path)
	}
	return c, nil
}

// createTask создаёт задачу в проекте и возвращает ссылку на неё.
func createTask(c Config, title, description string, priority int) (string, error) {
	body, err := json.Marshal(map[string]any{
		"title":       title,
		"description": description,
		"priority":    priority,
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/api/v1/projects/%d/tasks", c.URL, c.ProjectID)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("vikunja %s: %s", resp.Status, respBody)
	}

	var created struct {
		ID int `json:"id"`
	}

	if err := json.Unmarshal(respBody, &created); err != nil || created.ID == 0 {
		return "", fmt.Errorf("vikunja: нет id в ответе: %s", respBody)
	}

	return fmt.Sprintf("%s/tasks/%d", c.URL, created.ID), nil
}
