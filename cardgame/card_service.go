package cardgame

type CardService interface {
	// TODO: This should not deal with Decks.
	BuildIntroDeck() *Deck
	BuildBaseDeck() *Deck
	BuildSocietyCrumbledDeck() *Deck
	PlayerCard(pid PID) CardRef
	Card(CardRef) Card
	CardRequiresInput(CardRef) bool
	CardPriorityByType(CardRef) int
}
