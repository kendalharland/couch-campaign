package couchcampaign

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
