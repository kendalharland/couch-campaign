def _get_state():
    return _builtins_get_state()


def _deck_pop_card():
    return _builtins_deck_pop_card()


def _deck_push_card(card_ref):
    return _builtins_deck_push_card(card_ref)


def _deck_top_card():
    return _builtins_deck_top_card()


# The core API provided by the interpreter.
core = struct(
    deck=struct(pop=_deck_pop_card, push=_deck_push_card, top=_deck_top_card),
    state=_get_state,
)
