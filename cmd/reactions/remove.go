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

// CmdReactionRemove removes a reaction from an issue or comment
var CmdReactionRemove = cli.Command{
	Name:        "remove",
	Aliases:     []string{"rm", "delete", "-"},
	Usage:       "Remove a reaction from an issue or comment",
	Description: "Remove an emoji reaction from an issue, PR, or comment.\n" + ReactionHelp,
	ArgsUsage:   "<reaction>",
	Action:      runReactionRemove,
	Flags: append([]cli.Flag{
		&cli.IntFlag{
			Name:    "issue",
			Aliases: []string{"i"},
			Usage:   "Issue or PR index",
		},
		&cli.IntFlag{
			Name:    "comment",
			Aliases: []string{"c"},
			Usage:   "Comment ID (if removing reaction from a comment)",
		},
	}, flags.AllDefaultFlags...),
}

func runReactionRemove(_ stdctx.Context, cmd *cli.Command) error {
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
		// Remove reaction from a comment
		_, err := client.DeleteIssueCommentReaction(ctx.Owner, ctx.Repo, int64(commentID), reaction)
		if err != nil {
			return fmt.Errorf("failed to remove reaction from comment: %w", err)
		}
		fmt.Printf("Removed %s reaction from comment %d\n", reaction, commentID)
	} else {
		// Remove reaction from an issue/PR
		index, err := utils.ArgToIndex(fmt.Sprintf("%d", issueIndex))
		if err != nil {
			return err
		}
		_, err = client.DeleteIssueReaction(ctx.Owner, ctx.Repo, index, reaction)
		if err != nil {
			return fmt.Errorf("failed to remove reaction from issue: %w", err)
		}
		fmt.Printf("Removed %s reaction from issue #%d\n", reaction, issueIndex)
	}

	return nil
}
