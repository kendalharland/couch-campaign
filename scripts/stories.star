load("src/stories/basic.star", "basic")

cards = {card.id: card for card in basic.cards}

stories = struct(basic=basic)
