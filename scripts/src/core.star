def _add_player_health(id, amount):
    return _builtins_add_player_health(id, amount)


def _add_player_wealth(id, amount):
    return _builtins_add_player_health(id, amount)


def _add_player_stability(id, amount):
    return _builtins_add_player_health(id, amount)


def _get_player_ids():
    return _builtins_get_player_ids()


def _get_player_card_id(id):
    return _builtins_get_player_card_id(id)


def _set_player_card_id(player_id, card_id):
    return _builtins_set_player_card_id(player_id, card_id)


core = struct(
    add_player_health=_add_player_health,
    add_player_wealth=_add_player_wealth,
    add_player_stability=_add_player_stability,
    get_player_ids=_get_player_ids,
    get_player_card_id=_get_player_card_id,
    set_player_card_id=_set_player_card_id,
)
