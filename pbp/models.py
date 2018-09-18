from django.db import models
from django.contrib.auth.models import AbstractBaseUser, PermissionsMixin
from django.contrib.auth.validators import UnicodeUsernameValidator


class User(AbstractBaseUser, PermissionsMixin):

    username_validator = UnicodeUsernameValidator()

    name = models.CharField(
        max_length=150,
        blank=False,
        unique=True,
        validators=[username_validator]
    )
    email = models.EmailField(blank=True)
    is_staff = models.BooleanField(default=False)
    is_active = models.BooleanField(default=True)
    is_superuser = models.BooleanField(default=False)
    date_joined = models.DateTimeField()

    USERNAME_FIELD = 'name'
    EMAIL_FIELD = 'email'
    REQUIRED_FIELDS = []


# TODO user settings model, with one-to-one with User
# contains stuff like posts_per_page, newest_at_top, tag, if they want email notifications, etc.


class Campaign(models.Model):

    name = models.CharField(max_length=200)
    description = models.TextField()
    created = models.DateTimeField()
    locked = models.BooleanField(default=False)
    # TODO users? one to many -> 'User

    def __str__(self):
        return f'<Campaign {self.id}>'


class Post(models.Model):

    user = models.ForeignKey(
        'User',
        on_delete=models.CASCADE
    )
    campaign = models.ForeignKey(
        'Campaign',
        on_delete=models.CASCADE
    )
    date = models.DateTimeField()
    tag = models.CharField(max_length=200)
    content = models.TextField()

    def __str__(self):
        return f'<Post {self.id}>'


class Roll(models.Model):

    user = models.ForeignKey(
        'User',
        on_delete=models.CASCADE
    )
    post = models.ForeignKey(
        'Post',
        on_delete=models.CASCADE
    )
    pending = models.BooleanField()
    string = models.CharField(max_length=200)
    value = models.IntegerField()
    is_crit = models.BooleanField()

    def __str__(self):
        return f'<Roll {self.id}>'


class Glossary(models.Model):

    campaign = models.ForeignKey(
        'Campaign',
        on_delete=models.CASCADE
    )
    content = models.TextField()
    updated = models.DateTimeField()

    def __str__(self):
        return f'<Glossary {self.id}>'
