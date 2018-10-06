from flask import (
    Blueprint,
    flash,
    redirect,
    render_template,
    request,
    url_for
)
from flask_login import login_user, logout_user
import re

from .shared import db
from .models import User

blueprint = Blueprint('base', __name__, template_folder='templates')


@blueprint.route('/')
def index():
    return render_template('index.jinja2')


@blueprint.route('/campaigns')
def campaigns():
    return 'VIEW: campaigns'


@blueprint.route('/search')
def search():
    return 'VIEW: search'


@blueprint.route('/glossary')
def glossary():
    return 'VIEW: glossary'


@blueprint.route('/help')
def help():
    return 'VIEW: help'


@blueprint.route('/profile/login', methods=['GET', 'POST'])
def profile_login():
    return 'VIEW: profile/login'


@blueprint.route('/profile/register', methods=['GET', 'POST'])
def profile_register():
    if request.method == 'POST':
        email = request.form['email']
        password = request.form['password']
        if not re.match(r'.+@(?:.+){2,}\.(?:.+){2,}', email):
            flash('Email does meet basic requirements', 'error')
            return redirect(url_for('.profile_register'))
        if len(password) < 5:
            flash('Password must be at least 5 characters long', 'error')
            return redirect(url_for('.profile_register'))
        new_user = User(email=email)
        new_user.set_password(password)
        db.session.add(new_user)
        db.session.commit()
        login_user(new_user)
        return redirect(url_for('.profile_settings'))
    return render_template('register.jinja2')


@blueprint.route('/profile/settings')
def profile_settings():
    return 'VIEW: profile/settings'


@blueprint.route('/profile/logout')
def profile_logout():
    logout_user()
    return redirect(url_for('.profile_login'))
