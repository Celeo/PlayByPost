from datetime import datetime

from flask import (
    Blueprint,
    flash,
    jsonify,
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
    User,
    Roll
)
from .shared import db
from .util import (
    is_safe_url,
    is_valid_email,
    pagination_pages,
    roll_dice
)

blueprint = Blueprint('base', __name__, template_folder='templates')


@blueprint.route('/')
def index():
    return render_template('index.jinja2')


@blueprint.route('/campaigns')
def campaigns():
    campaigns = Campaign.query.all()
    return render_template('campaigns.jinja2', campaigns=campaigns)


@blueprint.route('/campaigns/create', methods=['GET', 'POST'])
@login_required
def campaign_create():
    if request.method == 'POST':
        new_campaign = Campaign(
            creator_user_id=current_user.id,
            name=request.form['name'],
            description=request.form['description'],
            date_created=datetime.utcnow()
        )
        new_dm = Character(
            user_id=current_user.id,
            name='DM',
            tag='Dungeon Master',
            campaign_approved=True,
        )
        new_campaign.dm_character = new_dm
        db.session.add(new_campaign)
        db.session.add(new_dm)
        db.session.commit()
        new_dm.campaign_id = new_campaign.id
        db.session.commit()
        flash('New campaign created')
        return redirect(url_for('.campaigns'))
    campaigns = Campaign.query.all()
    return render_template('campaigns_create.jinja2', campaigns=campaigns)


@blueprint.route('/campaign/<int:campaign_id>/posts')
@blueprint.route('/campaign/<int:campaign_id>/posts/<int:page>')
def campaign_posts(campaign_id, page=1):
    campaign = Campaign.query.get(campaign_id)
    if not campaign:
        flash('Could not find campaign with that id', 'error')
        return redirect(url_for('.campaigns'))
    query = Post.query.filter_by(campaign_id=campaign_id)
    if current_user.posts_newest_first:
        query = query.order_by(Post.id.desc())
    pagination = query.paginate(page=page, per_page=current_user.posts_per_page)
    return render_template(
        'campaign_posts.jinja2',
        campaign=campaign,
        posts=pagination.items,
        pages=pagination.pages,
        page=page,
        pagination_pages=pagination_pages
    )


@blueprint.route('/campaign/<int:campaign_id>/info')
def campaign_info(campaign_id):
    campaign = Campaign.query.get(campaign_id)
    if not campaign:
        flash('Could not find campaign with that id', 'error')
        return redirect(url_for('.campaigns'))
    return render_template('campaign_info.jinja2', campaign=campaign)


@blueprint.route('/campaign/<int:campaign_id>/new_post', methods=['POST'])
@login_required
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
    pending_rolls = Roll.query.filter_by(character_id=character.id, post_id=None).all()
    for roll in pending_rolls:
        roll.post = post
    db.session.commit()
    flash('New post added')
    return redirect(url_for('.campaign_posts', campaign_id=campaign_id))


@blueprint.route('/campaign/<int:campaign_id>/roll', methods=['GET', 'POST'])
@login_required
def campaign_rolls(campaign_id):
    campaign = Campaign.query.get(campaign_id)
    if not campaign:
        return 'Could not find campaign with that id', 404
    character = current_user.get_character_in_campaign(campaign)
    if not character:
        return 'You are not a member of that campaign', 403
    if request.method == 'POST':
        roll = request.json.get('roll')
        if not roll:
            return '', 400
        roll = roll_dice(character, roll)
        db.session.add(roll)
        db.session.commit()
    rolls = Roll.query.filter_by(character_id=current_user.get_character_in_campaign(campaign).id, post_id=None).all()
    return jsonify([r.to_dict() for r in rolls])


@blueprint.route('/post/<int:post_id>/edit', methods=['GET', 'POST'])
@login_required
def campaign_edit_post(post_id):
    post = Post.query.get(post_id)
    if not post:
        flash('Could not find a post with that id', 'error')
        return redirect(url_for('.campaigns'))
    if not post.character.user_id == current_user.id:
        flash('That isn\'t your post', 'error')
        return redirect(url_for('.campaign_posts', campaign_id=post.campaign_id))
    if not post.can_be_edited:
        flash('This post can no longer be edited', 'error')
        return redirect(url_for('.campaign_posts', campaign_id=post.campaign_id))
    if request.method == 'POST':
        content = request.form['content']
        post.content = content
        db.session.commit()
        flash('Content saved')
        return redirect(url_for('.campaign_posts', campaign_id=post.campaign_id))
    return render_template('campaign_edit_post.jinja2', post=post)


@blueprint.route('/campaign/<int:campaign_id>/join', methods=['GET', 'POST'])
@login_required
def campaign_join(campaign_id):
    campaign = Campaign.query.get(campaign_id)
    if not current_user.should_show_join_link(campaign):
        return redirect(url_for('.campaign_posts', campaign_id=campaign_id))
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


@blueprint.route('/campaign/<int:campaign_id>/dm_controls', methods=['GET', 'POST'])
@login_required
def campaign_dm_controls(campaign_id):
    campaign = Campaign.query.get(campaign_id)
    if not campaign:
        flash('Could not find campaign with that id', 'error')
        return redirect(url_for('.campaigns'))
    if not current_user.is_dm_to_campaign(campaign):
        flash('You are not a DM of that campaign', 'error')
        return redirect(url_for('.campaign_posts', campaign_id=campaign_id))
    if request.method == 'POST':
        form_type = request.form['type']
        if form_type == 'applicant':
            character = Character.query.get(request.form['character_id'])
            if not character:
                flash('Unknown character', 'error')
                return redirect(url_for('.campaign_dm_controls', campaign_id=campaign_id))
            if not character.campaign_id == campaign_id:
                flash('That character has not applied for this campaign', 'error')
                return redirect(url_for('.campaign_dm_controls', campaign_id=campaign_id))
            if character.campaign_approved:
                flash('That character is already approved for this campaign', 'error')
                return redirect(url_for('.campaign_dm_controls', campaign_id=campaign_id))
            if request.form['action'] == 'accept':
                character.campaign_approved = True
                flash('Character accepted')
            else:
                db.session.delete(character)
                flash('Character denied')
            db.session.commit()
            return redirect(url_for('.campaign_dm_controls', campaign_id=campaign_id))
        elif form_type == 'name_description':
            campaign.name = request.form['name']
            campaign.description = request.form['description']
            db.session.commit()
            flash('Campaign name/desciption updated')
            return redirect(url_for('.campaign_dm_controls', campaign_id=campaign_id))
        else:
            flash('Unknown form submission', 'error')
            return redirect(url_for('.campaign_dm_controls', campaign_id=campaign_id))
    applicants = Character.query.filter_by(campaign_id=campaign_id, campaign_approved=False).all()
    members = Character.query.filter_by(campaign_id=campaign_id, campaign_approved=True).all()
    return render_template('campaign_dm_controls.jinja2', campaign=campaign, applicants=applicants, members=members)


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
            return redirect(url_for('.campaigns'))
        return redirect(next_url or url_for('.campaigns'))
    return render_template('login.jinja2')


@blueprint.route('/profile/register', methods=['GET', 'POST'])
def profile_register():
    if request.method == 'POST':
        email = request.form['email']
        if User.query.filter_by(email=email).first():
            flash('Email already in use', 'error')
            return redirect(url_for('.profile_register'))
        password = request.form['password']
        if not is_valid_email(email):
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
        return redirect(url_for('.campaigns'))
    return render_template('register.jinja2')


@blueprint.route('/profile/characters', methods=['GET', 'POST'])
@login_required
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
            if character.campaign_approved:
                flash('You cannot delete a character that\'s part of a campaign', 'error')
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
                if character.name == 'DM':
                    flash('You cannot rename a DM character', 'error')
                    return redirect(url_for('.profile_characters'))
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
        settings_type = request.form['settings_type']
        if settings_type == 'posts':
            current_user.posts_per_page = request.form['posts_per_page']
            current_user.posts_newest_first = request.form['posts_newest_first'] == 'newest'
            db.session.commit()
            flash('Settings saved')
            return redirect(url_for('base.profile_settings'))
        elif settings_type == 'email':
            new_email = request.form['email']
            if current_user.email == new_email:
                flash('That\'s already your email', 'error')
                return redirect(url_for('base.profile_settings'))
            if User.query.filter_by(email=new_email).first():
                flash('Email already in use', 'error')
                return redirect(url_for('base.profile_settings'))
            if is_valid_email(new_email):
                current_user.email = new_email
                db.session.commit()
                flash('Settings saved')
            else:
                flash('Email does meet basic requirements', 'error')
            return redirect(url_for('base.profile_settings'))
        elif settings_type == 'password':
            if not current_user.check_password(request.form['old_password']):
                flash('Incorrect current password', 'error')
                return redirect(url_for('base.profile_settings'))
            if not request.form['new_password'] == request.form['new_password_confirm']:
                flash('New passwords don\'t match', 'error')
                return redirect(url_for('base.profile_settings'))
            if not len(request.form['new_password']) > 5:
                flash('Password must be at least 5 characters long', 'error')
                return redirect(url_for('base.profile_settings'))
            current_user.set_password(request.form['new_password'])
            db.session.commit()
            flash('New password saved')
            return redirect(url_for('base.profile_settings'))
        else:
            flash('Unknown setting value', 'error')
            return redirect(url_for('base.profile_settings'))
    return render_template('profile_settings.jinja2')


@blueprint.route('/profile/logout')
def profile_logout():
    logout_user()
    return redirect(url_for('.profile_login'))
