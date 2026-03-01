package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-git/internal/git"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// GitLogSchema returns the JSON Schema for the git_log tool.
func GitLogSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Working directory path (defaults to current directory)",
			},
			"count": map[string]any{
				"type":        "number",
				"description": "Number of commits to show (default 10)",
			},
			"author": map[string]any{
				"type":        "string",
				"description": "Filter by author name or email",
			},
			"since": map[string]any{
				"type":        "string",
				"description": "Show commits after this date (e.g. 2024-01-01)",
			},
			"until": map[string]any{
				"type":        "string",
				"description": "Show commits before this date",
			},
			"file_path": map[string]any{
				"type":        "string",
				"description": "Show commits that modified this file",
			},
		},
	})
	return s
}

// logEntry represents a single parsed commit from git log.
type logEntry struct {
	Hash    string `json:"hash"`
	Author  string `json:"author"`
	Date    string `json:"date"`
	Subject string `json:"subject"`
}

// GitLog returns a tool handler that shows commit history.
func GitLog() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		path := helpers.GetString(req.Arguments, "path")
		count := helpers.GetInt(req.Arguments, "count")
		author := helpers.GetString(req.Arguments, "author")
		since := helpers.GetString(req.Arguments, "since")
		until := helpers.GetString(req.Arguments, "until")
		filePath := helpers.GetString(req.Arguments, "file_path")

		if count <= 0 {
			count = 10
		}

		// Use a delimiter to parse fields reliably.
		const sep = "<<<SEP>>>"
		format := strings.Join([]string{"%H", "%an", "%ai", "%s"}, sep)

		args := []string{"log", fmt.Sprintf("-n%d", count), fmt.Sprintf("--format=%s", format)}

		if author != "" {
			args = append(args, fmt.Sprintf("--author=%s", author))
		}
		if since != "" {
			args = append(args, fmt.Sprintf("--since=%s", since))
		}
		if until != "" {
			args = append(args, fmt.Sprintf("--until=%s", until))
		}

		// file_path must come after "--" to avoid ambiguity.
		if filePath != "" {
			args = append(args, "--", filePath)
		}

		output, err := git.Run(ctx, path, args...)
		if err != nil {
			return helpers.ErrorResult("git_error", err.Error()), nil
		}
		if output == "" {
			return helpers.TextResult("No commits found"), nil
		}

		var entries []logEntry
		for _, line := range strings.Split(output, "\n") {
			parts := strings.SplitN(line, sep, 4)
			if len(parts) == 4 {
				entries = append(entries, logEntry{
					Hash:    parts[0],
					Author:  parts[1],
					Date:    parts[2],
					Subject: parts[3],
				})
			}
		}

		data, err := json.MarshalIndent(entries, "", "  ")
		if err != nil {
			return helpers.ErrorResult("parse_error", err.Error()), nil
		}
		return helpers.TextResult(string(data)), nil
	}
}
