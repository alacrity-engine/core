package ecs

import (
	"fmt"

	"github.com/alacrity-engine/core/tasking"
)

// Scene is a collection of game objects
// to be updated and drawn.
type (
	Scene struct {
		name              string
		gmobs             []*GameObject
		additionBuffer    []*GameObject
		destructionBuffer []string
		pasteBuffer       []*paste
		systems           map[string]System
		taskMgr           *tasking.TaskManager
	}

	paste struct {
		movedGmob  string
		targetGmob string
		action     string
	}
)

// Name returns the name of the scene.
func (scene *Scene) Name() string {
	return scene.name
}

// TaskManager returns the task manager of the scene.
func (scene *Scene) TaskManager() *tasking.TaskManager {
	return scene.taskMgr
}

// FindSystem returns the system with the specified name.
func (scene *Scene) FindSystem(name string) (System, error) {
	system, exists := scene.systems[name]

	if !exists {
		return nil, fmt.Errorf("scene '%s' has no system '%s'",
			scene.name, name)
	}

	return system, nil
}

// AddSystem adds the system to the scene and assigns the name to it.
func (scene *Scene) AddSystem(name string, system System) error {
	_, err := scene.FindSystem(name)

	if err == nil {
		return fmt.Errorf("scene '%s' already has system '%s'",
			scene.name, name)
	}

	scene.systems[name] = system

	return nil
}

// RemoveSystem removes the system from the scene by its name.
func (scene *Scene) RemoveSystem(name string) error {
	_, err := scene.FindSystem(name)

	if err != nil {
		return fmt.Errorf("scene '%s' doesn't have system '%s'",
			scene.name, name)
	}

	delete(scene.systems, name)

	return nil
}

// Start starts components of all
// the game objects on the scene.
func (scene *Scene) Start() error {
	for _, gmob := range scene.gmobs {
		err := gmob.Start()

		if err != nil {
			return err
		}
	}

	return nil
}

// Update calls update method on all
// game objects of the scene.
func (scene *Scene) Update() error {
	err := scene.removeDestroyedGameObjects()

	if err != nil {
		return err
	}

	err = scene.placeGameObjects()

	if err != nil {
		return err
	}

	err = scene.addBufferedGameObjects()

	if err != nil {
		return err
	}

	// Update all the game objects.
	for _, gmob := range scene.gmobs {
		err = gmob.Update()

		if err != nil {
			return err
		}
	}

	// Perform the next iteration of the tasks.
	err = scene.taskMgr.Update()

	if err != nil {
		return err
	}

	return nil
}

// placeGameObjects changes positions of the game
// objects requested during runtime.
func (scene *Scene) placeGameObjects() error {
	for _, paste := range scene.pasteBuffer {
		switch paste.action {
		case "before":
			err := scene.SetGameObjectPriorityBefore(
				paste.movedGmob, paste.targetGmob)

			if err != nil {
				switch err.(type) {
				case *ErrorNoGameObjectOnScene:
					continue

				default:
					return err
				}
			}

		case "after":
			err := scene.SetGameObjectPriorityAfter(
				paste.movedGmob, paste.targetGmob)

			if err != nil {
				switch err.(type) {
				case *ErrorNoGameObjectOnScene:
					continue

				default:
					return err
				}
			}

		default:
			return fmt.Errorf("unknown action '%s'",
				paste.action)
		}
	}

	scene.pasteBuffer = []*paste{}

	return nil
}

// addBufferedGameObjects adds all the buffered
// game objects to the scene.
func (scene *Scene) addBufferedGameObjects() error {
	scene.gmobs = append(scene.gmobs, scene.additionBuffer...)

	for _, gmob := range scene.additionBuffer {
		gmob.SetScene(scene)
		err := gmob.Start()

		if err != nil {
			return err
		}
	}

	scene.additionBuffer = []*GameObject{}

	return nil
}

// findDestroyedGameObject the game object with
// the specified name marked as destroyed.
func (scene *Scene) findDestroyedGameObject(name string) (int, *GameObject) {
	ind := -1
	var gameObject *GameObject

	for i, gmob := range scene.gmobs {
		if name == gmob.Name() && gmob.destroyed {
			ind = i
			gameObject = gmob

			break
		}
	}

	return ind, gameObject
}

// removeDestroyedGameObject removes the game object
// with the specified name marked as destroyed.
func (scene *Scene) removeDestroyedGameObject(name string) error {
	i, gmob := scene.findDestroyedGameObject(name)

	if gmob == nil {
		return fmt.Errorf("scene '%s' doesn't have destroyed game object '%s'",
			scene.name, name)
	}

	scene.gmobs = append(scene.gmobs[:i], scene.gmobs[i+1:]...)

	return nil
}

// removeDestroyedGameObjects removes all the game
// objects set for destruction from the scene.
//
// ATTENTION: this method must not be called from
// any ecs.Component. Call it after scene.Update().
func (scene *Scene) removeDestroyedGameObjects() error {
	for _, gmob := range scene.destructionBuffer {
		err := scene.removeDestroyedGameObject(gmob)

		if err != nil {
			return err
		}
	}

	scene.destructionBuffer = []string{}

	return nil
}

// PlaceGameObjectBefore sets the game object to be
// placed before the other game object in the next frame.
func (scene *Scene) PlaceGameObjectBefore(name, beforeName string) error {
	_, gmob := scene.FindGameObject(name)

	if gmob == nil {
		return RaiseErrorNoGameObjectOnScene(scene, name)
	}

	_, beforeGmob := scene.FindGameObject(beforeName)

	if beforeGmob == nil {
		return RaiseErrorNoGameObjectOnScene(scene, beforeName)
	}

	p := &paste{
		movedGmob:  name,
		targetGmob: beforeName,
		action:     "before",
	}

	scene.pasteBuffer = append(scene.pasteBuffer, p)

	return nil
}

// PlaceGameObjectAfter sets the game object to be
// placed after the other game object in the next frame.
func (scene *Scene) PlaceGameObjectAfter(name, afterName string) error {
	_, gmob := scene.FindGameObject(name)

	if gmob == nil {
		return RaiseErrorNoGameObjectOnScene(scene, name)
	}

	_, beforeGmob := scene.FindGameObject(afterName)

	if beforeGmob == nil {
		return RaiseErrorNoGameObjectOnScene(scene, afterName)
	}

	p := &paste{
		movedGmob:  name,
		targetGmob: afterName,
		action:     "after",
	}

	scene.pasteBuffer = append(scene.pasteBuffer, p)

	return nil
}

// SetGameObjectPriorityBefore sets the game object
// priority to be updated before the specified game object.
func (scene *Scene) SetGameObjectPriorityBefore(name, beforeName string) error {
	i, gmob := scene.FindGameObject(name)

	if gmob == nil {
		return RaiseErrorNoGameObjectOnScene(scene, name)
	}

	j, beforeGmob := scene.FindGameObject(beforeName)

	if beforeGmob == nil {
		return RaiseErrorNoGameObjectOnScene(scene, beforeName)
	}

	j--

	if j < 0 {
		j = 0
	}

	scene.gmobs = append(scene.gmobs[:i], scene.gmobs[i+1:]...)
	buffer := make([]*GameObject, j+1)
	copy(buffer, scene.gmobs[:j])
	buffer[j] = gmob
	scene.gmobs = append(buffer, scene.gmobs[j:]...)

	return nil
}

// SetGameObjectPriorityAfter sets the game object
// priority to be updated after the specified game object.
func (scene *Scene) SetGameObjectPriorityAfter(name, afterName string) error {
	i, gmob := scene.FindGameObject(name)

	if gmob == nil {
		return RaiseErrorNoGameObjectOnScene(scene, name)
	}

	j, beforeGmob := scene.FindGameObject(afterName)

	if beforeGmob == nil {
		return RaiseErrorNoGameObjectOnScene(scene, afterName)
	}

	j++
	lastInd := len(scene.gmobs) - 1

	if j > lastInd {
		j = lastInd
	}

	scene.gmobs = append(scene.gmobs[:i], scene.gmobs[i+1:]...)
	buffer := make([]*GameObject, j+1)
	copy(buffer, scene.gmobs[:j])
	buffer[j] = gmob
	scene.gmobs = append(buffer, scene.gmobs[j:]...)

	return nil
}

// FindGameObject finds the game object on the scene.
func (scene *Scene) FindGameObject(name string) (int, *GameObject) {
	ind := -1
	var gameObject *GameObject

	for i, gmob := range scene.gmobs {
		if name == gmob.Name() && !gmob.destroyed {
			ind = i
			gameObject = gmob

			break
		}
	}

	return ind, gameObject
}

// findGameObjectInDestroyed searches for a game object
// with the specified name in the buffer where game objects
// set to be destroyed reside.
func (scene *Scene) hasGameObjectInDestroyed(name string) bool {
	found := false

	for _, gmob := range scene.destructionBuffer {
		if name == gmob {
			found = true

			break
		}
	}

	return found
}

// findGameObjectInAdded searches for a game object
// with the specified name in the buffer where game objects
// set to be added reside.
func (scene *Scene) findGameObjectInAdded(name string) (int, *GameObject) {
	ind := -1
	var gameObject *GameObject

	for i, gmob := range scene.additionBuffer {
		if name == gmob.Name() {
			ind = i
			gameObject = gmob

			break
		}
	}

	return ind, gameObject
}

// HasGameObject returns true if the scene has a game
// object with the specified name, and false otherwise.
func (scene *Scene) HasGameObject(name string) bool {
	_, gmob := scene.FindGameObject(name)

	return gmob != nil
}

// AddGameObject adds a new game object on the scene.
func (scene *Scene) AddGameObject(gmob *GameObject, priority int) error {
	if scene.HasGameObject(gmob.name) {
		return fmt.Errorf("scene '%s' already has game object '%s'",
			scene.name, gmob.name)
	}

	length := len(scene.gmobs)

	if length <= 0 || priority >= length {
		scene.gmobs = append(scene.gmobs, gmob)
	} else if priority < 0 {
		scene.gmobs = append([]*GameObject{gmob},
			scene.gmobs...)
	} else {
		scene.gmobs = append(scene.gmobs[:priority+1],
			scene.gmobs[priority:]...)
		scene.gmobs[priority] = gmob
	}

	gmob.SetScene(scene)

	return nil
}

// AddGameObjectToBuffer must be called when the game
// object is created in the game loop and should be added
// to the scene.
func (scene *Scene) AddGameObjectToBuffer(gmob *GameObject) error {
	if scene.HasGameObject(gmob.name) {
		return fmt.Errorf("scene '%s' already has game object '%s'",
			scene.name, gmob.name)
	}

	i, _ := scene.findGameObjectInAdded(gmob.name)

	if i > 0 {
		return fmt.Errorf("game object '%s' is already set to be added",
			gmob.name)
	}

	scene.additionBuffer = append(scene.additionBuffer, gmob)

	return nil
}

// PasteGameObjectBefore pastes the game object
// before the game object with the specified name
// in the game object list of the scene.
func (scene *Scene) PasteGameObjectBefore(gmobName string, gmob *GameObject) error {
	i, gmobReq := scene.FindGameObject(gmobName)

	if gmobReq == nil {
		return fmt.Errorf("scene '%s' doesn't have game object '%s'",
			scene.name, gmobName)
	}

	if scene.HasGameObject(gmob.name) {
		return fmt.Errorf("scene '%s' already has game object '%s'",
			scene.name, gmob.name)
	}

	i--

	return scene.AddGameObject(gmob, i)
}

// PasteGameObjectAfter pastes the game object
// after the game object with the specified name
// in the game object list of the scene.
func (scene *Scene) PasteGameObjectAfter(gmobName string, gmob *GameObject) error {
	i, gmobReq := scene.FindGameObject(gmobName)

	if gmobReq == nil {
		return fmt.Errorf("scene '%s' doesn't have game object '%s'",
			scene.name, gmobName)
	}

	if scene.HasGameObject(gmob.name) {
		return fmt.Errorf("scene '%s' already has game object '%s'",
			scene.name, gmob.name)
	}

	i++

	return scene.AddGameObject(gmob, i)
}

// RemoveGameObject removes the game object from the scene.
func (scene *Scene) RemoveGameObject(name string) error {
	i, gmob := scene.FindGameObject(name)

	if gmob == nil {
		return fmt.Errorf("scene '%s' doesn't have game object '%s'",
			scene.name, name)
	}

	scene.gmobs = append(scene.gmobs[:i], scene.gmobs[i+1:]...)

	return nil
}

// DestroyGameObject destroys the game object, i.e.
// deactivates all its components and stops drawing it.
func (scene *Scene) DestroyGameObject(name string) error {
	_, gmob := scene.FindGameObject(name)

	if gmob == nil {
		return fmt.Errorf("scene '%s' doesn't have game object '%s'",
			scene.name, name)
	}

	// Deactivate all the game object components.
	for _, component := range gmob.components {
		err := component.Destroy()

		if err != nil {
			return err
		}

		component.SetActive(false)
	}

	gmob.Transform().SetParent(nil)
	gmob.SetScene(nil)
	gmob.SetDraw(false)
	gmob.destroyed = true

	scene.destructionBuffer = append(scene.
		destructionBuffer, name)

	return nil
}

// NewScene creates a new scene to
// place game objects onto.
func NewScene(name string) *Scene {
	return &Scene{
		name:              name,
		additionBuffer:    []*GameObject{},
		gmobs:             []*GameObject{},
		destructionBuffer: []string{},
		systems:           map[string]System{},
		taskMgr:           tasking.NewTaskManager(),
	}
}
