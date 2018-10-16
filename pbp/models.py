from datetime import (
    datetime,
    timedelta
)

import bcrypt

from .shared import db


class User(db.Model):
    """A site user.

    This represents a single real person using the site.

    A user can have multiple characters.
    """

    __tablename__ = 'users'

    id = db.Column(db.Integer, primary_key=True)
    email = db.Column(db.String(200), unique=True)
    password = db.Column(db.String(500))
    date_joined = db.Column(db.DateTime)
    is_banned = db.Column(db.Boolean, default=False)
    is_admin = db.Column(db.Boolean, default=False)
    posts_per_page = db.Column(db.Integer, default=20, nullable=False)
    posts_newest_first = db.Column(db.Boolean, default=True, nullable=False)

    @property
    def is_authenticated(self):
        return True

    @property
    def is_anonymous(self):
        return False

    @property
    def is_active(self):
        return not self.is_banned

    def get_id(self):
        return str(self.id)

    def set_password(self, string):
        self.password = bcrypt.hashpw(string.encode('utf8'), bcrypt.gensalt())

    def check_password(self, string):
        return bcrypt.checkpw(string.encode('utf8'), self.password)

    def get_character_in_campaign(self, campaign):
        for character in self.characters:
            if character.campaign_id == campaign.id and character.campaign_approved:
                return character
        return None

    def get_character_applied_to_campaign(self, campaign):
        for character in self.characters:
            if character.campaign_id == campaign.id and not character.campaign_approved:
                return character
        return None

    def should_show_join_link(self, campaign):
        for character in self.characters:
            if character.campaign_id == campaign.id:
                return False
        return True

    def is_dm_to_campaign(self, campaign):
        return campaign.dm_character.user_id == self.id


class Campaign(db.Model):
    """A collection of posts.

    A campaign is a series of posts both from a DM and other characters
    telling a story.

    A campaign can have multiple posts and multiple characters.
    """

    __tablename__ = 'campaigns'

    id = db.Column(db.Integer, primary_key=True)
    creator_user_id = db.Column(db.Integer, db.ForeignKey('users.id'), nullable=False)
    dm_character_id = db.Column(db.Integer, db.ForeignKey('characters.id'))
    name = db.Column(db.String(100), nullable=False)
    description = db.Column(db.String(5000))
    date_created = db.Column(db.DateTime)
    is_locked = db.Column(db.Boolean, default=False)
    is_posts_public = db.Column(db.Boolean, default=True)

    created_by_user = db.relationship('User', foreign_keys=[creator_user_id], backref=db.backref('created_campaigns', lazy=True))
    dm_character = db.relationship('Character', uselist=False, foreign_keys=[dm_character_id], backref=db.backref('dm_campaign', lazy=True))

    @property
    def dm_posts(self):
        return [post for post in self.posts if post.character.user_id == self.dm_character_id]


class Character(db.Model):
    """A fictional persona.

    A character is a persona that a person "steps into" in order
    to participate in a story.

    A character can have multiple posts and multiple rolls, and
    is a member of either no or one campaign.
    """

    __tablename__ = 'characters'

    id = db.Column(db.Integer, primary_key=True)
    user_id = db.Column(db.Integer, db.ForeignKey('users.id'), nullable=False)
    campaign_id = db.Column(db.Integer, db.ForeignKey('campaigns.id'), nullable=True)
    name = db.Column(db.String(100))
    tag = db.Column(db.String(100))
    campaign_approved = db.Column(db.Boolean, default=False)

    user = db.relationship('User', foreign_keys=[user_id], backref=db.backref('characters', lazy=True))
    campaign = db.relationship('Campaign', foreign_keys=[campaign_id], backref=db.backref('characters', lazy=True))


class Post(db.Model):
    """An article.

    A post is a collection of text and rolls that represents a character
    taking action in the story.

    A post can have multiple rolls and is in a single campaign by a single character.
    """

    __tablename__ = 'posts'

    id = db.Column(db.Integer, primary_key=True)
    character_id = db.Column(db.Integer, db.ForeignKey('characters.id'), nullable=False)
    campaign_id = db.Column(db.Integer, db.ForeignKey('campaigns.id'), nullable=False)
    date = db.Column(db.DateTime)
    tag = db.Column(db.String(100))
    content = db.Column(db.String(10000))

    character = db.relationship('Character', backref=db.backref('posts', lazy=True))
    campaign = db.relationship('Campaign', backref=db.backref('posts', lazy=True))

    @property
    def user(self):
        return self.character.user

    @property
    def can_be_edited(self):
        return (datetime.utcnow() - self.date) <= timedelta(minutes=30)


class Roll(db.Model):
    """A random chance.

    A roll is the random representation of resolving success or failure.

    A roll is part of a single post and is made by a single character.
    """

    __tablename__ = 'rolls'

    id = db.Column(db.Integer, primary_key=True)
    character_id = db.Column(db.Integer, db.ForeignKey('characters.id'), nullable=False)
    post_id = db.Column(db.Integer, db.ForeignKey('posts.id'))
    string = db.Column(db.String(100))
    value = db.Column(db.Integer, default=0)
    is_crit = db.Column(db.Boolean, default=False)

    post = db.relationship('Post', backref=db.backref('rolls', lazy=True))
    character = db.relationship('Character', backref=db.backref('rolls', lazy=True))

    def to_dict(self):
        return {
            'id': self.id,
            'character_id': self.character_id,
            'post_id': self.post_id,
            'pending': self.pending,
            'string': self.string,
            'value': self.value,
            'is_crit': self.is_crit
        }

    @property
    def pending(self):
        return self.post_id is None
