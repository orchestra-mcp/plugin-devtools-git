package internal

import (
	"github.com/orchestra-mcp/plugin-devtools-git/internal/tools"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

// ToolsPlugin registers all Git and GitHub tools.
type ToolsPlugin struct{}

// RegisterTools registers all 20 tools with the plugin builder.
func (tp *ToolsPlugin) RegisterTools(builder *plugin.PluginBuilder) {
	// Git operations (10 tools)
	builder.RegisterTool("git_status",
		"Show working tree status (staged, unstaged, untracked files)",
		tools.GitStatusSchema(), tools.GitStatus())

	builder.RegisterTool("git_diff",
		"Show diff between refs or working tree changes (staged or unstaged)",
		tools.GitDiffSchema(), tools.GitDiff())

	builder.RegisterTool("git_log",
		"Show commit history with optional filters (author, date, file)",
		tools.GitLogSchema(), tools.GitLog())

	builder.RegisterTool("git_commit",
		"Stage files and create a commit with a message",
		tools.GitCommitSchema(), tools.GitCommit())

	builder.RegisterTool("git_branch",
		"Create, list, delete, or switch branches",
		tools.GitBranchSchema(), tools.GitBranch())

	builder.RegisterTool("git_merge",
		"Merge a branch into the current branch",
		tools.GitMergeSchema(), tools.GitMerge())

	builder.RegisterTool("git_stash",
		"Stash or unstash working directory changes",
		tools.GitStashSchema(), tools.GitStash())

	builder.RegisterTool("git_blame",
		"Show line-by-line blame annotation for a file",
		tools.GitBlameSchema(), tools.GitBlame())

	builder.RegisterTool("git_checkout",
		"Checkout a branch, tag, commit, or specific files",
		tools.GitCheckoutSchema(), tools.GitCheckout())

	builder.RegisterTool("git_tag",
		"Create, list, or delete tags",
		tools.GitTagSchema(), tools.GitTag())

	// GitHub operations (10 tools via gh CLI)
	builder.RegisterTool("gh_pr_create",
		"Create a GitHub pull request",
		tools.GHPRCreateSchema(), tools.GHPRCreate())

	builder.RegisterTool("gh_pr_list",
		"List GitHub pull requests with optional filters",
		tools.GHPRListSchema(), tools.GHPRList())

	builder.RegisterTool("gh_pr_review",
		"Submit a review on a GitHub pull request (approve, comment, request changes)",
		tools.GHPRReviewSchema(), tools.GHPRReview())

	builder.RegisterTool("gh_pr_merge",
		"Merge a GitHub pull request (merge, squash, or rebase)",
		tools.GHPRMergeSchema(), tools.GHPRMerge())

	builder.RegisterTool("gh_issue_create",
		"Create a GitHub issue with optional labels and assignees",
		tools.GHIssueCreateSchema(), tools.GHIssueCreate())

	builder.RegisterTool("gh_issue_list",
		"List GitHub issues with optional filters",
		tools.GHIssueListSchema(), tools.GHIssueList())

	builder.RegisterTool("gh_issue_comment",
		"Add a comment to a GitHub issue",
		tools.GHIssueCommentSchema(), tools.GHIssueComment())

	builder.RegisterTool("gh_actions_status",
		"Show GitHub Actions CI/CD workflow run status",
		tools.GHActionsStatusSchema(), tools.GHActionsStatus())

	builder.RegisterTool("gh_release_create",
		"Create a GitHub release with optional notes",
		tools.GHReleaseCreateSchema(), tools.GHReleaseCreate())

	builder.RegisterTool("gh_repo_info",
		"Show GitHub repository metadata (stars, forks, languages, etc.)",
		tools.GHRepoInfoSchema(), tools.GHRepoInfo())
}
