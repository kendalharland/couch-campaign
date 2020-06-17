load("engine.star", "engine")
load("stories.star", "stories", "cards")

INPUT_SHOW_CARD = "show"
INPUT_ACCEPT_CARD = "accept"
INPUT_REJECT_CARD = "reject"


def new_game():
    game = engine.game.new()
    for id in engine.api.get_player_ids():
        engine.game.add_player(game, id)
        # TODO: Initialize the player's deck instead, and add the id of the top card.
        engine.api.set_player_card_id(id, "viral_infection")

    return game


def handle_input(game, player_id, input):
    if input == INPUT_SHOW_CARD:
        _on_card_shown(game, player_id)
    elif input == INPUT_ACCEPT_CARD:
        _on_card_accepted(game, player_id)
    elif input == INPUT_REJECT_CARD:
        _on_card_rejected(game, player_id)
    else:
        print("invalid input %s" % input)  # TODO: Add error fn.
        pass


def _on_card_shown(game, player_id):
    card = cards[engine.api.get_player_card_id(player_id)]
    for effect in card.on_show:
        engine.effect.apply(player_id, effect)


def _on_card_accepted(game, player_id):
    card = cards[engine.api.get_player_card_id(player_id)]
    for effect in card.on_accept:
        engine.effect.apply(player_id, effect)


def _on_card_rejected(game, player_id):
    card = cards[engine.api.get_player_card_id(player_id)]
    for effect in card.on_reject:
        engine.effect.apply(player_id, effect)
