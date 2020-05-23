package cardgame

type CardRef = string

type cardType int

const (
	actionCardType cardType = iota
	infoCardType
	votingCardType
)

type Card struct {
	ID         CardRef `json:"id"`
	Image      string  `json:"image"`
	Speaker    string  `json:"speaker"`
	Text       string  `json:"text"`
	AcceptText string  `json:"accept_text"`
	RejectText string  `json:"reject"`
	OnShown    func(GameLogicApi, PID, CardRef)
	OnAccept   func(GameLogicApi, PID, CardRef)
	OnReject   func(GameLogicApi, PID, CardRef)
}
