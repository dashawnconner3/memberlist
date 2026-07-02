package main

import "sync"

type NodeState int

const (
	StateAlive NodeState = iota
	StateSuspect
	StateDead
)

type Node struct {
	Name        string
	Incarnation uint32
	State       NodeState
}

type Memberlist struct {
	mu    sync.RWMutex
	nodes map[string]*Node
}

func (m *Memberlist) UpdateNode(name string, incarnation uint32, state NodeState) {
	m.mu.Lock()
	defer m.mu.Unlock()

	existing, ok := m.nodes[name]
	if !ok {
		m.nodes[name] = &Node{Name: name, Incarnation: incarnation, State: state}
		return
	}

	// Strict incarnation validation: 
	// If incoming incarnation is lower, ignore.
	// If equal, only allow transition if current is Dead and new is Alive (rejoin).
	if incarnation < existing.Incarnation {
		return
	}

	if incarnation == existing.Incarnation && state == StateAlive && existing.State == StateDead {
		// Allow rejoin if incarnation matches but state is transitioning from Dead to Alive
		existing.State = StateAlive
		return
	}

	existing.Incarnation = incarnation
	existing.State = state
}