package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-git/internal/git"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// GHPRReviewSchema returns the JSON Schema for the gh_pr_review tool.
func GHPRReviewSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"pr_number": map[string]any{
				"type":        "number",
				"description": "Pull request number",
			},
			"action": map[string]any{
				"type":        "string",
				"description": "Review action: approve, comment, or request-changes",
				"enum":        []any{"approve", "comment", "request-changes"},
			},
			"path": map[string]any{
				"type":        "string",
				"description": "Working directory path (defaults to current directory)",
			},
			"body": map[string]any{
				"type":        "string",
				"description": "Review comment body",
			},
		},
		"required": []any{"pr_number", "action"},
	})
	return s
}

// GHPRReview returns a tool handler that submits a PR review.
func GHPRReview() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "action"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		prNumber := helpers.GetInt(req.Arguments, "pr_number")
		if prNumber == 0 {
			return helpers.ErrorResult("validation_error", "pr_number is required"), nil
		}

		action := helpers.GetString(req.Arguments, "action")
		path := helpers.GetString(req.Arguments, "path")
		body := helpers.GetString(req.Arguments, "body")

		if err := helpers.ValidateOneOf(action, "approve", "comment", "request-changes"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		args := []string{"pr", "review", fmt.Sprintf("%d", prNumber)}

		switch action {
		case "approve":
			args = append(args, "--approve")
		case "comment":
			args = append(args, "--comment")
		case "request-changes":
			args = append(args, "--request-changes")
		}

		if body != "" {
			args = append(args, "--body", body)
		}

		output, err := git.GH(ctx, path, args...)
		if err != nil {
			return helpers.ErrorResult("gh_error", err.Error()), nil
		}
		if output == "" {
			output = fmt.Sprintf("Review submitted for PR #%d: %s", prNumber, action)
		}
		return helpers.TextResult(output), nil
	}
}
