package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go-app/internal/task"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) ListTasks(group string) ([]task.Task, error) {
	endpoint := c.BaseURL + "/tasks"
	if group != "" {
		endpoint += "?group=" + url.QueryEscape(group)
	}

	resp, err := c.HTTPClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("server unreachable: %w", err)
	}
	defer resp.Body.Close()

	var tasks []task.Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (c *Client) AddTask(title, group string) (*task.Task, error) {
	payload := task.Task{Title: title, Group: group}
	body, _ := json.Marshal(payload)

	resp, err := c.HTTPClient.Post(c.BaseURL+"/tasks", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("server unreachable: %w", err)
	}
	defer resp.Body.Close()

	var t task.Task
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (c *Client) ToggleTask(id string) (*task.Task, error) {
	req, err := http.NewRequest(http.MethodPatch, c.BaseURL+"/tasks/"+id+"/toggle", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("server unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("task ID %s not found", id)
	}

	var t task.Task
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (c *Client) DeleteTask(id string) error {
	req, err := http.NewRequest(http.MethodDelete, c.BaseURL+"/tasks/"+id, nil)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("server unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("task ID %s not found", id)
	}
	return nil
}

func (c *Client) ListGroups() (map[string]int, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/groups")
	if err != nil {
		return nil, fmt.Errorf("server unreachable: %w", err)
	}
	defer resp.Body.Close()

	var groups map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&groups); err != nil {
		return nil, err
	}
	return groups, nil
}
