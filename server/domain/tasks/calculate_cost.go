package tasks

import (
	"github.com/google/uuid"
)

type CalculateCost struct{}

func NewCalculateCost() *CalculateCost { return &CalculateCost{} }

type tree struct {
	Nodes []node
}

type node struct {
	Task     Task
	Children []node
}

func buildTree(tasks []Task) tree {
	taskMap := make(map[uuid.UUID]Task)
	childrenMap := make(map[uuid.UUID][]uuid.UUID)
	rootTasks := make([]Task, 0)
	for _, task := range tasks {
		taskMap[task.ID] = task

		if task.ParentID == uuid.Nil {
			rootTasks = append(rootTasks, task)
			continue
		}

		childrenMap[task.ParentID] = append(childrenMap[task.ParentID], task.ID)

	}

	rootNodes := make([]node, 0, len(rootTasks))
	for _, rootTask := range rootTasks {
		rootNodes = append(rootNodes, buildNode(rootTask, childrenMap, taskMap))
	}

	return tree{Nodes: rootNodes}
}

func buildNode(task Task, childMap map[uuid.UUID][]uuid.UUID, taskMap map[uuid.UUID]Task) node {
	childIDs := childMap[task.ID]
	children := make([]node, 0, len(childIDs))
	totalCost := task.Cost

	for _, childID := range childIDs {
		child := buildNode(taskMap[childID], childMap, taskMap)
		children = append(children, child)
		totalCost += child.Task.Cost
	}

	task.Cost = totalCost

	return node{
		Task:     task,
		Children: children,
	}
}

func flatten(t tree) []Task {
	// Allocate for atleast the size of root nodes
	tasks := make([]Task, 0, len(t.Nodes))
	stack := make([]node, 0, len(t.Nodes))

	// Populate the stack with root nodes
	for _, n := range t.Nodes {
		stack = append(stack, n)
	}

	for len(stack) > 0 {

		// pop from the stack
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		tasks = append(tasks, current.Task)

		// Push all the children to the stack
		for _, child := range current.Children {
			stack = append(stack, child)
		}
	}

	return tasks
}

func (c *CalculateCost) Run(tasks []Task) []Task {
	return flatten(buildTree(tasks))
}
