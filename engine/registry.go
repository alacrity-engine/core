package engine

import "fmt"

// sceneRegistries contains component type registries
// from all the scenes of a game.
var sceneRegistries map[string]map[string]ComponentTypeEntry

// RegisterScene registers all the
// components of the scene scripts.
func RegisterScene(sceneID string, compRegistry map[string]ComponentTypeEntry) error {
	if _, ok := sceneRegistries[sceneID]; ok {
		return fmt.Errorf(
			"a scene with an ID=%s is already registered",
			sceneID)
	}

	sceneRegistries[sceneID] = compRegistry

	return nil
}
