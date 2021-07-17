package ecs

import (
	"fmt"
	"sort"

	"github.com/alacrity-engine/core/tasking"
)

// Scene is a collection of game objects
// to be updated and drawn.
type (
	Scene struct {
		name              string
		gmobs             []*GameObject
		addBuffer         []add
		destructionBuffer []string
		changeZBuffer     []changeZ
		systems           map[string]System
		taskMgr           *tasking.TaskManager
		layout            *DrawLayout
	}

	changeZ struct {
		gmobName string
		targetZ  float64
	}

	add struct {
		gmob *GameObject
		zUpd float64
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

// DrawLayout returns the draw layout og the scene.
func (scene *Scene) DrawLayout() *DrawLayout {
	return scene.layout
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

	/*err = scene.placeGameObjects()

	if err != nil {
		return err
	}*/

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

// insertGameObject inserts the game object
// into the sorted Z-buffer using binary search.
func (scene *Scene) insertGameObject(gmob *GameObject, zUpd float64) {
	// Use binary search to insert the game
	// object into the Z-sorted update buffer.
	length := len(scene.gmobs)
	ind := sort.Search(length, func(i int) bool {
		return scene.gmobs[i].zUpdate >= zUpd
	})

	if ind == 0 {
		scene.gmobs = append([]*GameObject{gmob},
			scene.gmobs...)
	} else if ind < length {
		scene.gmobs = append(scene.gmobs[:ind+1],
			scene.gmobs[ind:]...)
		scene.gmobs[ind] = gmob
	} else {
		scene.gmobs = append(scene.gmobs, gmob)
	}
}

// placeGameObjects changes Z update coordinate
// of all the requested game objects.
func (scene *Scene) placeGameObjects() error {
	for _, changeZ := range scene.changeZBuffer {
		pos, gmob := scene.FindGameObject(changeZ.gmobName)

		if pos < 0 {
			return RaiseErrorNoGameObjectOnScene(
				scene, changeZ.gmobName)
		}

		// Remove the game object from
		// its previous position.
		scene.gmobs = append(scene.gmobs[:pos],
			scene.gmobs[pos+1:]...)
		// Insert the game object back
		// into the Z-buffer.
		scene.insertGameObject(gmob, changeZ.targetZ)
	}

	return nil
}

// addBufferedGameObjects adds all the buffered
// game objects to the scene.
func (scene *Scene) addBufferedGameObjects() error {
	// Add all the game objects from the buffer
	// and start them all.
	for _, gmobAdd := range scene.addBuffer {
		err := scene.AddGameObject(gmobAdd.gmob,
			gmobAdd.zUpd)

		if err != nil {
			return err
		}

		err = gmobAdd.gmob.Start()

		if err != nil {
			return err
		}
	}

	scene.addBuffer = []add{}

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
func (scene *Scene) findGameObjectInAdded(name string) (int, add) {
	ind := -1
	var gameObjectAdd add

	for i, gmobAdd := range scene.addBuffer {
		if name == gmobAdd.gmob.Name() {
			ind = i
			gameObjectAdd = gmobAdd

			break
		}
	}

	return ind, gameObjectAdd
}

// HasGameObject returns true if the scene has a game
// object with the specified name, and false otherwise.
func (scene *Scene) HasGameObject(name string) bool {
	_, gmob := scene.FindGameObject(name)

	return gmob != nil
}

// AddGameObject adds a new game object on the scene.
//
// Should be used before the scene is started.
func (scene *Scene) AddGameObject(gmob *GameObject, zUpd float64) error {
	if scene.HasGameObject(gmob.name) {
		return fmt.Errorf("scene '%s' already has game object '%s'",
			scene.name, gmob.name)
	}

	scene.insertGameObject(gmob, zUpd)
	gmob.SetScene(scene)

	return nil
}

// AddGameObjectInRuntime must be called when the game
// object is created in the game loop and should be added
// to the scene.
func (scene *Scene) AddGameObjectInRuntime(gmob *GameObject, zUpd float64) error {
	if scene.HasGameObject(gmob.name) {
		return fmt.Errorf("scene '%s' already has game object '%s'",
			scene.name, gmob.name)
	}

	i, _ := scene.findGameObjectInAdded(gmob.name)

	if i > 0 {
		return fmt.Errorf("game object '%s' is already set to be added",
			gmob.name)
	}

	scene.addBuffer = append(scene.addBuffer, add{gmob, zUpd})

	return nil
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
		addBuffer:         []add{},
		changeZBuffer:     []changeZ{},
		gmobs:             []*GameObject{},
		destructionBuffer: []string{},
		systems:           map[string]System{},
		taskMgr:           tasking.NewTaskManager(),
		layout:            NewLayout(),
	}
}
