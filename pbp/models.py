from .shared import db


class User(db.Model):

    __tablename__ = 'users'

    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(100))
    email = db.Column(db.String(100))
    date_joined = db.Column(db.DateTime)


class Campaign(db.Model):

    __tablename__ = 'campaigns'

    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(100))
    created = db.Column(db.DateTime)
    locked = db.Column(db.Boolean, default=False)


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
