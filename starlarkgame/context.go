package starlarkgame

type Context struct {
	player *PlayerState
	deck   Deck
}

func NewContext(deck Deck) *Context {
	return &Context{deck: deck, player: &PlayerState{}}
}

func (g *Context) GetPlayerState() *PlayerState {
	return g.player
}

func (g *Context) SetPlayerState(state *PlayerState) {
	g.player = state
}

func (g *Context) DeckPushCard(ref CardRef) {
	g.deck.Insert(CardNode{
		Card:     ref,
		Priority: MinCardPriority,
	})
	if g.player.CardRef == "" {
		g.player.CardRef = g.deck.Top()
	}
}

func (g *Context) DeckPopCard() {
	g.deck.Pop()
	g.player.CardRef = g.deck.Top()
}

func (g *Context) TopCard() CardRef {
	return g.deck.Top()
}
