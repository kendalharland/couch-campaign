package starlarkgame

type Context struct {
	ps   PlayerState
	deck Deck
}

func NewContext(deck Deck) *Context {
	return &Context{deck: deck}
}

func (g *Context) GetPlayerState() PlayerState {
	return g.ps
}

func (g *Context) SetPlayerState(state PlayerState) {
	g.ps = state
}

func (g *Context) DeckPushCard(ref CardRef) {
	g.deck.Insert(CardNode{
		Card:     ref,
		Priority: MinCardPriority,
	})
	if g.ps.CardRef == "" {
		g.ps.CardRef = g.deck.Top()
	}
}

func (g *Context) DeckPopCard() {
	g.deck.Pop()
	g.ps.CardRef = g.deck.Top()
}
