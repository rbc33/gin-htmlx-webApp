package plugins

import lua "github.com/yuin/gopher-lua"

// Subject interface
type Hook interface {
	Register(plugin Plugin)
	Deregister(plugin Plugin)
	NotifyAll(lua_states map[string]*lua.LState) []string
}
