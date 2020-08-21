def _add_player_health(amount):
    return _builtins_add_player_health(amount)


def _add_player_wealth(amount):
    return _builtins_add_player_wealth(amount)


def _add_player_stability(amount):
    return _builtins_add_player_stability(amount)


def _get_player_card_id():
    return _builtins_get_player_card_id()


def _set_player_card_id(card_id):
    return _builtins_set_player_card_id(card_id)


# The core API provided by the interpreter.
core = struct(
    add_player_health=_add_player_health,
    add_player_wealth=_add_player_wealth,
    add_player_stability=_add_player_stability,
    get_player_card_id=_get_player_card_id,
    set_player_card_id=_set_player_card_id,
)
