package starlarkgame

type Context struct {
	state *State
	deck  Deck
}

func NewContext(deck Deck) *Context {
	return &Context{deck: deck, state: &State{}}
}

func (g *Context) GetState() *State {
	return g.state
}

func (g *Context) SetState(state *State) {
	g.state = state
}

func (g *Context) DeckPushCard(ref CardRef) {
	g.deck.Insert(CardNode{
		Card:     ref,
		Priority: MinCardPriority,
	})
}

func (g *Context) DeckPopCard() {
	g.deck.Pop()
}

func (g *Context) DeckTopCard() CardRef {
	return g.deck.Top()
}
