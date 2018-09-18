from django.db import models


class Campaign(models.Model):

    name = models.CharField()
    tag = models.CharField()
    description = models.TextField()
    created = models.DateTimeField()
    locked = models.BooleanField(default=False)


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
    tag = models.CharField()
    content = models.TextField()


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
    string = models.CharField()
    value = models.IntegerField()
    is_crit = models.BooleanField()


class Glossary(models.Model):

    campaign = models.ForeignKey(
        'Campaign',
        on_delete=models.CASCADE
    )
    content = models.TextField()
    updated = models.DateTimeField()
