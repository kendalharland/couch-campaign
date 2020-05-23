package cardgame

import (
	"container/heap"
	"fmt"
	"io"
	"log"
)

var _ heap.Interface = &cardHeap{}

// Card priority.
//
// The lower the value the higher the priority.
const (
	maxCardPriority = iota
	infoCardPriority
	actionCardPriority
	minCardPriority
)

// Deck implements a min-heap of Card references.
//
// This is used to order the cards that are presented to the player. Each card
// is pushed onto the deck with an integer priority. Cards with equal priority
// are ordered based on their insertion time. Cards that have been in the deck
// longer than others with the same priority are popped first.
type Deck struct {
	h *cardHeap
}

// NewDeck ...
func NewDeck(cards []CardRef, cs CardService) *Deck {
	return &Deck{
		h: newCardHeap(cards, cs),
	}
}

// Clear ...
func (d *Deck) Clear() {
	d.h.nodes = []*cardHeapNode{}
}

// IsEmpty ...
func (d *Deck) IsEmpty() bool {
	return d.h.Len() == 0
}

// InsertCard ...
func (d *Deck) InsertCard(c CardRef, priority int) {
	heap.Push(d.h, &cardHeapNode{card: c, priority: priority})
	heap.Init(d.h)
}

// RemoveCard ...
func (d *Deck) RemoveCard(ref CardRef) {
	pos := -1
	for i, n := range d.h.nodes {
		if n.card == ref {
			pos = i
			break
		}
	}
	if pos >= 0 {
		heap.Remove(d.h, pos)
	}
}

// TopCard ...
func (d *Deck) TopCard() CardRef {
	if d.h.Len() > 0 {
		return d.h.Top().card
	}
	return ""
}

// PopTopCard ...
func (d *Deck) PopTopCard() CardRef {
	return heap.Pop(d.h).(*cardHeapNode).card
}

// DebugDump ...
func (d *Deck) DebugDump(w io.Writer) {
	d.h.debugDump(w)
}

type cardHeap struct {
	nodes             []*cardHeapNode
	lastInsertionTime int
}

type cardHeapNode struct {
	card          CardRef
	priority      int
	insertionTime int
}

func newCardHeap(cards []CardRef, cs CardService) *cardHeap {
	h := &cardHeap{}
	for _, card := range cards {
		h.nodes = append(h.nodes, &cardHeapNode{
			card:          card,
			priority:      cs.CardPriorityByType(card),
			insertionTime: h.lastInsertionTime,
		})
	}
	h.lastInsertionTime++
	return h
}

func (h *cardHeap) Init() {
	heap.Init(h)
}

func (h *cardHeap) Push(x interface{}) {
	n := x.(*cardHeapNode)
	n.insertionTime = h.lastInsertionTime
	h.lastInsertionTime++
	h.nodes = append(h.nodes, n)
}

func (h *cardHeap) Pop() interface{} {
	old := h.nodes
	n := len(old)
	item := old[n-1]
	h.nodes = old[0 : n-1]
	return item
}

func (h *cardHeap) Len() int {
	return len(h.nodes)
}

func (h *cardHeap) Less(i, j int) bool {
	ni, nj := h.nodes[i], h.nodes[j]
	return ni.priority < nj.priority ||
		ni.priority == nj.priority && ni.insertionTime < nj.insertionTime
}

func (h *cardHeap) Swap(i, j int) {
	h.nodes[i], h.nodes[j] = h.nodes[j], h.nodes[i]
}

func (h *cardHeap) Top() *cardHeapNode {
	if h.Len() == 0 {
		log.Fatal("deck is empty")
	}
	return h.nodes[0]
}

func (h *cardHeap) debugDump(w io.Writer) {
	for _, c := range h.nodes {
		fmt.Fprintf(w, "%+v\n", c)
	}
}
