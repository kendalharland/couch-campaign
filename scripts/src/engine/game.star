def _new():
    return struct(player_ids=[])


def _add_player(game, id):
    game.player_ids.append(id)


game = struct(new=_new, add_player=_add_player,)
