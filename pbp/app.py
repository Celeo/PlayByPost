from flask import Flask
from flask_login import LoginManager

from .shared import db
from .models import User
from .blueprint import blueprint as base_blueprint


app = Flask(__name__)
app.config.from_json('config.json')

db.app = app
db.init_app(app)

login_manager = LoginManager()
login_manager.init_app(app)


@login_manager.user_loader
def load_user(user_id):
    return User.get(user_id)


app.register_blueprint(base_blueprint)
