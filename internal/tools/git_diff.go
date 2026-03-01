package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-git/internal/git"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// GitDiffSchema returns the JSON Schema for the git_diff tool.
func GitDiffSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Working directory path (defaults to current directory)",
			},
			"ref1": map[string]any{
				"type":        "string",
				"description": "First ref to compare (e.g. HEAD~1, branch name, commit SHA)",
			},
			"ref2": map[string]any{
				"type":        "string",
				"description": "Second ref to compare",
			},
			"staged": map[string]any{
				"type":        "boolean",
				"description": "Show staged (cached) changes instead of unstaged",
			},
		},
	})
	return s
}

// GitDiff returns a tool handler that shows diffs between refs or working tree.
func GitDiff() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		path := helpers.GetString(req.Arguments, "path")
		ref1 := helpers.GetString(req.Arguments, "ref1")
		ref2 := helpers.GetString(req.Arguments, "ref2")
		staged := helpers.GetBool(req.Arguments, "staged")

		args := []string{"diff"}

		if staged {
			args = append(args, "--cached")
		}

		if ref1 != "" {
			args = append(args, ref1)
		}
		if ref2 != "" {
			args = append(args, ref2)
		}

		output, err := git.Run(ctx, path, args...)
		if err != nil {
			return helpers.ErrorResult("git_error", err.Error()), nil
		}
		if output == "" {
			output = "No differences found"
		}
		return helpers.TextResult(output), nil
	}
}
