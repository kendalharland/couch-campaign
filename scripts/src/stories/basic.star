load("//engine.star", "engine")

_cards = [
    engine.card.new(
        id="viral_infection",
        speaker="SurgeonGeneral",
        text="A novel viral infection broke out among a cargo ship waiting to come ashore. What should we do?",
        accept_text="Quarantine the ship in the harbor.",
        reject_text="Have them deliver their goods, then see them to the nearest hospital immediately.",
        on_accept=[engine.effect.add_health(-2), engine.effect.add_wealth(2),],
        on_reject=[engine.effect.add_health(-3), engine.effect.add_wealth(-2),],
    ),
    engine.card.new(
        id="tobbacco_ad",
        speaker="SurgeonGeneral",
        text="A tobacco company wants to advertise at the community center. They're offering to cut us in on the profits...",
        accept_text="Our coffers have been running a little dry...",
        reject_text="I won't sacrifice the public health for financial gain.",
        on_accept=[engine.effect.add_health(-2), engine.effect.add_wealth(2),],
        on_reject=[engine.effect.add_health(-2), engine.effect.add_wealth(2),],
    ),
]

basic = engine.story.new(id="basic", cards=_cards)
