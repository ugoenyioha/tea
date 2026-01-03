// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package files

import (
	stdctx "context"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v3"
)

// CmdFilesCreate creates a new file in the repository
var CmdFilesCreate = cli.Command{
	Name:        "create",
	Aliases:     []string{"add", "new"},
	Usage:       "Create a new file in the repository",
	Description: "Create a new file in the repository via API (commits directly)",
	ArgsUsage:   "<path>",
	Action:      runFilesCreate,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:    "message",
			Aliases: []string{"m"},
			Usage:   "Commit message (required)",
		},
		&cli.StringFlag{
			Name:    "branch",
			Aliases: []string{"b"},
			Usage:   "Branch to create file in (defaults to repo default branch)",
		},
		&cli.StringFlag{
			Name:  "new-branch",
			Usage: "Create a new branch for the commit",
		},
		&cli.StringFlag{
			Name:    "content",
			Aliases: []string{"c"},
			Usage:   "File content (use - for stdin, or provide directly)",
		},
		&cli.StringFlag{
			Name:    "from-file",
			Aliases: []string{"f"},
			Usage:   "Read content from local file",
		},
	}, flags.AllDefaultFlags...),
}

func runFilesCreate(_ stdctx.Context, cmd *cli.Command) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	if !cmd.Args().Present() {
		return fmt.Errorf("must specify a file path")
	}

	filePath := cmd.Args().First()
	message := cmd.String("message")
	branch := cmd.String("branch")
	newBranch := cmd.String("new-branch")
	content := cmd.String("content")
	fromFile := cmd.String("from-file")

	if message == "" {
		message = fmt.Sprintf("Create %s", filePath)
	}

	// Get content
	var fileContent []byte
	var err error

	if fromFile != "" {
		fileContent, err = os.ReadFile(fromFile)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
	} else if content == "-" {
		fileContent, err = io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}
	} else if content != "" {
		fileContent = []byte(content)
	} else {
		return fmt.Errorf("must specify --content or --from-file")
	}

	client := ctx.Login.Client()

	opts := gitea.CreateFileOptions{
		FileOptions: gitea.FileOptions{
			Message:       message,
			BranchName:    branch,
			NewBranchName: newBranch,
		},
		Content: base64.StdEncoding.EncodeToString(fileContent),
	}

	resp, _, err := client.CreateFile(ctx.Owner, ctx.Repo, filePath, opts)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	fmt.Printf("Created %s\n", filePath)
	if resp.Commit != nil {
		fmt.Printf("Commit: %s\n", resp.Commit.SHA)
	}

	return nil
}
