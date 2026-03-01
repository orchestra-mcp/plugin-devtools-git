package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-git/internal/git"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// GHPRMergeSchema returns the JSON Schema for the gh_pr_merge tool.
func GHPRMergeSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"pr_number": map[string]any{
				"type":        "number",
				"description": "Pull request number",
			},
			"path": map[string]any{
				"type":        "string",
				"description": "Working directory path (defaults to current directory)",
			},
			"method": map[string]any{
				"type":        "string",
				"description": "Merge method: merge, squash, or rebase (default merge)",
				"enum":        []any{"merge", "squash", "rebase"},
			},
		},
		"required": []any{"pr_number"},
	})
	return s
}

// GHPRMerge returns a tool handler that merges a GitHub pull request.
func GHPRMerge() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		prNumber := helpers.GetInt(req.Arguments, "pr_number")
		if prNumber == 0 {
			return helpers.ErrorResult("validation_error", "pr_number is required"), nil
		}

		path := helpers.GetString(req.Arguments, "path")
		method := helpers.GetStringOr(req.Arguments, "method", "merge")

		if err := helpers.ValidateOneOf(method, "merge", "squash", "rebase"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		args := []string{"pr", "merge", fmt.Sprintf("%d", prNumber)}

		switch method {
		case "squash":
			args = append(args, "--squash")
		case "rebase":
			args = append(args, "--rebase")
		default:
			args = append(args, "--merge")
		}

		output, err := git.GH(ctx, path, args...)
		if err != nil {
			return helpers.ErrorResult("gh_error", err.Error()), nil
		}
		if output == "" {
			output = fmt.Sprintf("PR #%d merged via %s", prNumber, method)
		}
		return helpers.TextResult(output), nil
	}
}
