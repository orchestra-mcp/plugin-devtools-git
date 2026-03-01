package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-git/internal/git"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// GHReleaseCreateSchema returns the JSON Schema for the gh_release_create tool.
func GHReleaseCreateSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"tag": map[string]any{
				"type":        "string",
				"description": "Tag name for the release (e.g. v1.0.0)",
			},
			"path": map[string]any{
				"type":        "string",
				"description": "Working directory path (defaults to current directory)",
			},
			"title": map[string]any{
				"type":        "string",
				"description": "Release title (defaults to tag name)",
			},
			"body": map[string]any{
				"type":        "string",
				"description": "Release notes body",
			},
			"draft": map[string]any{
				"type":        "boolean",
				"description": "Create as draft release",
			},
			"prerelease": map[string]any{
				"type":        "boolean",
				"description": "Mark as prerelease",
			},
		},
		"required": []any{"tag"},
	})
	return s
}

// GHReleaseCreate returns a tool handler that creates a GitHub release.
func GHReleaseCreate() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "tag"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		tag := helpers.GetString(req.Arguments, "tag")
		path := helpers.GetString(req.Arguments, "path")
		title := helpers.GetString(req.Arguments, "title")
		body := helpers.GetString(req.Arguments, "body")
		draft := helpers.GetBool(req.Arguments, "draft")
		prerelease := helpers.GetBool(req.Arguments, "prerelease")

		args := []string{"release", "create", tag}

		if title != "" {
			args = append(args, "--title", title)
		}
		if body != "" {
			args = append(args, "--notes", body)
		} else {
			args = append(args, "--generate-notes")
		}
		if draft {
			args = append(args, "--draft")
		}
		if prerelease {
			args = append(args, "--prerelease")
		}

		output, err := git.GH(ctx, path, args...)
		if err != nil {
			return helpers.ErrorResult("gh_error", err.Error()), nil
		}
		return helpers.TextResult(output), nil
	}
}
