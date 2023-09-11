package scripting

import "github.com/alacrity-engine/core/engine"

// ComponentTypeEntry contains the info
// about the game component type.
type ComponentTypeEntry struct {
	Name        string                             // Name is a name of the component.
	PkgPath     string                             // PkgPath is a path to the Go package where the component is located.
	Fields      map[string]ComponentTypeFieldEntry // Fields is a collection of all the component's exported fields.
	Constructor func() engine.Component            // Constructor is a function that returns an empty component of the type.
}

// ComponentTypeFieldEntry contains the info
// about the component's fields.
type ComponentTypeFieldEntry struct {
	Name   string                                         // Name is a name of the field.
	Type   string                                         // Type is a type of the field as how it would be written in Go code.
	Getter func(comp engine.Component) interface{}        // Getter is a function to obtain a value of the field.
	Setter func(comp engine.Component, value interface{}) // Setter is a function to change the value of the field.
}
