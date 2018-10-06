import bcrypt

from .shared import db


# Will need membership of characters to campaign
# Will need ability to track requests and (dis)approvals for those requests


class User(db.Model):

    __tablename__ = 'users'

    id = db.Column(db.Integer, primary_key=True)
    email = db.Column(db.String(200), unique=True)
    password = db.Column(db.String)
    date_joined = db.Column(db.DateTime)
    is_banned = db.Column(db.Boolean, default=False)
    is_admin = db.Column(db.Boolean, default=False)

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


class Character(db.Model):

    __tablename__ = 'characters'

    id = db.Column(db.Integer, primary_key=True)
    user_id = db.Column(db.Integer, db.ForeignKey('users.id'), nullable=False)
    name = db.Column(db.String(100))

    user = db.relationship('User', backref=db.backref('characters', lazy=True))


class Campaign(db.Model):

    __tablename__ = 'campaigns'

    id = db.Column(db.Integer, primary_key=True)
    created_by_user_id = db.Column(db.Integer, db.ForeignKey('users.id'), nullable=False)
    dm_user_id = db.Column(db.Integer, db.ForeignKey('users.id'), nullable=False)
    name = db.Column(db.String(100))
    description = db.Column(db.String)
    date_created = db.Column(db.DateTime)
    is_locked = db.Column(db.Boolean, default=False)
    is_posts_public = db.Column(db.Boolean, default=True)

    created_by_user = db.relationship('User', foreign_keys=[created_by_user_id], backref=db.backref('created_campaigns', lazy=True))
    dm_user = db.relationship('User', foreign_keys=[dm_user_id], backref=db.backref('dm_campaigns', lazy=True))


class Post(db.Model):

    __tablename__ = 'posts'

    id = db.Column(db.Integer, primary_key=True)
    user_id = db.Column(db.Integer, db.ForeignKey('users.id'), nullable=False)
    campaign_id = db.Column(db.Integer, db.ForeignKey('campaigns.id'), nullable=False)
    date = db.Column(db.DateTime)
    tag = db.Column(db.String(100))
    content = db.Column(db.String)

    user = db.relationship('User', backref=db.backref('posts', lazy=True))
    campaign = db.relationship('Campaign', backref=db.backref('posts', lazy=True))


class Roll(db.Model):

    __tablename__ = 'rolls'

    id = db.Column(db.Integer, primary_key=True)
    user_id = db.Column(db.Integer, db.ForeignKey('users.id'), nullable=False)
    post_id = db.Column(db.Integer, db.ForeignKey('posts.id'), nullable=False)
    pending = db.Column(db.Boolean, default=True)
    string = db.Column(db.String(100))
    value = db.Column(db.Integer, default=0)
    is_crit = db.Column(db.Boolean, default=False)

    post = db.relationship('Post', backref=db.backref('rolls', lazy=True))
    user = db.relationship('User', backref=db.backref('rolls', lazy=True))
