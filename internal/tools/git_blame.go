package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-git/internal/git"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// GitBlameSchema returns the JSON Schema for the git_blame tool.
func GitBlameSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"file": map[string]any{
				"type":        "string",
				"description": "File to blame",
			},
			"path": map[string]any{
				"type":        "string",
				"description": "Working directory path (defaults to current directory)",
			},
			"start_line": map[string]any{
				"type":        "number",
				"description": "Start line number for blame range",
			},
			"end_line": map[string]any{
				"type":        "number",
				"description": "End line number for blame range",
			},
		},
		"required": []any{"file"},
	})
	return s
}

// GitBlame returns a tool handler that shows line-by-line blame.
func GitBlame() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "file"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		file := helpers.GetString(req.Arguments, "file")
		path := helpers.GetString(req.Arguments, "path")
		startLine := helpers.GetInt(req.Arguments, "start_line")
		endLine := helpers.GetInt(req.Arguments, "end_line")

		args := []string{"blame", "--no-color"}

		if startLine > 0 && endLine > 0 {
			args = append(args, fmt.Sprintf("-L%d,%d", startLine, endLine))
		} else if startLine > 0 {
			args = append(args, fmt.Sprintf("-L%d,", startLine))
		}

		args = append(args, file)

		output, err := git.Run(ctx, path, args...)
		if err != nil {
			return helpers.ErrorResult("git_error", err.Error()), nil
		}
		return helpers.TextResult(output), nil
	}
}
