package tasks

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/tekofx/crossposter/internal/logger"
)

type TasksManager struct {
	mu    sync.Mutex
	tasks map[string]context.CancelFunc
}

func newTasksManager() *TasksManager {
	return &TasksManager{
		tasks: make(map[string]context.CancelFunc),
	}
}

// StartTask starts a new goroutine and stores its CancelFunc
func (tm *TasksManager) StartTask(id string, taskFunc func(context.Context)) {
	ctx, cancel := context.WithCancel(context.Background())

	tm.mu.Lock()
	tm.tasks[id] = cancel
	tm.mu.Unlock()

	go func() {
		defer cancel()
		taskFunc(ctx)
	}()
}

func (tm *TasksManager) GetAllTasks() string {
	var msg strings.Builder
	for id, _ := range tm.tasks {
		msg.WriteString(fmt.Sprintf("%s\n", id))
	}
	return msg.String()
}

// StopTask stops a specific goroutine by ID
func (tm *TasksManager) StopTask(id string) {
	logger.Log("Stopping task", id)
	tm.mu.Lock()
	cancel, exists := tm.tasks[id]
	if exists {
		cancel()             // Trigger cancellation
		delete(tm.tasks, id) // Clean up
	}
	tm.mu.Unlock()
}

// StopAll stops all running tasks
func (tm *TasksManager) StopAll() {
	tm.mu.Lock()
	for id, cancel := range tm.tasks {
		cancel()
		delete(tm.tasks, id)
	}
	tm.mu.Unlock()
}
