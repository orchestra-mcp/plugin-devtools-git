package devtoolsgit

import (
	"github.com/orchestra-mcp/plugin-devtools-git/internal"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

// Register adds all Git and GitHub tools to the builder.
func Register(builder *plugin.PluginBuilder) {
	tp := &internal.ToolsPlugin{}
	tp.RegisterTools(builder)
}
