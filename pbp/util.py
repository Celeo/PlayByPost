from random import randint
import re

from flask import request
from urllib.parse import (
    urlparse,
    urljoin
)

from .models import Roll


regex_dice = re.compile(r'([+-]?\d+)d([+-]?\d+)')
regex_mod = re.compile(r'([+-])(\d+)')


class DiceException(Exception):
    pass


def is_safe_url(target):
    """Return true if the URL is save to navigate to.

    Args:
        target (str): url

    Returns:
        bool: true if safe
    """
    ref_url = urlparse(request.host_url)
    test_url = urlparse(urljoin(request.host_url, target))
    return test_url.scheme in ('http', 'https') and ref_url.netloc == test_url.netloc


def roll_dice(character, s):
    """Roll some dice.

    Args:
        character (pbp.models.Character): the character rolling
        s (str): the roll string from the client

    Returns:
        pbp.models.Roll: roll object for the dice (unsaved)
    """
    roll = Roll(character_id=character.id, string=s, is_crit=False)
    final_value = 0
    s = s.replace(' ', '')
    last_colon_index = s.rindex(':')
    dice_section = s[last_colon_index + 1:]
    dice = dice_section.split(',')
    for die in dice:
        die_result = 0
        groups = regex_dice.findall(die)
        count = int(groups[0][0])
        if count < 1:
            raise DiceException('Minimum roll count is 1')
        sides = int(groups[0][1])
        if sides < 2:
            raise DiceException('Minimum side count is 2')
        if sides > 100:
            raise DiceException('Maximum side count is 100')
        for i in range(count):
            roll_result = randint(1, sides)
            if sides == 20 and roll_result == 20:
                roll.is_crit = True
            die_result += roll_result
        groups = regex_mod.findall(die)
        if len(groups) > 0:
            delta = int(groups[0][1])
            if groups[0][0] == '+':
                die_result += delta
            else:
                die_result -= delta
        final_value += die_result
    roll.value = final_value
    return roll
