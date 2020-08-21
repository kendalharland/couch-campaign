load("//src/core.star", "core")
load("//engine.star", "engine")
load("//stories.star", "stories", "cards")

INPUT_SHOW_CARD = "show"
INPUT_ACCEPT_CARD = "accept"
INPUT_REJECT_CARD = "reject"


def new_game():
    game = engine.game.new()
    # TODO: Initialize the player's deck instead, and add the id of the top card.
    core.set_player_card_id("viral_infection")

    return game


def handle_input(game, input):
    if input == INPUT_SHOW_CARD:
        _on_card_shown(game)
    elif input == INPUT_ACCEPT_CARD:
        _on_card_accepted(game)
    elif input == INPUT_REJECT_CARD:
        _on_card_rejected(game)
    else:
        print("invalid input %s" % input)  # TODO: Add error fn.
        pass


def _on_card_shown(game):
    card = cards[core.get_player_card_id()]
    for effect in card.on_show:
        engine.effect.apply(effect)


def _on_card_accepted(game):
    card = cards[core.get_player_card_id()]
    for effect in card.on_accept:
        engine.effect.apply(effect)


def _on_card_rejected(game):
    card = cards[core.get_player_card_id()]
    for effect in card.on_reject:
        engine.effect.apply(effect)
