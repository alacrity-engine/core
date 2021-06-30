package tasking

import "fmt"

// ErrorTaskAlreadyExists is raised
// when the task with the specified
// name already exists in the task manager.
type ErrorTaskAlreadyExists struct {
	taskName string
}

// Error returns the error message.
func (err *ErrorTaskAlreadyExists) Error() string {
	return fmt.Sprintf("task '%s' already exists in the task manager",
		err.taskName)
}

// NewErrorTaskAlreadyExists creates a new
// error of type ErrorTaskAlreadyExists.
func NewErrorTaskAlreadyExists(taskName string) *ErrorTaskAlreadyExists {
	return &ErrorTaskAlreadyExists{
		taskName: taskName,
	}
}

/*******************************************************************************/

// ErrorTaskNotExists is returned when the task
// with the specified name doesn't exist in the
// task manager.
type ErrorTaskNotExists struct {
	taskName string
}

// Error returns the error message.
func (err *ErrorTaskNotExists) Error() string {
	return fmt.Sprintf("task '%s' doesn't exist in the task manager",
		err.taskName)
}

// NewErrorTaskNotExists creates a new error
// of type ErrorTaskNotExists.
func NewErrorTaskNotExists(taskName string) *ErrorTaskNotExists {
	return &ErrorTaskNotExists{
		taskName: taskName,
	}
}

/*******************************************************************************/

// ErrorStoppedTaskNotExists is raised when
// there's no stopped task with the specified name.
type ErrorStoppedTaskNotExists struct {
	taskName string
}

// Error returns the error message.
func (err *ErrorStoppedTaskNotExists) Error() string {
	return fmt.Sprintf("stopped task '%s' doesn't exist in the task manager",
		err.taskName)
}

// NewErrorStoppedTaskNotExists returns a new error
// of type ErrorStoppedTaskNotExists.
func NewErrorStoppedTaskNotExists(taskName string) *ErrorStoppedTaskNotExists {
	return &ErrorStoppedTaskNotExists{
		taskName: taskName,
	}
}

/*******************************************************************************/

// ErrorTaskAlreadyStarted is raised when the
// task with the specified name has already been
// started.
type ErrorTaskAlreadyStarted struct {
	taskName string
}

// Errorreturns the error message.
func (err *ErrorTaskAlreadyStarted) Error() string {
	return fmt.Sprintf("task '%s' has already been started",
		err.taskName)
}

// NewErrorTaskAlreadyStarted returns a new error
// of type ErrorTaskAlreadyStarted,
func NewErrorTaskAlreadyStarted(taskName string) *ErrorTaskAlreadyStarted {
	return &ErrorTaskAlreadyStarted{
		taskName: taskName,
	}
}

/*******************************************************************************/

// ErrorAfterTaskAlreadyStarted is raised when the
// task with the specified name has already been
// started.
type ErrorAfterTaskAlreadyStarted struct {
	taskName string
}

// Errorreturns the error message.
func (err *ErrorAfterTaskAlreadyStarted) Error() string {
	return fmt.Sprintf("after task '%s' has already been started",
		err.taskName)
}

// NewErrorTaskAlreadyStarted returns a new error
// of type ErrorTaskAlreadyStarted,
func NewErrorAfterTaskAlreadyStarted(taskName string) *ErrorAfterTaskAlreadyStarted {
	return &ErrorAfterTaskAlreadyStarted{
		taskName: taskName,
	}
}
