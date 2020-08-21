package couchcampaign

import (
	"container/heap"
	"log"
	"math"
)

// CardRef is a reference to a card.
type CardRef = string

// CardPriority determines how soon a card is popped from a Deck.
//
// The lower the value the higher the priority.
type CardPriority = int64

// Card priority constants.
const (
	MaxCardPriority CardPriority = math.MinInt64
	MinCardPriority CardPriority = math.MaxInt64
)

// Deck implements a min-heap of Card references.
//
// This is used to order the cards that are presented to the player. Each card
// is pushed onto the deck with an integer priority. Cards with equal priority
// are ordered based on their insertion time. Cards that have been in the deck
// longer than others with the same priority are popped first.
type Deck struct {
	h *cardMinHeap
}

func newDeck() *Deck {
	return &Deck{h: newCardHeap()}
}

// Clear empties the deck.
func (d *Deck) Clear() {
	d.h.nodes = []*CardNode{}
}

// IsEmpty returns true if the deck has no cards.
func (d *Deck) IsEmpty() bool {
	return d.h.Len() == 0
}

// Insert adds a card to the deck.
func (d *Deck) Insert(n CardNode) {
	heap.Push(d.h, &n)
	heap.Init(d.h)
}

// Remove removes a card from the deck.
func (d *Deck) Remove(c CardRef) CardRef {
	pos := -1
	for i, n := range d.h.nodes {
		if n.Card == c {
			pos = i
			break
		}
	}
	if pos >= 0 {
		node := heap.Remove(d.h, pos).(*CardNode)
		return node.Card
	}
	return ""
}

// Top returns the top card from the deck.
func (d *Deck) Top() CardRef {
	if d.h.Len() > 0 {
		return d.h.Top().Card
	}
	return ""
}

// Pop removes and returns the top card from the deck.
func (d *Deck) Pop() CardRef {
	return heap.Pop(d.h).(*CardNode).Card
}

// CardNode represents a Card stored in a Deck.
type CardNode struct {
	Card          CardRef
	Priority      int64
	insertionTime int
}

var _ heap.Interface = &cardMinHeap{}

type cardMinHeap struct {
	nodes             []*CardNode
	lastInsertionTime int
}

func newCardHeap() *cardMinHeap {
	return &cardMinHeap{
		nodes: []*CardNode{},
	}
}

func (h *cardMinHeap) Init() {
	heap.Init(h)
}

func (h *cardMinHeap) Push(x interface{}) {
	n := x.(*CardNode)
	n.insertionTime = h.lastInsertionTime
	h.lastInsertionTime++
	h.nodes = append(h.nodes, n)
}

func (h *cardMinHeap) Pop() interface{} {
	old := h.nodes
	n := len(old)
	item := old[n-1]
	h.nodes = old[0 : n-1]
	return item
}

func (h *cardMinHeap) Len() int {
	return len(h.nodes)
}

func (h *cardMinHeap) Less(i, j int) bool {
	ni, nj := h.nodes[i], h.nodes[j]
	return ni.Priority < nj.Priority ||
		ni.Priority == nj.Priority && ni.insertionTime < nj.insertionTime
}

func (h *cardMinHeap) Swap(i, j int) {
	h.nodes[i], h.nodes[j] = h.nodes[j], h.nodes[i]
}

func (h *cardMinHeap) Top() *CardNode {
	if h.Len() == 0 {
		log.Fatal("deck is empty")
	}
	return h.nodes[0]
}
