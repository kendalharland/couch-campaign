
type CardService struct {
	Cards []Card
}

func (s *CardService) BuildBaseDeck() *Deck {
	cards := make([]CardRef, len(s.Cards))
	for _, card := range s.Cards {
		cards = append(cards, card.ID)
	}
	deck := NewDeck(cards, s.ards)
	return deck
}

func (s *CardService) Card(ref CardRef) Card {
	for _, card := range s.Cards {
		if card.ID == ref {
			return card
		}
	}

	log.Fatalf("missing card %v", ref)
	return Card{} // never gets here.
}

func (s *CardService) CardTypeOf(c Card) cardType {
	switch {
	case c.ID == TheVotingCard.ID:
		return votingCardType
	case c.AcceptText == "" && c.RejectText == "":
		return infoCardType
	default:
		return actionCardType
	}
}

func (s *CardService) CardPriorityByType(ref CardRef) int {
	switch s.CardTypeOf(s.Card(ref)) {
	case actionCardType:
		return actionCardPriority
	case infoCardType, votingCardType:
		return infoCardPriority
	default:
		return minCardPriority
	}
}

func (s *CardService) CardRequiresInput(c Card) bool {
	return s.CardTypeOf(c) != votingCardType
}
