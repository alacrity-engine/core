package ecs

import (
	"fmt"
)

// ErrorNoGameObjectOnScene is returned
// when there is no game object with
// certain name on the scene.
type ErrorNoGameObjectOnScene struct {
	scene    *Scene
	gmobName string
}

// Scene returns the scene which doesn't have the game object.
func (err *ErrorNoGameObjectOnScene) Scene() *Scene {
	return err.scene
}

// GameObject returns the name of the game object.
func (err *ErrorNoGameObjectOnScene) GameObject() string {
	return err.gmobName
}

// Error returns the error message.
func (err *ErrorNoGameObjectOnScene) Error() string {
	return fmt.Sprintf("scene '%s' doen't have game object '%s'",
		err.scene.Name(), err.gmobName)
}

// RaiseErrorNoGameObjectOnScene returns a new error
// about absence of the game object on the scene.
func RaiseErrorNoGameObjectOnScene(scene *Scene, gmobName string) *ErrorNoGameObjectOnScene {
	return &ErrorNoGameObjectOnScene{
		scene:    scene,
		gmobName: gmobName,
	}
}

/*****************************************************************************************************************/

// ErrorNoComponentOnGameObject is returned
// when there is no component with certain name
// on the game object.
type ErrorNoComponentOnGameObject struct {
	gmob          *GameObject
	componentName string
}

// GameObject returns the game object which
// doesn't have the component.
func (err *ErrorNoComponentOnGameObject) GameObject() *GameObject {
	return err.gmob
}

// Component returns the name of the component
// which is absent on the game object.
func (err *ErrorNoComponentOnGameObject) Component() string {
	return err.componentName
}

// Error returns the error message.
func (err *ErrorNoComponentOnGameObject) Error() string {
	return fmt.Sprintf("gamne object '%s' has no component '%s'",
		err.gmob.Name(), err.componentName)
}

// RaiseErrorNoComponentOnGameObject returns a new error
// about absence of the component on the game object.
func RaiseErrorNoComponentOnGameObject(gmob *GameObject, componentName string) *ErrorNoComponentOnGameObject {
	return &ErrorNoComponentOnGameObject{
		gmob:          gmob,
		componentName: componentName,
	}
}

/*****************************************************************************************************************/

// ErrorWrongComponentType is returned
// when the component is not of the type
// it should be.
type ErrorWrongComponentType struct {
	component    Component
	assertedType string
	actualType   string
}

// Component returns the component
// whose type was improperly asserted.
func (err *ErrorWrongComponentType) Component() Component {
	return err.component
}

// AssertedType returns the name of the type
// that was asserted for the component.
func (err *ErrorWrongComponentType) AssertedType() string {
	return err.assertedType
}

// ActualType returns the name of the
// actual type of the component.
func (err *ErrorWrongComponentType) ActualType() string {
	return err.actualType
}

// Error returns the error message.
func (err *ErrorWrongComponentType) Error() string {
	return fmt.Sprintf("component '%s' is of type '%s', but should be '%s'",
		err.component.Name(), err.actualType, err.assertedType)
}

// RaiseErrorWrongComponentType returns a new error
// about wrong type assertion of the component.
func RaiseErrorWrongComponentType(component Component, assertedType, actualType string) *ErrorWrongComponentType {
	return &ErrorWrongComponentType{
		component:    component,
		assertedType: assertedType,
		actualType:   actualType,
	}
}
