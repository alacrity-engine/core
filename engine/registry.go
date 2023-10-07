package engine

import "fmt"

// sceneRegistries contains component type registries
// from all the scenes of a game.
var compTypeRegistry map[string]ComponentTypeEntry

// RegisterScene registers all the
// components of the scene scripts.
func SetRegistry(registry map[string]ComponentTypeEntry) error {
	if compTypeRegistry != nil {
		return fmt.Errorf(
			"the component type registry is already set")
	}

	compTypeRegistry = registry

	return nil
}
