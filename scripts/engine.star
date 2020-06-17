load("//src/engine/api.star", "api")
load("//src/engine/card.star", "card")
load("//src/engine/effect.star", "effect")
load("//src/engine/game.star", "game")
load("//src/engine/story.star", "story")

engine = struct(api=api, card=card, effect=effect, game=game, story=story,)
