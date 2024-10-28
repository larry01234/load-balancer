package loadbalancer

import (
	"sync"
	"time"
)

type Node struct {
	Address string
	BPM     int // Bytes Per Minute
	RPM     int // Requests Per Minute

	mu           sync.Mutex
	currentBPM   int
	currentRPM   int
	lastResetBPM time.Time
	lastResetRPM time.Time
}

func NewNode(address string, bpm int, rpm int) *Node {
	return &Node{
		Address:      address,
		BPM:          bpm,
		RPM:          rpm,
		lastResetBPM: time.Now(),
		lastResetRPM: time.Now(),
	}
}

func (n *Node) AllowRequest(bodySize int) bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	now := time.Now()

	// Reset BPM counter if a minute has passed
	if now.Sub(n.lastResetBPM) >= time.Minute {
		n.currentBPM = 0
		n.lastResetBPM = now
	}

	// Reset RPM counter if a minute has passed
	if now.Sub(n.lastResetRPM) >= time.Minute {
		n.currentRPM = 0
		n.lastResetRPM = now
	}

	// Check if adding this request would exceed BPM or RPM
	if n.currentBPM+bodySize > n.BPM || n.currentRPM+1 > n.RPM {
		return false
	}

	// Update counters
	n.currentBPM += bodySize
	n.currentRPM += 1
	return true
}