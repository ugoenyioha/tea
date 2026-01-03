// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runs

import (
	"time"

	"code.gitea.io/sdk/gitea"
)

// ActionRun represents a workflow run
type ActionRun struct {
	ID           int64       `json:"id"`
	URL          string      `json:"url"`
	HTMLURL      string      `json:"html_url"`
	DisplayTitle string      `json:"display_title"`
	Path         string      `json:"path"`
	Event        string      `json:"event"`
	RunAttempt   int         `json:"run_attempt"`
	RunNumber    int64       `json:"run_number"`
	HeadSha      string      `json:"head_sha"`
	HeadBranch   string      `json:"head_branch"`
	Status       string      `json:"status"`
	Conclusion   string      `json:"conclusion"`
	StartedAt    time.Time   `json:"started_at"`
	CompletedAt  time.Time   `json:"completed_at"`
	Actor        *gitea.User `json:"actor"`
}

// ActionRunList represents a list of workflow runs
type ActionRunList struct {
	WorkflowRuns []*ActionRun `json:"workflow_runs"`
	TotalCount   int64        `json:"total_count"`
}

// ActionJob represents a job within a workflow run
type ActionJob struct {
	ID          int64         `json:"id"`
	RunID       int64         `json:"run_id"`
	Name        string        `json:"name"`
	Status      string        `json:"status"`
	Conclusion  string        `json:"conclusion"`
	StartedAt   time.Time     `json:"started_at"`
	CompletedAt time.Time     `json:"completed_at"`
	Steps       []*ActionStep `json:"steps"`
}

// ActionStep represents a step within a job
type ActionStep struct {
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Conclusion  string    `json:"conclusion"`
	Number      int64     `json:"number"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`
}

// ActionJobList represents a list of jobs
type ActionJobList struct {
	Jobs       []*ActionJob `json:"jobs"`
	TotalCount int64        `json:"total_count"`
}
