// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package files

import (
	stdctx "context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"

	"github.com/urfave/cli/v3"
)

// CmdFilesGet gets a file from the repository
var CmdFilesGet = cli.Command{
	Name:        "get",
	Aliases:     []string{"cat", "show", "read"},
	Usage:       "Get a file from the repository",
	Description: "Retrieve the contents of a file from the repository via API",
	ArgsUsage:   "<path>",
	Action:      runFilesGet,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:    "ref",
			Aliases: []string{"b", "branch"},
			Usage:   "Branch, tag, or commit to get file from",
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "Write to file instead of stdout",
		},
		&cli.BoolFlag{
			Name:  "raw",
			Usage: "Output raw file contents without formatting",
		},
	}, flags.AllDefaultFlags...),
}

func runFilesGet(_ stdctx.Context, cmd *cli.Command) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	if !cmd.Args().Present() {
		return fmt.Errorf("must specify a file path")
	}

	filePath := cmd.Args().First()
	ref := cmd.String("ref")
	outputFile := cmd.String("output")
	raw := cmd.Bool("raw")

	client := ctx.Login.Client()

	// Get file contents
	contents, _, err := client.GetContents(ctx.Owner, ctx.Repo, ref, filePath)
	if err != nil {
		return fmt.Errorf("failed to get file: %w", err)
	}

	if contents.Type != "file" {
		return fmt.Errorf("path is a %s, not a file", contents.Type)
	}

	// Decode content (base64)
	if contents.Content == nil {
		return fmt.Errorf("file has no content")
	}

	decoded, err := base64.StdEncoding.DecodeString(*contents.Content)
	if err != nil {
		return fmt.Errorf("failed to decode file content: %w", err)
	}

	// Write to file or stdout
	var writer io.Writer = os.Stdout
	if outputFile != "" {
		f, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer f.Close()
		writer = f
	}

	if raw || outputFile != "" {
		_, err = writer.Write(decoded)
		return err
	}

	// Pretty output with metadata
	if ctx.Output == "json" {
		data, _ := json.MarshalIndent(map[string]interface{}{
			"path":     contents.Path,
			"sha":      contents.SHA,
			"size":     contents.Size,
			"encoding": contents.Encoding,
			"content":  string(decoded),
		}, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	fmt.Printf("Path: %s\n", contents.Path)
	fmt.Printf("SHA:  %s\n", contents.SHA)
	fmt.Printf("Size: %d bytes\n", contents.Size)
	fmt.Println("---")
	fmt.Println(string(decoded))

	return nil
}
