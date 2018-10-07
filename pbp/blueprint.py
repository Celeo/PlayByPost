import re

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

from .models import (
    Campaign,
    Post,
    User
)
from .shared import db
from .util import is_safe_url

blueprint = Blueprint('base', __name__, template_folder='templates')


@blueprint.route('/')
def index():
    return render_template('index.jinja2')


@blueprint.route('/campaigns')
def campaigns():
    campaigns = Campaign.query.all()
    return render_template('campaigns.jinja2', campaigns=campaigns)


@blueprint.route('/campaign/<int:campaign_id>/posts')
def campaign_posts(campaign_id):
    # TODO pagination
    campaign = Campaign.query.get(campaign_id)
    if not campaign:
        flash('Could not find campaign with that id', 'error')
        return redirect(url_for('.campaigns'))
    posts = Post.query.filter_by(campaign_id=campaign_id).all()
    return render_template('campaign_posts.jinja2', campaign=campaign, posts=posts)


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
        next_url = request.args.get('next')
        if next_url and not is_safe_url(next_url):
            return redirect(url_for('.profile_settings'))
        return redirect(next_url or url_for('.profile_settings'))
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
