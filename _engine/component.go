package engine

// Component is a single script
// which contains data and instructions.
// Updated once per frame.
type Component interface {
	Name() string
	SetName(string)
	Start() error
	Update() error
	Destroy() error
	GameObject() *GameObject
	SetGameObject(*GameObject)
	Active() bool
	SetActive(bool)
}

type RegisteredComponent interface {
	TypeID() string
}

// BaseComponent is the base type
// to be included into any component.
type BaseComponent struct {
	typeID string
	name   string
	gmob   *GameObject
	active bool
}

// Name returns the name of the component.
func (bc *BaseComponent) Name() string {
	return bc.name
}

// SetName sets the name of the component.
func (bc *BaseComponent) SetName(name string) {
	bc.name = name
}

// GameObject returns the game object the component
// is currently attached to.
func (bc *BaseComponent) GameObject() *GameObject {
	return bc.gmob
}

// SetGameObject changes the game object of the component.
func (bc *BaseComponent) SetGameObject(gmob *GameObject) {
	bc.gmob = gmob
}

// Active indicates if the component is
// currently active.
//
// If the component is not active, it's
// Update method is not being called in
// the application loop.
func (bc *BaseComponent) Active() bool {
	return bc.active
}

// SetActive changes the activity status of the component.
func (bc *BaseComponent) SetActive(active bool) {
	bc.active = active
}
