import os

import arrow
from flask import Flask, render_template
from flask_login import LoginManager

from .shared import db, redis
from .models import User
from .blueprint import blueprint as base_blueprint


app = Flask(__name__)
app.config.from_json('config.json')

db.app = app
db.init_app(app)

redis.init_app(app)

login_manager = LoginManager()
login_manager.init_app(app)
login_manager.login_view = 'base.profile_login'
login_manager.login_message_category = 'error'

app.jinja_env.globals['FLASK_ENV'] = os.getenv('FLASK_ENV', 'production')


@login_manager.user_loader
def load_user(user_id):
    return User.query.get(int(user_id))


@app.template_filter('format_date')
def filter_format_date(dt):
    return arrow.get(dt).to('America/Los_Angeles').strftime('%b %m, %y @ %I:%M:%S %p')


app.register_blueprint(base_blueprint)


@app.errorhandler(404)
def error_404(error):
    print(f'404 ERROR: {error}')
    return render_template('error_404.jinja2')


@app.errorhandler(400)
def error_400(error):
    print(f'400 ERROR: {error}')
    return render_template('error_400.jinja2')


@app.errorhandler(500)
def error_500(error):
    print(f'500 ERROR: {error}')
    return render_template('error_500.jinja2')
