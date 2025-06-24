package plugins

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"

	lua "github.com/yuin/gopher-lua"
)

// Oserver interface
type Plugin struct {
	ScriptName string
	Id         string
}

func (p *Plugin) Update(args []string, lua_state map[string]*lua.LState) []string {
	if handler, ok := lua_state[p.ScriptName]; ok {

		err := handler.DoString(fmt.Sprintf("result = HandleShortcode({%s})", strings.Join(args, ",")))
		if err != nil {
			log.Error().Msgf("Could not execute the plugin: %s", err)
			return args
		}
		value := handler.GetGlobal("result")
		if ret_type := value.Type().String(); ret_type != "string" {
			return args

		}
	}

	log.Error().Msgf("Could find a plugin: for: %s", p.ScriptName)
	return args
}
