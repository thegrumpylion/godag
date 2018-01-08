package dag

import (
	"errors"
)

// Node represents a node in the dag
type Node struct {
	ID       string
	Value    interface{}
	Children []*Node
	In       int
}

// NodeCBFunc callback function for traversing
type NodeCBFunc func(node *Node) error

// WalkType determines the kind of the traversal
type WalkType int

const (
	// WalkBreadthFirst flag for walk function
	WalkBreadthFirst WalkType = iota
	// WalkDepthFirst flag for walk function
	WalkDepthFirst
)

// IsLeaf returns true if node is a leaf
func (node *Node) IsLeaf() bool {
	return len(node.Children) == 0
}

// IsRoot returns true if node is a root
func (node *Node) IsRoot() bool {
	return node.In == 0
}

// Walk traverses the subgraph from a given node
func (node *Node) Walk(wt WalkType, cb NodeCBFunc) error {
	if wt == WalkBreadthFirst {
		return walkBreadthFirst(node, cb)
	}
	// for DFS we explicitly visit the root node first and then
	// let the recursive function do is thing
	if err := cb(node); err != nil {
		return err
	}
	return walkDepthFirst(node, cb)
}

func walkDepthFirst(node *Node, cb NodeCBFunc) error {
	for _, child := range node.Children {
		if err := cb(child); err != nil {
			return err
		}
		if !child.IsLeaf() {
			if err := walkDepthFirst(child, cb); err != nil {
				return err
			}
		}
	}
	return nil
}

// minimal queue implementation for BFS using a slice
type queue struct {
	nodes []*Node
}

func (q *queue) push(node *Node) {
	q.nodes = append(q.nodes, node)
}

func (q *queue) pop() *Node {
	node := q.nodes[0]
	q.nodes = q.nodes[1:]
	return node
}

func (q *queue) isEmpty() bool {
	return len(q.nodes) == 0
}

// BreadthFirstSearch algorithm
// queue has more ? -> pop -> visit -> push children in queue -> repeat
func walkBreadthFirst(node *Node, cb NodeCBFunc) error {

	q := queue{}
	// push the root node
	q.push(node)

	for {
		if q.isEmpty() {
			return nil
		}
		currNode := q.pop()
		if err := cb(currNode); err != nil {
			return err
		}
		if !currNode.IsLeaf() {
			for _, child := range currNode.Children {
				q.push(child)
			}
		}
	}
}

// Dag the main object of the graph
type Dag struct {
	nodes    map[string]*Node
	revEdges map[string]map[string]struct{}
}

// AddNode creates and adds a new node to the dag and returns the instance
func (dag *Dag) AddNode(id string, value interface{}) (*Node, error) {
	if _, inMap := dag.nodes[id]; inMap {
		return nil, errors.New("A node with id " + id + " already exist")
	}
	node := &Node{
		ID:    id,
		Value: value,
	}
	dag.nodes[id] = node
	return node, nil
}

// AddEdge adds an edge between nodes
func (dag *Dag) AddEdge(from, to string) error {
	// chack if we have reverse edges for the "to" node
	if fromEdges, inMap := dag.revEdges[to]; inMap {
		// check if edge exists
		if _, fromExist := fromEdges[from]; fromExist {
			return errors.New("Edge from: " + from + " to: " + to + " already exist")
		}
	} else {
		// initialize the inner map
		dag.revEdges[to] = map[string]struct{}{}
	}

	dag.revEdges[to][from] = struct{}{}
	fromNode := dag.nodes[from]
	toNode := dag.nodes[to]
	fromNode.Children = append(fromNode.Children, toNode)
	toNode.In++

	return nil
}

// Roots returns all roots in the dag, i.e. no in edges
func (dag *Dag) Roots() []*Node {
	roots := []*Node{}
	for _, node := range dag.nodes {
		if node.IsRoot() {
			roots = append(roots, node)
		}
	}
	return roots
}

// IsAcyclic checks the Dag for loops
func (dag *Dag) IsAcyclic() bool {

	roots := dag.Roots()

	for _, root := range roots {
		visitedSet := map[string]struct{}{}
		err := root.Walk(WalkDepthFirst, func(node *Node) error {
			if _, visited := visitedSet[node.ID]; visited {
				return errors.New("The graph has cycles")
			}
			visitedSet[node.ID] = struct{}{}
			return nil
		})
		if err != nil {
			return false
		}
	}
	return true
}

// New creates new dag
func New() *Dag {
	dag := &Dag{}
	dag.nodes = make(map[string]*Node)
	dag.revEdges = make(map[string]map[string]struct{})
	return dag
}
