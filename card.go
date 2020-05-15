package couchcampaign

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

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

func loadBaseCards() ([]Card, error) {
	res, err := http.Get(deckURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var cards []Card
	// Parse the deck. It's stored in tab-separated-values (TVS) format. Each
	// card is on a new line.
	tsvLines := strings.SplitN(string(body), "\n", deckMaxSize)
	for _, line := range tsvLines {
		cardFields := strings.SplitN(line, "\t", cardFieldCount)
		c, err := parseCardTSV(cardFields)
		if err != nil {
			return nil, fmt.Errorf("invalid card: %w", err)
		}
		cards = append(cards, c)
	}

	return cards, nil
}

func parseCardTSV(fields []string) (Card, error) {
	parseInt := func(input string) (int, error) {
		id, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return -1, err
		}
		return int(id), nil
	}

	var c Card

	for i := range fields {
		fields[i] = strings.TrimSpace(fields[i])
	}

	id, err := parseInt(fields[cardID])
	if err != nil {
		return c, fmt.Errorf("invalid id: %w", err)
	}

	awe, err := parseInt(fields[cardAWE])
	if err != nil {
		return c, fmt.Errorf("invalid ams: %w", err)
	}
	ahe, err := parseInt(fields[cardAHE])
	if err != nil {
		return c, fmt.Errorf("invalid aps: %w", err)
	}
	ase, err := parseInt(fields[cardASE])
	if err != nil {
		return c, fmt.Errorf("invalid ajs: %w", err)
	}

	rwe, err := parseInt(fields[cardRWE])
	if err != nil {
		return c, fmt.Errorf("invalid rms: %w", err)
	}
	rhe, err := parseInt(fields[cardRHE])
	if err != nil {
		return c, fmt.Errorf("invalid rps: %w", err)
	}
	rse, err := parseInt(fields[cardRSE])
	if err != nil {
		return c, fmt.Errorf("invalid rjs: %w", err)
	}

	return actionCard{
		ID:                 id,
		Advisor:            fields[cardAdvisor],
		Content:            fields[cardContent],
		AcceptText:         fields[cardAcceptText],
		AccWealthEffect:    awe,
		AccHealthEffect:    ahe,
		AccStabilityEffect: ase,
		RejectText:         fields[cardRejectText],
		RejWealthEffect:    rwe,
		RejHealthEffect:    rhe,
		RejStabilityEffect: rse,
	}, nil
}
