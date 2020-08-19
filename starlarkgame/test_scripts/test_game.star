def new_game():
    return struct()


def handle_input(game, player_id, input):
    # Errs if the player does not exist
    _builtins_get_player_card_id(player_id)
    return game
