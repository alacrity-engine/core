package tasking

type (
	// TaskManager holds tasks
	// that are being performed
	// throughout frames.
	TaskManager struct {
		startedTasks []*task
		tasks        []*task
		stoppedTasks []string
		afterTasks   map[string][]*task
	}

	// task is a function
	// being performed
	// throughout frames.
	task struct {
		name    string
		fn      func() (bool, error)
		stopped bool
	}
)

// addStartedTasks transfers all the newly
// started tasks from the buffer to the
// array of tasks.
func (mgr *TaskManager) addStartedTasks() {
	mgr.tasks = append(mgr.tasks,
		mgr.startedTasks...)
	mgr.startedTasks = []*task{}
}

// findStartedTask searches for task in the
// buffer of the newly started tasks.
func (mgr *TaskManager) findStartedTask(name string) (int, *task) {
	ind := -1
	var t *task

	for i, tsk := range mgr.startedTasks {
		if tsk.name == name {
			ind = i
			t = tsk

			break
		}
	}

	return ind, t
}

// findTask finds a currently being performed task
// in the array of tasks.
func (mgr *TaskManager) findTask(name string) (int, *task) {
	ind := -1
	var t *task

	for i, tsk := range mgr.tasks {
		if tsk.name == name && !tsk.stopped {
			ind = i
			t = tsk

			break
		}
	}

	return ind, t
}

// findStoppedTask searches for task in the buffer
// where all the stopped tasks reside.
func (mgr *TaskManager) findStoppedTask(name string) (int, *task) {
	ind := -1
	var t *task

	for i, tsk := range mgr.tasks {
		if tsk.name == name && tsk.stopped {
			ind = i
			t = tsk

			break
		}
	}

	return ind, t
}

// removeStoppedTask removes the stopped task
// from the array of tasks.
func (mgr *TaskManager) removeStoppedTask(name string) error {
	i, stoppedTask := mgr.findStoppedTask(name)

	if stoppedTask == nil {
		return NewErrorStoppedTaskNotExists(name)
	}

	mgr.tasks = append(mgr.tasks[:i], mgr.tasks[i+1:]...)

	return nil
}

// removeStoppedTasks removes all the stopped tasks
// from the array of tasks.
func (mgr *TaskManager) removeStoppedTasks() error {
	for _, taskName := range mgr.stoppedTasks {
		err := mgr.removeStoppedTask(taskName)

		if err != nil {
			return err
		}
	}

	mgr.stoppedTasks = []string{}

	return nil
}

// Start starts the task manager.
func (mgr *TaskManager) Start() error {
	return nil
}

// Update calls all the active tasks once.
func (mgr *TaskManager) Update() error {
	err := mgr.removeStoppedTasks()

	if err != nil {
		return err
	}

	mgr.addStartedTasks()

	for _, tsk := range mgr.tasks {
		if tsk.stopped {
			continue
		}

		shouldContinue, err := tsk.fn()

		if err != nil {
			return err
		}

		if !shouldContinue {
			err = mgr.StopTask(tsk.name)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Destroy stops all the tasks and
// removes them from the task manager.
func (mgr *TaskManager) Destroy() error {
	for _, tsk := range mgr.tasks {
		tsk.stopped = true
	}

	return mgr.removeStoppedTasks()
}

// HasTask returns true if the task with the
// specified name exists in the array of
// currently being performed tasks.
func (mgr *TaskManager) HasTask(name string) bool {
	_, tsk := mgr.findTask(name)

	return tsk != nil
}

// StartTask starts a new task with the
// specified name and function.
func (mgr *TaskManager) StartTask(name string, fn func() (bool, error)) error {
	if mgr.HasTask(name) {
		return NewErrorTaskAlreadyExists(name)
	}

	_, tsk := mgr.findStartedTask(name)

	if tsk != nil {
		return NewErrorTaskAlreadyStarted(name)
	}

	t := &task{
		name:    name,
		fn:      fn,
		stopped: false,
	}

	mgr.startedTasks = append(mgr.startedTasks, t)

	return nil
}

// StartTaskAfter starts a new task when the
// specified task stops.
func (mgr *TaskManager) StartTaskAfter(afterName, taskName string, fn func() (bool, error)) error {
	if mgr.HasTask(taskName) {
		return NewErrorTaskAlreadyExists(taskName)
	}

	_, tsk := mgr.findStartedTask(taskName)

	if tsk != nil {
		return NewErrorTaskAlreadyStarted(taskName)
	}

	if _, ok := mgr.afterTasks[afterName]; !ok {
		mgr.afterTasks[afterName] = []*task{}
	}

	// Check if the after task already exists.
	found := false

	for _, afterTask := range mgr.afterTasks[afterName] {
		if afterTask.name == taskName {
			found = true
			break
		}
	}

	if found {
		return NewErrorAfterTaskAlreadyStarted(taskName)
	}

	tsk = &task{
		name:    taskName,
		fn:      fn,
		stopped: false,
	}
	mgr.afterTasks[afterName] = append(mgr.afterTasks[afterName], tsk)

	return nil
}

// StopTask stops the task with the specified name.
func (mgr *TaskManager) StopTask(name string) error {
	_, tsk := mgr.findTask(name)

	if tsk == nil {
		return NewErrorTaskNotExists(name)
	}

	mgr.stoppedTasks = append(mgr.stoppedTasks, name)
	tsk.stopped = true

	// Check if there are tasks that must be started
	// after the stopped one.
	if afterTasks, ok := mgr.afterTasks[tsk.name]; ok {
		for _, afterTask := range afterTasks {
			err := mgr.StartTask(afterTask.name, afterTask.fn)

			if err != nil {
				return err
			}
		}
	}

	mgr.afterTasks[tsk.name] = []*task{}

	return nil
}

// NewTaskManager returns a new task manager
// to hold active tasks performed through frames.
func NewTaskManager() *TaskManager {
	taskMgr := &TaskManager{
		startedTasks: []*task{},
		tasks:        []*task{},
		stoppedTasks: []string{},
		afterTasks:   map[string][]*task{},
	}

	return taskMgr
}
