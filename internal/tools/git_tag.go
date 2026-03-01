package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-git/internal/git"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// GitTagSchema returns the JSON Schema for the git_tag tool.
func GitTagSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"action": map[string]any{
				"type":        "string",
				"description": "Tag action: list, create, or delete",
				"enum":        []any{"list", "create", "delete"},
			},
			"path": map[string]any{
				"type":        "string",
				"description": "Working directory path (defaults to current directory)",
			},
			"name": map[string]any{
				"type":        "string",
				"description": "Tag name (required for create and delete)",
			},
			"message": map[string]any{
				"type":        "string",
				"description": "Tag message (creates annotated tag)",
			},
			"ref": map[string]any{
				"type":        "string",
				"description": "Ref to tag (defaults to HEAD)",
			},
		},
		"required": []any{"action"},
	})
	return s
}

// GitTag returns a tool handler for tag operations.
func GitTag() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "action"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		action := helpers.GetString(req.Arguments, "action")
		path := helpers.GetString(req.Arguments, "path")
		name := helpers.GetString(req.Arguments, "name")
		message := helpers.GetString(req.Arguments, "message")
		ref := helpers.GetString(req.Arguments, "ref")

		if err := helpers.ValidateOneOf(action, "list", "create", "delete"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		switch action {
		case "list":
			output, err := git.Run(ctx, path, "tag", "-l", "--sort=-creatordate")
			if err != nil {
				return helpers.ErrorResult("git_error", err.Error()), nil
			}
			if output == "" {
				output = "No tags found"
			}
			return helpers.TextResult(output), nil

		case "create":
			if name == "" {
				return helpers.ErrorResult("validation_error", "name is required for create"), nil
			}
			args := []string{"tag"}
			if message != "" {
				args = append(args, "-a", name, "-m", message)
			} else {
				args = append(args, name)
			}
			if ref != "" {
				args = append(args, ref)
			}
			output, err := git.Run(ctx, path, args...)
			if err != nil {
				return helpers.ErrorResult("git_error", err.Error()), nil
			}
			if output == "" {
				output = "Tag '" + name + "' created"
			}
			return helpers.TextResult(output), nil

		case "delete":
			if name == "" {
				return helpers.ErrorResult("validation_error", "name is required for delete"), nil
			}
			output, err := git.Run(ctx, path, "tag", "-d", name)
			if err != nil {
				return helpers.ErrorResult("git_error", err.Error()), nil
			}
			return helpers.TextResult(output), nil
		}

		return helpers.ErrorResult("validation_error", "unknown action: "+action), nil
	}
}
