package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-git/internal/git"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// GitStashSchema returns the JSON Schema for the git_stash tool.
func GitStashSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"action": map[string]any{
				"type":        "string",
				"description": "Stash action: push, pop, list, or drop",
				"enum":        []any{"push", "pop", "list", "drop"},
			},
			"path": map[string]any{
				"type":        "string",
				"description": "Working directory path (defaults to current directory)",
			},
			"message": map[string]any{
				"type":        "string",
				"description": "Stash message (for push action)",
			},
			"index": map[string]any{
				"type":        "number",
				"description": "Stash index to pop or drop (default 0)",
			},
		},
		"required": []any{"action"},
	})
	return s
}

// GitStash returns a tool handler for stash operations.
func GitStash() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "action"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		action := helpers.GetString(req.Arguments, "action")
		path := helpers.GetString(req.Arguments, "path")
		message := helpers.GetString(req.Arguments, "message")
		index := helpers.GetInt(req.Arguments, "index")

		if err := helpers.ValidateOneOf(action, "push", "pop", "list", "drop"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		var args []string

		switch action {
		case "push":
			args = []string{"stash", "push"}
			if message != "" {
				args = append(args, "-m", message)
			}
		case "pop":
			args = []string{"stash", "pop", fmt.Sprintf("stash@{%d}", index)}
		case "list":
			args = []string{"stash", "list"}
		case "drop":
			args = []string{"stash", "drop", fmt.Sprintf("stash@{%d}", index)}
		}

		output, err := git.Run(ctx, path, args...)
		if err != nil {
			return helpers.ErrorResult("git_error", err.Error()), nil
		}
		if output == "" {
			output = "Stash operation completed"
		}
		return helpers.TextResult(output), nil
	}
}
