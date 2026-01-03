// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package reactions

import (
	stdctx "context"
	"fmt"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/utils"

	"github.com/urfave/cli/v3"
)

// CmdReactionAdd adds a reaction to an issue or comment
var CmdReactionAdd = cli.Command{
	Name:        "add",
	Aliases:     []string{"a", "+"},
	Usage:       "Add a reaction to an issue or comment",
	Description: "Add an emoji reaction to an issue, PR, or comment.\n" + ReactionHelp,
	ArgsUsage:   "<reaction>",
	Action:      runReactionAdd,
	Flags: append([]cli.Flag{
		&cli.IntFlag{
			Name:    "issue",
			Aliases: []string{"i"},
			Usage:   "Issue or PR index",
		},
		&cli.IntFlag{
			Name:    "comment",
			Aliases: []string{"c"},
			Usage:   "Comment ID (if reacting to a comment)",
		},
	}, flags.AllDefaultFlags...),
}

func runReactionAdd(_ stdctx.Context, cmd *cli.Command) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	if !cmd.Args().Present() {
		return fmt.Errorf("must specify a reaction. %s", ReactionHelp)
	}

	reaction := NormalizeReaction(cmd.Args().First())
	issueIndex := cmd.Int("issue")
	commentID := cmd.Int("comment")

	if issueIndex == 0 && commentID == 0 {
		return fmt.Errorf("must specify --issue or --comment")
	}

	client := ctx.Login.Client()

	if commentID > 0 {
		// React to a comment
		_, _, err := client.PostIssueCommentReaction(ctx.Owner, ctx.Repo, int64(commentID), reaction)
		if err != nil {
			return fmt.Errorf("failed to add reaction to comment: %w", err)
		}
		fmt.Printf("Added %s reaction to comment %d\n", reaction, commentID)
	} else {
		// React to an issue/PR
		index, err := utils.ArgToIndex(fmt.Sprintf("%d", issueIndex))
		if err != nil {
			return err
		}
		_, _, err = client.PostIssueReaction(ctx.Owner, ctx.Repo, index, reaction)
		if err != nil {
			return fmt.Errorf("failed to add reaction to issue: %w", err)
		}
		fmt.Printf("Added %s reaction to issue #%d\n", reaction, issueIndex)
	}

	return nil
}
