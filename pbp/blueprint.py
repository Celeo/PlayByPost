from flask import (
    Blueprint,
    flash,
    redirect,
    render_template,
    request,
    url_for
)
from flask_login import (
    login_user,
    login_required,
    logout_user
)
import re

from .shared import db
from .models import User

blueprint = Blueprint('base', __name__, template_folder='templates')


@blueprint.route('/')
def index():
    return render_template('index.jinja2')


@blueprint.route('/campaigns')
def campaigns():
    return render_template('campaigns.jinja2')


@blueprint.route('/search')
def search():
    return render_template('search.jinja2')


@blueprint.route('/glossary')
def glossary():
    return render_template('glossary.jinja2')


@blueprint.route('/help')
def help():
    return render_template('help.jinja2')


@blueprint.route('/profile/login', methods=['GET', 'POST'])
def profile_login():
    if request.method == 'POST':
        email = request.form['email']
        password = request.form['password']
        user = User.query.filter_by(email=email).first()
        if not user or not user.check_password(password):
            flash('Login failed', 'error')
            return redirect(url_for('.profile_login'))
        flash('Login successful')
        login_user(user, remember=True)
        # TODO https://flask-login.readthedocs.io/en/latest/#login-example
        return redirect(url_for('.profile_settings'))
    return render_template('login.jinja2')


@blueprint.route('/profile/register', methods=['GET', 'POST'])
def profile_register():
    if request.method == 'POST':
        email = request.form['email']
        if User.query.filter_by(email=email).first():
            flash('Email already in use', 'error')
            return redirect(url_for('.profile_register'))
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
        flash('Login successful')
        login_user(new_user, remember=True)
        return redirect(url_for('.profile_settings'))
    return render_template('register.jinja2')


@blueprint.route('/profile/characters', methods=['GET', 'POST'])
def profile_characters():
    if request.method == 'POST':
        return 'TODO'
    return render_template('profile_characters.jinja2')


@blueprint.route('/profile/settings', methods=['GET', 'POST'])
@login_required
def profile_settings():
    if request.method == 'POST':
        return 'TODO'
    return render_template('profile_settings.jinja2')


@blueprint.route('/profile/logout')
def profile_logout():
    logout_user()
    return redirect(url_for('.profile_login'))
