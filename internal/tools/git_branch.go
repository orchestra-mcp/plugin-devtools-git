package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-git/internal/git"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// GitBranchSchema returns the JSON Schema for the git_branch tool.
func GitBranchSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"action": map[string]any{
				"type":        "string",
				"description": "Branch action: list, create, delete, or switch",
				"enum":        []any{"list", "create", "delete", "switch"},
			},
			"path": map[string]any{
				"type":        "string",
				"description": "Working directory path (defaults to current directory)",
			},
			"name": map[string]any{
				"type":        "string",
				"description": "Branch name (required for create, delete, switch)",
			},
			"base": map[string]any{
				"type":        "string",
				"description": "Base branch or ref for creating a new branch",
			},
		},
		"required": []any{"action"},
	})
	return s
}

// GitBranch returns a tool handler for branch operations.
func GitBranch() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "action"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		action := helpers.GetString(req.Arguments, "action")
		path := helpers.GetString(req.Arguments, "path")
		name := helpers.GetString(req.Arguments, "name")
		base := helpers.GetString(req.Arguments, "base")

		if err := helpers.ValidateOneOf(action, "list", "create", "delete", "switch"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		switch action {
		case "list":
			output, err := git.Run(ctx, path, "branch", "-a", "--no-color")
			if err != nil {
				return helpers.ErrorResult("git_error", err.Error()), nil
			}
			return helpers.TextResult(output), nil

		case "create":
			if name == "" {
				return helpers.ErrorResult("validation_error", "name is required for create"), nil
			}
			args := []string{"branch", name}
			if base != "" {
				args = append(args, base)
			}
			output, err := git.Run(ctx, path, args...)
			if err != nil {
				return helpers.ErrorResult("git_error", err.Error()), nil
			}
			if output == "" {
				output = "Branch '" + name + "' created"
			}
			return helpers.TextResult(output), nil

		case "delete":
			if name == "" {
				return helpers.ErrorResult("validation_error", "name is required for delete"), nil
			}
			output, err := git.Run(ctx, path, "branch", "-d", name)
			if err != nil {
				return helpers.ErrorResult("git_error", err.Error()), nil
			}
			return helpers.TextResult(output), nil

		case "switch":
			if name == "" {
				return helpers.ErrorResult("validation_error", "name is required for switch"), nil
			}
			output, err := git.Run(ctx, path, "switch", name)
			if err != nil {
				return helpers.ErrorResult("git_error", err.Error()), nil
			}
			if output == "" {
				output = "Switched to branch '" + name + "'"
			}
			return helpers.TextResult(output), nil
		}

		return helpers.ErrorResult("validation_error", "unknown action: "+action), nil
	}
}
