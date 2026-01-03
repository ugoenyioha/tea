// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package comments

import (
	stdctx "context"
	"errors"
	"fmt"
	"io"
	"strings"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/interact"
	"code.gitea.io/tea/modules/print"
	"code.gitea.io/tea/modules/theme"

	"code.gitea.io/sdk/gitea"
	"github.com/charmbracelet/huh"
	"github.com/urfave/cli/v3"
)

// CmdCommentsUpdate updates a comment
var CmdCommentsUpdate = cli.Command{
	Name:        "update",
	Aliases:     []string{"edit", "e"},
	Usage:       "Update a comment",
	Description: "Update an existing comment by ID",
	ArgsUsage:   "<comment id> [<new body>]",
	Action:      runCommentsUpdate,
	Flags:       flags.AllDefaultFlags,
}

func runCommentsUpdate(_ stdctx.Context, cmd *cli.Command) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	args := ctx.Args()
	if args.Len() == 0 {
		return fmt.Errorf("must specify comment ID")
	}

	commentID, err := parseCommentID(args.First())
	if err != nil {
		return err
	}

	body := strings.Join(args.Tail(), " ")
	if interact.IsStdinPiped() {
		if bodyStdin, err := io.ReadAll(ctx.Reader); err != nil {
			return err
		} else if len(bodyStdin) != 0 {
			body = strings.Join([]string{body, string(bodyStdin)}, "\n\n")
		}
	} else if len(body) == 0 {
		// Get existing comment body for editing
		client := ctx.Login.Client()
		existingComment, _, err := client.GetIssueComment(ctx.Owner, ctx.Repo, commentID)
		if err != nil {
			return fmt.Errorf("failed to get existing comment: %w", err)
		}
		body = existingComment.Body

		if err := huh.NewForm(
			huh.NewGroup(
				huh.NewText().
					Title("Edit comment (markdown):").
					ExternalEditor(config.GetPreferences().Editor).
					EditorExtension("md").
					Value(&body),
			),
		).WithTheme(theme.GetTheme()).
			Run(); err != nil {
			return err
		}
	}

	if len(body) == 0 {
		return errors.New("no comment content provided")
	}

	client := ctx.Login.Client()
	comment, _, err := client.EditIssueComment(ctx.Owner, ctx.Repo, commentID, gitea.EditIssueCommentOption{
		Body: body,
	})
	if err != nil {
		return err
	}

	print.Comment(comment)
	return nil
}

func parseCommentID(s string) (int64, error) {
	var id int64
	_, err := fmt.Sscanf(s, "%d", &id)
	if err != nil {
		return 0, fmt.Errorf("invalid comment ID: %s", s)
	}
	return id, nil
}
