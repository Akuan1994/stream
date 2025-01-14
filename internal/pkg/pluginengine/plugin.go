package pluginengine

import (
	"fmt"
	"plugin"

	"github.com/merico-dev/stream/internal/pkg/configloader"
)

// DevStreamPlugin is a struct, on which install/reinstall/uninstall interfaces are defined.
type DevStreamPlugin interface {
	// Install will return (true, nil) if there is no error occurred. Otherwise (false, error) will be returned.
	Install(*map[string]interface{}) (bool, error)
	Reinstall(*map[string]interface{}) (bool, error)
	Uninstall(*map[string]interface{}) (bool, error)
}

// Install loads the plugin and calls the Install method of that plugin.
func Install(tool *configloader.Tool) (bool, error) {
	p, err := loadPlugin(tool)
	if err != nil {
		return false, err
	}
	return p.Install(&tool.Options)
}

// Reinstall loads the plugin and calls the Reinstall method of that plugin.
func Reinstall(tool *configloader.Tool) (bool, error) {
	p, err := loadPlugin(tool)
	if err != nil {
		return false, err
	}
	return p.Reinstall(&tool.Options)
}

// Uninstall loads the plugin and calls the Uninstall method of that plugin.
func Uninstall(tool *configloader.Tool) (bool, error) {
	p, err := loadPlugin(tool)
	if err != nil {
		return false, err
	}
	return p.Uninstall(&tool.Options)
}

func loadPlugin(tool *configloader.Tool) (DevStreamPlugin, error) {
	mod := fmt.Sprintf("plugins/%s_%s.so", tool.Name, tool.Version)
	plug, err := plugin.Open(mod)
	if err != nil {
		return nil, err
	}

	var devStreamPlugin DevStreamPlugin
	symDevStreamPlugin, err := plug.Lookup("DevStreamPlugin")
	if err != nil {
		return nil, err
	}

	devStreamPlugin, ok := symDevStreamPlugin.(DevStreamPlugin)
	if !ok {
		return nil, err
	}

	return devStreamPlugin, nil
}
