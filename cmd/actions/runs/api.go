// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runs

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"code.gitea.io/tea/modules/config"
)

// makeAPIRequest makes a direct HTTP request to the Gitea API
// This is needed because the SDK doesn't support workflow runs endpoints
func makeAPIRequest(login *config.Login, method, path string) ([]byte, error) {
	url := login.URL + "/api/v1" + path

	client := &http.Client{}
	if login.Insecure {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "token "+login.Token)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// getWorkflowRuns fetches workflow runs from the API
func getWorkflowRuns(login *config.Login, owner, repo, queryParams string) (*ActionRunList, error) {
	path := fmt.Sprintf("/repos/%s/%s/actions/runs", owner, repo)
	if queryParams != "" {
		path += "?" + queryParams
	}

	body, err := makeAPIRequest(login, "GET", path)
	if err != nil {
		return nil, err
	}

	var result ActionRunList
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// getWorkflowRun fetches a single workflow run
func getWorkflowRun(login *config.Login, owner, repo string, runID int64) (*ActionRun, error) {
	path := fmt.Sprintf("/repos/%s/%s/actions/runs/%d", owner, repo, runID)

	body, err := makeAPIRequest(login, "GET", path)
	if err != nil {
		return nil, err
	}

	var result ActionRun
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// getWorkflowRunJobs fetches jobs for a workflow run
func getWorkflowRunJobs(login *config.Login, owner, repo string, runID int64) (*ActionJobList, error) {
	path := fmt.Sprintf("/repos/%s/%s/actions/runs/%d/jobs", owner, repo, runID)

	body, err := makeAPIRequest(login, "GET", path)
	if err != nil {
		return nil, err
	}

	var result ActionJobList
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}
