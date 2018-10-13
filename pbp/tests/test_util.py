import pytest

from ..models import Character
from ..util import roll_dice


@pytest.fixture
def character():
    return Character(id=5)


def test_roll_dice_simple(character):
    s = 'Perception: 1d20'
    r = roll_dice(character, s)
    assert r.character_id == 5
    assert r.string == s
    assert 1 <= r.value <= 20
    assert r.is_crit is not None


def test_roll_dice_with_mod(character):
    s = 'Perception: 1d20 + 3'
    r = roll_dice(character, s)
    assert r.character_id == 5
    assert r.string == s
    assert 4 <= r.value <= 23
    assert r.is_crit is not None


def test_roll_weird_string(character):
    s = 'Foo: bar: $%sdf24d@#&++: 2d6 - 3'
    r = roll_dice(character, s)
    assert r.character_id == 5
    assert r.string == s
    assert -1 <= r.value <= 9
    assert r.is_crit is False
