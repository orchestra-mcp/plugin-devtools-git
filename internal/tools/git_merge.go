package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-git/internal/git"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// GitMergeSchema returns the JSON Schema for the git_merge tool.
func GitMergeSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"branch": map[string]any{
				"type":        "string",
				"description": "Branch to merge into the current branch",
			},
			"path": map[string]any{
				"type":        "string",
				"description": "Working directory path (defaults to current directory)",
			},
			"no_ff": map[string]any{
				"type":        "boolean",
				"description": "Create a merge commit even for fast-forward merges",
			},
			"message": map[string]any{
				"type":        "string",
				"description": "Custom merge commit message",
			},
		},
		"required": []any{"branch"},
	})
	return s
}

// GitMerge returns a tool handler that merges branches.
func GitMerge() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "branch"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		branch := helpers.GetString(req.Arguments, "branch")
		path := helpers.GetString(req.Arguments, "path")
		noFF := helpers.GetBool(req.Arguments, "no_ff")
		message := helpers.GetString(req.Arguments, "message")

		args := []string{"merge"}
		if noFF {
			args = append(args, "--no-ff")
		}
		if message != "" {
			args = append(args, "-m", message)
		}
		args = append(args, branch)

		output, err := git.Run(ctx, path, args...)
		if err != nil {
			return helpers.ErrorResult("git_error", err.Error()), nil
		}
		return helpers.TextResult(output), nil
	}
}
