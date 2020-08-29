def _get_player():
    return _builtins_get_player()


def _deck_pop_card():
    return _builtins_deck_pop_card()


def _deck_push_card(card_ref):
    return _builtins_deck_push_card(card_ref)


# The core API provided by the interpreter.
core = struct(
    deck=struct(pop=_deck_pop_card, push=_deck_push_card,),
    player=_get_player,
)
