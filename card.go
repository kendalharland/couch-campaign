package couchcampaign

import "log"

// Cards are stored in a google spreadsheet and exported in TSV format.
// these constants are the array indices of the card fields in the TSV file.
const (
	cardID = iota
	cardAdvisor
	cardContent
	cardAcceptText
	cardAWE
	cardAHE
	cardASE
	cardRejectText
	cardRWE
	cardRHE
	cardRSE

	cardFieldCount
)

type cardType int

const (
	actionCardType cardType = iota
	infoCardType
	votingCardType
)

func cardTypeOf(c Card) cardType {
	switch t := c.(type) {
	case actionCard:
		return actionCardType
	case infoCard:
		return infoCardType
	case votingCard:
		return votingCardType
	default:
		log.Fatalf("cardTypeOf: unknown card type %v: %v", t, c)
		return -1 // satisfies the analyzer.
	}
}

type Card interface{}

type actionCard struct {
	ID                 int
	Advisor            string
	Content            string
	AcceptText         string
	AccWealthEffect    int
	AccHealthEffect    int
	AccStabilityEffect int
	RejectText         string
	RejWealthEffect    int
	RejHealthEffect    int
	RejStabilityEffect int
}

func (c actionCard) toProto() *ActionCard {
	return &ActionCard{
		Advisor:    c.Advisor,
		Content:    c.Content,
		AcceptText: c.AcceptText,
		RejectText: c.RejectText,
	}
}

func (c actionCard) accept(s *stats) {
	s.Wealth += c.AccWealthEffect
	s.Health += c.AccHealthEffect
	s.Stability += c.AccStabilityEffect
}

func (c actionCard) reject(s *stats) {
	s.Wealth += c.RejWealthEffect
	s.Health += c.RejHealthEffect
	s.Stability += c.RejStabilityEffect
}

type infoCard struct {
	Info string
}

func (c infoCard) toProto() *InfoCard {
	return &InfoCard{Text: c.Info}
}

var theVotingCard = votingCard("waiting for votes...")

type votingCard string

func (c votingCard) toProto() *VotingCard {
	return &VotingCard{Text: string(c)}
}
