package plugins

import (
	"strconv"

	"github.com/rbc33/gocms/common"
	lua "github.com/yuin/gopher-lua"
)

type PostHook struct {
	plugins    []Plugin
	latestPost common.Post
}

// runs at start up
func (p PostHook) Register(plugin Plugin) {
	p.plugins = append(p.plugins, plugin)
}

func (p PostHook) Deregister(plugin Plugin) {
	// look each plugin and remove the
	// one we don't want
}

func (p PostHook) NotifyAll(lua_state map[string]*lua.LState) []string {
	// iterate over each plugin

	args := []string{
		strconv.Itoa(p.latestPost.Id),
		p.latestPost.Title,
		p.latestPost.Excerpt,
		p.latestPost.Content,
	}
	for _, plugin := range p.plugins {
		args = plugin.Update(args, lua_state)
	}
	return args
}

func (p PostHook) UpdatePost(title string, content string, excerpt string, lua_states map[string]*lua.LState) common.Post {
	p.latestPost = common.Post{
		Title:   title,
		Content: content,
		Excerpt: excerpt,
		Id:      -1,
	}
	args := p.NotifyAll(lua_states)
	return common.Post{
		Title:   args[1],
		Content: args[2],
		Excerpt: args[3],
		Id:      -1,
	}
}
