package engine

import (
	"fmt"
	"path/filepath"

	"github.com/alacrity-engine/core/render"
	"github.com/alacrity-engine/core/resources"
	"github.com/alacrity-engine/core/system/collections"
	"github.com/alacrity-engine/core/tasking"
)

// TODO: mitigate memory allocations by
// replacing slices with pooled linked
// lists for addBuffer and changeZBuffer.
// Maybe it's also worth trying to replace
// standard collections (maps, slices)
// with pooled custom counetrparts.

// TODO: use int instead of float64 for zUpd.

// Scene is a collection of game objects
// to be updated and drawn.
type (
	Scene struct {
		name              string
		gmobs             collections.SortedDictionary[float64, map[*GameObject]struct{}]
		gmobNameIndex     map[string]*GameObject
		addBuffer         map[*GameObject]float64
		destructionBuffer map[string]struct{}
		changeZBuffer     []changeZ
		systems           map[string]System
		taskMgr           *tasking.TaskManager
		layout            *render.Layout
		resourceLoaders   map[string]*resources.ResourceLoader
	}

	changeZ struct {
		gmobName string
		targetZ  float64
	}

	add struct {
		gmob *GameObject
		zUpd float64
	}

	record struct {
		gmob *GameObject
		pos  int
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
func (scene *Scene) DrawLayout() *render.Layout {
	return scene.layout
}

// openResourceFile opens a new resource file for the scene.
func (scene *Scene) openResourceFile(fname string) (*resources.ResourceLoader, error) {
	loaderID, err := filepath.Abs(fname)

	if err != nil {
		return nil, err
	}

	if _, ok := scene.resourceLoaders[loaderID]; ok {
		return nil, NewErrorLoaderAlreadyExists(scene, loaderID)
	}

	resourceLoader, err := resources.NewResourceLoader(loaderID)

	if err != nil {
		return nil, err
	}

	scene.resourceLoaders[loaderID] = resourceLoader

	return resourceLoader, nil
}

// GetResourceLoader opens or returns
// an already opened resource loader.
func (scene *Scene) GetResourceLoader(fname string) (*resources.ResourceLoader, error) {
	loaderID, err := filepath.Abs(fname)

	if err != nil {
		return nil, err
	}

	loader, ok := scene.resourceLoaders[loaderID]

	if !ok {
		loader, err = scene.openResourceFile(loaderID)

		if err != nil {
			return nil, err
		}
	}

	return loader, nil
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
	err := scene.gmobs.VisitInOrder(func(key float64, gmobs map[*GameObject]struct{}) error {
		for gmob := range gmobs {
			err := gmob.Start()

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	currentSceneName = scene.name

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
	err = scene.gmobs.VisitInOrder(func(key float64, gmobs map[*GameObject]struct{}) error {
		for gmob := range gmobs {
			err = gmob.Update()

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
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
func (scene *Scene) insertGameObject(gmob *GameObject, zUpd float64) error {
	// Use binary search to insert the game
	// object into the Z-sorted update buffer.
	err := scene.gmobs.AddOrUpdate(zUpd,
		map[*GameObject]struct{}{gmob: {}},
		func(gmobs map[*GameObject]struct{}) (map[*GameObject]struct{}, error) {
			gmobs[gmob] = struct{}{}
			return gmobs, nil
		})

	if err != nil {
		return err
	}

	gmob.zUpdate = zUpd

	return nil
}

// removeGameObject deletes the game
// object from the Z update index.
func (scene *Scene) removeGameObject(gmob *GameObject) error {
	err := scene.gmobs.Update(gmob.zUpdate,
		func(gmobs map[*GameObject]struct{}, found bool) (map[*GameObject]struct{}, error) {
			if !found {
				return nil, RaiseErrorNoGameObjectOnScene(scene, gmob.name)
			}

			delete(gmobs, gmob)

			return gmobs, nil
		})

	if err != nil {
		return err
	}

	return nil
}

// placeGameObjects changes Z update coordinate
// of all the requested game objects.
func (scene *Scene) placeGameObjects() error {
	for _, changeZ := range scene.changeZBuffer {
		gmob := scene.FindGameObject(changeZ.gmobName)

		if gmob == nil {
			return RaiseErrorNoGameObjectOnScene(
				scene, changeZ.gmobName)
		}

		// Remove the game object from
		// its previous position.
		err := scene.removeGameObject(gmob)

		if err != nil {
			return err
		}

		// Insert the game object back
		// into the Z-buffer.
		err = scene.insertGameObject(gmob, changeZ.targetZ)

		if err != nil {
			return err
		}
	}

	return nil
}

// addBufferedGameObjects adds all the buffered
// game objects to the scene.
func (scene *Scene) addBufferedGameObjects() error {
	// Add all the game objects from the buffer
	// and start them all.
	for gmob, zUpd := range scene.addBuffer {
		err := scene.AddGameObject(gmob, zUpd)

		if err != nil {
			return err
		}

		err = gmob.Start()

		if err != nil {
			return err
		}
	}

	scene.addBuffer = map[*GameObject]float64{}

	return nil
}

// findDestroyedGameObject the game object with
// the specified name marked as destroyed.
func (scene *Scene) findDestroyedGameObject(name string) *GameObject {
	if gmob, ok := scene.gmobNameIndex[name]; ok {
		if gmob.destroyed {
			return gmob
		}
	}

	return nil
}

// removeDestroyedGameObject removes the game object
// with the specified name marked as destroyed.
func (scene *Scene) removeDestroyedGameObject(name string) error {
	gmob := scene.findDestroyedGameObject(name)

	if gmob == nil {
		return fmt.Errorf("scene '%s' doesn't have destroyed game object '%s'",
			scene.name, name)
	}

	err := scene.removeGameObject(gmob)

	if err != nil {
		return err
	}

	delete(scene.gmobNameIndex, gmob.name)

	return nil
}

// removeDestroyedGameObjects removes all the game
// objects set for destruction from the scene.
//
// ATTENTION: this method must not be called from
// any ecs.Component. Call it after scene.Update().
func (scene *Scene) removeDestroyedGameObjects() error {
	for gmob := range scene.destructionBuffer {
		err := scene.removeDestroyedGameObject(gmob)

		if err != nil {
			return err
		}
	}

	scene.destructionBuffer = map[string]struct{}{}

	return nil
}

// FindGameObject finds the game object on the scene.
func (scene *Scene) FindGameObject(name string) *GameObject {
	if gmob, ok := scene.gmobNameIndex[name]; !ok {
		return nil
	} else {
		if !gmob.destroyed {
			return gmob
		}

		return nil
	}
}

// findGameObjectInAdded searches for a game object
// with the specified name in the buffer where game objects
// set to be added reside.
func (scene *Scene) findGameObjectInAdded(gmob *GameObject) (bool, float64) {
	if zUpd, ok := scene.addBuffer[gmob]; ok {
		return true, zUpd
	}

	return false, -1
}

// HasGameObject returns true if the scene has a game
// object with the specified name, and false otherwise.
func (scene *Scene) HasGameObject(name string) bool {
	gmob := scene.FindGameObject(name)
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

	err := scene.insertGameObject(gmob, zUpd)

	if err != nil {
		return err
	}

	scene.gmobNameIndex[gmob.name] = gmob
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

	found, _ := scene.findGameObjectInAdded(gmob)

	if found {
		return fmt.Errorf("game object '%s' is already set to be added",
			gmob.name)
	}

	scene.addBuffer[gmob] = zUpd

	return nil
}

// ChangeGameObjectZ changes the game object Z.
func (scene *Scene) ChangeGameObjectZ(name string, zUpd float64) error {
	if !scene.HasGameObject(name) {
		return RaiseErrorNoGameObjectOnScene(scene, name)
	}

	zChange := changeZ{
		gmobName: name,
		targetZ:  zUpd,
	}

	scene.changeZBuffer = append(scene.changeZBuffer, zChange)

	return nil
}

// RemoveGameObject removes the game object from the scene.
func (scene *Scene) RemoveGameObject(name string) error {
	gmob := scene.FindGameObject(name)

	if gmob == nil {
		return RaiseErrorNoGameObjectOnScene(scene, name)
	}

	err := scene.removeGameObject(gmob)

	if err != nil {
		return err
	}

	delete(scene.gmobNameIndex, name)

	return nil
}

// DestroyGameObject destroys the game object, i.e.
// deactivates all its components and stops drawing it.
func (scene *Scene) DestroyGameObject(name string) error {
	gmob := scene.FindGameObject(name)

	if gmob == nil {
		return RaiseErrorNoGameObjectOnScene(scene, name)
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

	scene.destructionBuffer[name] = struct{}{}

	return nil
}

// DontDestroyOnSceneSwitch sets the game object
// under the specified name to be not destroyed on
// scene switch.
func (scene *Scene) DontDestroyOnSceneSwitch(gmobName string) error {
	gmob := scene.FindGameObject(gmobName)

	if gmob == nil {
		return RaiseErrorNoGameObjectOnScene(scene, gmobName)
	}

	if _, ok := noDestroyOnSceneSwitch[gmob.name]; ok {
		return NewErrorObjectAlreadyNotDestroyedOnSceneSwitch(gmob)
	}

	noDestroyOnSceneSwitch[gmob.name] = gmob

	return nil
}

// SwitchTo starts playing a different scene
// under the specified name.
func (scene *Scene) SwitchTo(sceneName string) error {
	otherScene, ok := scenes[sceneName]

	if !ok {
		return NewErrorSceneDoesntExist(sceneName)
	}

	for _, gmob := range noDestroyOnSceneSwitch {
		err := otherScene.AddGameObject(
			gmob, gmob.zUpdate)

		if err != nil {
			return err
		}
	}

	err := otherScene.Start()

	if err != nil {
		return err
	}

	return nil
}

// NewScene creates a new scene to
// place game objects onto.
func NewScene(name string, gmobsDictProducer collections.SortedDictionaryProducer[float64, map[*GameObject]struct{}]) (*Scene, error) {
	gmobs, err := gmobsDictProducer.Produce()

	if err != nil {
		return nil, err
	}

	return &Scene{
		name:              name,
		addBuffer:         map[*GameObject]float64{},
		changeZBuffer:     []changeZ{},
		gmobs:             gmobs,
		destructionBuffer: map[string]struct{}{},
		gmobNameIndex:     map[string]*GameObject{},
		systems:           map[string]System{},
		taskMgr:           tasking.NewTaskManager(),
		layout:            render.NewLayout(),
		resourceLoaders:   map[string]*resources.ResourceLoader{},
	}, nil
}
