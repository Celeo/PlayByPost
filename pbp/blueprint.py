from datetime import datetime
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
    current_user,
    login_user,
    login_required,
    logout_user
)

from .models import (
    Campaign,
    Character,
    Post,
    User
)
from .shared import db
from .util import is_safe_url

blueprint = Blueprint('base', __name__, template_folder='templates')


@blueprint.route('/')
def index():
    return render_template('index.jinja2')


@blueprint.route('/campaigns', methods=['GET', 'POST'])
def campaigns():
    # TODO when the user goes and makes the campaign, create for them a "DM"
    # user that's automatically tied to the campaign without having to join,
    # and set that character as the DM.
    if request.method == 'POST':
        new_campaign = Campaign(
            created_by_user_id=current_user.id,
            dm_user_id=current_user.id,
            name=request.form['name'],
            description=request.form['description'],
            date_created=datetime.utcnow()
        )
        db.session.add(new_campaign)
        db.session.commit()
        flash('New campaign created')
        return redirect(url_for('.campaigns'))
    campaigns = Campaign.query.all()
    return render_template('campaigns.jinja2', campaigns=campaigns)


@blueprint.route('/campaign/<int:campaign_id>/posts')
@blueprint.route('/campaign/<int:campaign_id>/posts/<int:page>')
def campaign_posts(campaign_id, page=1):
    # TODO pagination
    campaign = Campaign.query.get(campaign_id)
    if not campaign:
        flash('Could not find campaign with that id', 'error')
        return redirect(url_for('.campaigns'))
    posts = Post.query.filter_by(campaign_id=campaign_id).all()
    return render_template('campaign_posts.jinja2', campaign=campaign, posts=posts)


@blueprint.route('/campaign/<int:campaign_id>/info')
def campaign_info(campaign_id):
    campaign = Campaign.query.get(campaign_id)
    if not campaign:
        flash('Could not find campaign with that id', 'error')
        return redirect(url_for('.campaigns'))
    return render_template('campaign_info.jinja2', campaign=campaign)


@blueprint.route('/campaign/<int:campaign_id>/new_post', methods=['POST'])
def campaign_new_post(campaign_id):
    campaign = Campaign.query.get(campaign_id)
    if not campaign:
        flash('Could not find campaign with that id', 'error')
        return redirect(url_for('.campaigns'))
    character = current_user.get_character_in_campaign(campaign)
    if not character:
        flash('You are not a member of that campaign', 'error')
        return redirect(url_for('.campaign_posts', campaign_id=campaign_id))
    post = Post(
        character_id=character.id,
        campaign_id=campaign.id,
        date=datetime.utcnow(),
        tag=character.tag,
        content=request.form['content']
    )
    db.session.add(post)
    db.session.commit()
    flash('New post added')
    return redirect(url_for('.campaign_posts', campaign_id=campaign_id))


@blueprint.route('/campaign/<int:campaign_id>/join', methods=['GET', 'POST'])
def campaign_join(campaign_id):
    campaign = Campaign.query.get(campaign_id)
    if request.method == 'POST':
        character = Character.query.get(int(request.form['character']))
        if character.campaign_id:
            if not character.campaign_approved:
                flash('Your membership to that campaign is pending', 'error')
                return redirect(url_for('.campaign_join', campaign_id=campaign_id))
            flash('That character is already a member of that campaign')
            return redirect(url_for('.campaign_join', campaign_id=campaign_id))
        for other_character in campaign.characters:
            if other_character.name == character.name:
                flash('There is already a character in that campaign with that name', 'error')
                return redirect(url_for('.campaign_join', campaign_id=campaign_id))
        character.campaign_id = campaign_id
        db.session.commit()
        flash('Membership request submitted')
        return redirect(url_for('.campaign_join', campaign_id=campaign_id))
    return render_template('campaign_join.jinja2', campaign=campaign)


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
        new_user = User(email=email, date_joined=datetime.utcnow())
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
        form_field = request.form.get('field')
        new_value = request.form.get('value')
        character_id = request.form.get('character_id', 0, type=int)
        if form_field == 'new_character':
            character = Character(user_id=current_user.id, name=new_value)
            db.session.add(character)
            db.session.commit()
            flash('New character created')
            return redirect(url_for('.profile_characters'))
        elif form_field == 'delete':
            character = Character.query.get(character_id)
            if not character:
                flash('Unknown character', 'error')
                return redirect(url_for('.profile_characters'))
            if not character.user_id == current_user.id:
                flash('You are not the owner of that character', 'error')
                return redirect(url_for('.profile_characters'))
            db.session.delete(character)
            db.session.commit()
            flash('Character deleted')
            return redirect(url_for('.profile_characters'))
        else:
            character = Character.query.get(character_id)
            if not character:
                flash('Unknown character', 'error')
                return redirect(url_for('.profile_characters'))
            if not character.user_id == current_user.id:
                flash('You are not the owner of that character', 'error')
                return redirect(url_for('.profile_characters'))
            if form_field == 'name':
                if character.campaign_id:
                    for other_character in Character.query.filter_by(campaign_id=character.campaign_id):
                        if other_character.character.name == new_value:
                            flash('A character with that name is already in the same campaign', 'error')
                            return redirect(url_for('.profile_characters'))
                character.name = new_value
            elif form_field == 'tag':
                character.tag = new_value
            else:
                flash('An error occurred', 'error')
            db.session.commit()
            return redirect(url_for('.profile_characters'))
    characters = Character.query.filter_by(user_id=current_user.id).all()
    return render_template('profile_characters.jinja2', characters=characters)


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
