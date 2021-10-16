package engine

import (
	"plugin"

	"github.com/alacrity-engine/core/tasking"
)

var (
	// noDestroyOnSceneSwitch contains all the objects
	// that should not be destroyed on switching to a new scene.
	noDestroyOnSceneSwitch map[string]*GameObject
	// scenes stores all the loaded scenes.
	scenes map[string]*Scene
)

// StartLoadingScene starts an asynchronous
// process of loading a new scene.
func StartLoadingScene(scenePath string) (*tasking.AsynchronousProcess, error) {
	lib, err := plugin.Open(scenePath)

	if err != nil {
		return nil, err
	}

	createSceneSymbol, err := lib.Lookup("CreateScene")

	if err != nil {
		return nil, err
	}

	createScene, ok := createSceneSymbol.(func() *tasking.AsynchronousProcess)

	if !ok {
		return nil, NewErrorIncorrectSceneLoader(scenePath)
	}

	return createScene(), nil
}

// AddScene adds a new scene to the buffer.
func AddScene(scene *Scene) error {
	if _, ok := scenes[scene.name]; ok {
		return NewErrorSceneAlreadyExists(scene)
	}

	scenes[scene.name] = scene

	return nil
}

// RemoveScene removes the existing scene
// from the buffer.
func RemoveScene(sceneName string) error {
	if _, ok := scenes[sceneName]; !ok {
		return NewErrorSceneDoesntExist(sceneName)
	}

	delete(scenes, sceneName)

	return nil
}
