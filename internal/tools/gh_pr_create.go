package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-git/internal/git"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// GHPRCreateSchema returns the JSON Schema for the gh_pr_create tool.
func GHPRCreateSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"title": map[string]any{
				"type":        "string",
				"description": "Pull request title",
			},
			"body": map[string]any{
				"type":        "string",
				"description": "Pull request body/description",
			},
			"path": map[string]any{
				"type":        "string",
				"description": "Working directory path (defaults to current directory)",
			},
			"base": map[string]any{
				"type":        "string",
				"description": "Base branch to merge into (defaults to repo default branch)",
			},
			"head": map[string]any{
				"type":        "string",
				"description": "Head branch to create PR from (defaults to current branch)",
			},
			"draft": map[string]any{
				"type":        "boolean",
				"description": "Create as draft pull request",
			},
		},
		"required": []any{"title", "body"},
	})
	return s
}

// GHPRCreate returns a tool handler that creates a GitHub pull request.
func GHPRCreate() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "title", "body"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		title := helpers.GetString(req.Arguments, "title")
		body := helpers.GetString(req.Arguments, "body")
		path := helpers.GetString(req.Arguments, "path")
		base := helpers.GetString(req.Arguments, "base")
		head := helpers.GetString(req.Arguments, "head")
		draft := helpers.GetBool(req.Arguments, "draft")

		args := []string{"pr", "create", "--title", title, "--body", body}
		if base != "" {
			args = append(args, "--base", base)
		}
		if head != "" {
			args = append(args, "--head", head)
		}
		if draft {
			args = append(args, "--draft")
		}

		output, err := git.GH(ctx, path, args...)
		if err != nil {
			return helpers.ErrorResult("gh_error", err.Error()), nil
		}
		return helpers.TextResult(output), nil
	}
}
