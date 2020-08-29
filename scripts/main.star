load("//src/core.star", "core")
load("//engine.star", "engine")
load("//stories.star", "stories", "cards")

INPUT_ACCEPT_CARD = "accept"
INPUT_REJECT_CARD = "reject"


def new_game():
    for card in stories.basic.cards:
        core.deck.push(card.id)
    game = engine.game.new()
    core.player().set_wealth(10)
    core.player().set_health(10)
    core.player().set_stability(10)
    print(core.player().wealth())
    return game


def handle_input(game, input):
    if input == INPUT_ACCEPT_CARD:
        _on_card_accepted(game)
    elif input == INPUT_REJECT_CARD:
        _on_card_rejected(game)
    else:
        print("invalid input %s" % input)  # TODO: Add error fn.
        pass
    core.deck.pop()
    

def _on_card_shown(game):
    card = cards[core.player().card_ref()]
    for effect in card.on_show:
        engine.effect.apply(effect)

def _on_card_accepted(game):
    card = cards[core.player().card_ref()]
    for effect in card.on_accept:
        engine.effect.apply(effect)


def _on_card_rejected(game):
    card = cards[core.player().card_ref()]
    for effect in card.on_reject:
        engine.effect.apply(effect)
