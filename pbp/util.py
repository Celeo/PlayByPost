from random import randint
import re
from uuid import uuid4

from celery import Celery
from flask import request
from urllib.parse import (
    urlparse,
    urljoin
)
import requests

from .models import Roll
from .shared import redis


regex_dice = re.compile(r'([+-]?\d+)d([+-]?\d+)')
regex_mod = re.compile(r'([+-])(\d+)')

celery = Celery('pbp', broker='redis://localhost:6379/0')


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


class DiceException(Exception):
    pass


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


def is_valid_email(email):
    return re.match(r'.+@(?:.+){2,}\.(?:.+){2,}', email)


def pagination_pages(current_page, page_count):
    return [
        page for page in
        [current_page - 2, current_page - 1, current_page, current_page + 1, current_page + 2]
        if 0 < page <= page_count
    ]


def create_password_reset_key(email):
    key = str(uuid4())
    redis.set(f'password_reset:{email}', key, ex=1 * 60 * 60)
    return key


def get_password_reset_key(email):
    key = redis.get(f'password_reset:{email}')
    return key.decode('UTF-8') if key else None


def clear_password_reset_keys(email):
    redis.delete(f'password_reset:{email}')


@celery.task()
def send_email(api_key, domain, from_, recipients, subject, body):
    return requests.post(
        'https://api.mailgun.net/v3/{}/messages'.format(domain),
        auth=('api', api_key),
        data={
            'from': from_,
            'to': recipients,
            'subject': subject,
            'text': body
        }
    )
