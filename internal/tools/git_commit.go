package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-git/internal/git"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// GitCommitSchema returns the JSON Schema for the git_commit tool.
func GitCommitSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"message": map[string]any{
				"type":        "string",
				"description": "Commit message",
			},
			"path": map[string]any{
				"type":        "string",
				"description": "Working directory path (defaults to current directory)",
			},
			"files": map[string]any{
				"type":        "array",
				"description": "Specific files to stage before committing",
				"items": map[string]any{
					"type": "string",
				},
			},
			"all": map[string]any{
				"type":        "boolean",
				"description": "Stage all tracked modified files before committing (git add -u)",
			},
		},
		"required": []any{"message"},
	})
	return s
}

// GitCommit returns a tool handler that stages files and creates a commit.
func GitCommit() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "message"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		message := helpers.GetString(req.Arguments, "message")
		path := helpers.GetString(req.Arguments, "path")
		files := helpers.GetStringSlice(req.Arguments, "files")
		all := helpers.GetBool(req.Arguments, "all")

		// Stage files first.
		if len(files) > 0 {
			addArgs := append([]string{"add"}, files...)
			if _, err := git.Run(ctx, path, addArgs...); err != nil {
				return helpers.ErrorResult("git_error", fmt.Sprintf("staging files: %s", err.Error())), nil
			}
		} else if all {
			if _, err := git.Run(ctx, path, "add", "-u"); err != nil {
				return helpers.ErrorResult("git_error", fmt.Sprintf("staging all: %s", err.Error())), nil
			}
		}

		// Create commit.
		output, err := git.Run(ctx, path, "commit", "-m", message)
		if err != nil {
			return helpers.ErrorResult("git_error", err.Error()), nil
		}
		return helpers.TextResult(output), nil
	}
}
