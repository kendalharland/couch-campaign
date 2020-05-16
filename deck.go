package couchcampaign

import (
	"container/heap"
	"fmt"
	"io"
	"log"
	"math/rand"
)

var _ heap.Interface = &cardHeap{}

const (
	deckURL     = "https://gist.githubusercontent.com/kendalharland/71a3b843f740099706b148e62a2dd8eb/raw/couch-campaign-cards"
	deckMaxSize = 1000
)

// Deck ...
type Deck struct {
	h *cardHeap
}

// NewDeck ...
func NewDeck(cards []Card) *Deck {
	return &Deck{
		h: newCardHeap(cards),
	}
}

// Clear ...
func (d *Deck) Clear() {
	d.h = newCardHeap([]Card{})
}

// IsEmpty ...
func (d *Deck) IsEmpty() bool {
	return d.h.Len() == 0
}

// InsertCard ...
func (d *Deck) InsertCard(c Card) {
	heap.Push(d.h, &cardHeapNode{card: c, priority: cardPriorityByType(c)})
	heap.Init(d.h)
}

// InsertCardWithPriority ...
func (d *Deck) InsertCardWithPriority(c Card, priority int) {
	heap.Push(d.h, &cardHeapNode{card: c, priority: priority})
	heap.Init(d.h)
}

// RemoveCard ...
func (d *Deck) RemoveCard(c Card) {
	pos := -1
	_ = d.h
	_ = d.h.nodes

	for i, n := range d.h.nodes {
		if n.card == c {
			pos = i
			break
		}
	}
	if pos >= 0 {
		heap.Remove(d.h, pos)
	}
}

// TopCard ...
func (d *Deck) TopCard() Card {
	if d.h.Len() > 0 {
		return d.h.Top().card
	}
	return nil
}

// PopTopCard ...
func (d *Deck) PopTopCard() Card {
	return heap.Pop(d.h).(*cardHeapNode).card
}

// ShuffleActionCards ...
func (d *Deck) ShuffleActionCards() {
	for _, n := range d.h.nodes {
		if _, ok := n.card.(ActionCard); !ok {
			continue
		}
		n.priority = rand.Intn(actionCardPriority + 1)
		n.insertionTime = d.h.lastInsertionTime
	}
	heap.Init(d.h)
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
	card          Card
	priority      int
	insertionTime int
}

func newCardHeap(cards []Card) *cardHeap {
	h := &cardHeap{}
	for _, card := range cards {
		h.nodes = append(h.nodes, &cardHeapNode{
			card:          card,
			priority:      cardPriorityByType(card),
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

// Card priority.
//
// The lower the value the higher the priority.
const (
	maxCardPriority = iota
	infoCardPriority
	actionCardPriority
	minCardPriority
)

func cardPriorityByType(c Card) int {
	switch c.(type) {
	case actionCard:
		return actionCardPriority
	case infoCard, votingCard:
		return infoCardPriority
	default:
		return minCardPriority
	}
}
