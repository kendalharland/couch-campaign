package starlarkgame

type Context struct {
	ps PlayerState
}

func NewContext() *Context {
	return &Context{}
}

func (g *Context) GetPlayerState() PlayerState {
	return g.ps
}

func (g *Context) SetPlayerState(state PlayerState) {
	g.ps = state
}
