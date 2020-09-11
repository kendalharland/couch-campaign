load("//src/core.star", "core")
load("//engine.star", "engine")
load("//stories.star", "stories", "cards")

INPUT_ACCEPT_CARD = "accept"
INPUT_REJECT_CARD = "reject"


def new_game():
    for card in stories.basic.cards:
        core.deck.push(card.id)
    state = core.state()
    state.set_wealth(10)
    state.set_health(10)
    state.set_stability(10)
    state.set_character("you")
    _update_card(state)
    return struct()


def handle_input(game, input):
    if input == INPUT_ACCEPT_CARD:
        _on_card_accepted(game)
    elif input == INPUT_REJECT_CARD:
        _on_card_rejected(game)
    else:
        print("invalid input %s" % input)  # TODO: Add error fn.
        pass
    core.deck.pop()
    state = core.state()
    _update_card(state)
    state.set_character_lifespan(state.character_lifespan()+1)
    

def _on_card_shown(game):
    card = cards[core.state().card_ref()]
    for effect in card.on_show:
        engine.effect.apply(effect)


def _on_card_accepted(game):
    card = cards[core.state().card_ref()]
    for effect in card.on_accept:
        engine.effect.apply(effect)


def _on_card_rejected(game):
    card = cards[core.state().card_ref()]
    for effect in card.on_reject:
        engine.effect.apply(effect)

def _update_card(state):
    top_card = cards[core.deck.top()]
    state.set_card_ref(top_card.id)
    state.set_card_text(top_card.text)
    state.set_card_accept_text(top_card.accept_text)
    state.set_card_reject_text(top_card.reject_text)
    state.set_card_speaker(top_card.speaker)