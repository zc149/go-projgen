package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func initGitAndPushAPI(projectName string) error {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("❌ GITHUB_TOKEN not set")
	}

	owner := os.Getenv("GITHUB_OWNER")

	// repo 생성 요청
	body := map[string]any{
		"name":    projectName,
		"private": false,
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "https://api.github.com/user/repos", bytes.NewReader(data))
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("failed to create repo: %s", resp.Status)
	}

	// git init & push
	if err := run("git", "-C", projectName, "init"); err != nil {
		return err
	}
	if err := run("git", "-C", projectName, "add", "."); err != nil {
		return err
	}
	if err := run("git", "-C", projectName, "commit", "-m", "first commit"); err != nil {
		return err
	}
	if err := run("git", "-C", projectName, "branch", "-M", "main"); err != nil {
		return err
	}
	if err := run("git", "-C", projectName, "remote", "add", "origin",
		fmt.Sprintf("https://github.com/%s/%s.git", owner, projectName)); err != nil {
		return err
	}
	if err := run("git", "-C", projectName, "push", "-u", "origin", "main"); err != nil {
		return err
	}

	return nil
}
