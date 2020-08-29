load("//src/core.star", "core")

ADD_HEALTH = "ADD_HEALTH"
ADD_WEALTH = "ADD_WEALTH"
ADD_STABILITY = "ADD_STABILITY"


def _add_health(amount=0):
    return struct(id=ADD_HEALTH, amount=amount)


def _add_wealth(amount=0):
    return struct(id=ADD_WEALTH, amount=amount)


def _add_stability(amount):
    return struct(id=ADD_STABILITY, amount=amount)


def _apply(effect):
    if effect.id == ADD_HEALTH:
        core.state().set_health(core.state().health() + effect.amount)
    elif effect.id == ADD_WEALTH:
        core.state().set_wealth(core.state().wealth() + effect.amount)
    elif effect.id == ADD_STABILITY:
        core.state().set_stability(core.state().stability() + effect.amount)
    else:
        # TODO: error
        pass


effect = struct(
    apply=_apply,
    add_health=_add_health,
    add_wealth=_add_wealth,
    add_stability=_add_stability,
)
