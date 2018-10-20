from datetime import datetime

from flask import (
    Blueprint,
    current_app,
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
from .shared import (
    csrf,
    db
)
from .util import (
    create_password_reset_key,
    clear_password_reset_keys,
    get_password_reset_key,
    is_safe_url,
    is_valid_email,
    pagination_pages,
    roll_dice,
    send_email as _send_email
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
        # TODO check if name is unique
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
        current_app.logger.info(f'User {current_user.id} created new campaign with name "{new_campaign.name}"')
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
    if current_user.is_authenticated and current_user.posts_newest_first:
        query = query.order_by(Post.id.desc())
    pagination = query.paginate(page=page, per_page=current_user.posts_per_page if current_user.is_authenticated else 20)
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
    current_app.logger.info(f'User {current_user.id} made new post in campaign {campaign.id}')
    flash('New post added')
    is_dm_post = character.campaign.dm_character_id == character.id
    for other_character in campaign.characters:
        if other_character.id == character.id:
            continue
        link = url_for('.campaign_posts', campaign_id=campaign.id, _external=True)
        if is_dm_post and other_character.user.email_for_dm_post:
            current_app.logger.info(f'Send DM post notify email to {other_character.user.id} for campaign {campaign.id}')
            send_email(
                [other_character.user.email],
                f'New DM post in "{campaign.name}"',
                f'The DM has made a new post in the campaign.\n\nCampaign link: {link}'
            )
        elif not is_dm_post and other_character.user.email_for_any_post:
            current_app.logger.info(f'Send generic post notify email to {other_character.user.id} for campaign {campaign.id}')
            send_email(
                [other_character.user.email],
                f'New post in "{campaign.name}"',
                f'{character.name} has made a new post in the campaign.\n\nCampaign link: {link}'
            )
    return redirect(url_for('.campaign_posts', campaign_id=campaign_id))


@blueprint.route('/campaign/<int:campaign_id>/roll', methods=['GET', 'POST'])
@login_required
@csrf.exempt
def campaign_rolls(campaign_id):
    campaign = Campaign.query.get(campaign_id)
    if not campaign:
        return 'Could not find campaign with that id', 404
    character = current_user.get_character_in_campaign(campaign)
    if not character:
        return 'You are not a member of that campaign', 403
    if request.method == 'POST':
        roll_str = request.json.get('roll')
        if not roll_str:
            return '', 400
        roll = roll_dice(character, roll_str)
        db.session.add(roll)
        current_app.logger.info(f'User {current_user.id} as character {character.id} rolled str "{roll_str}"')
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
        current_app.logger.info(f'User {current_user.id} edited post {post.id}')
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
        character.campaign_join_note = request.form['notes']
        db.session.commit()
        current_app.logger.info(f'User {current_user.id} as character {character.id} requested to join campaign {campaign.id}')
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
                if character.user.email_for_accepted:
                    send_email(
                        [character.user.email],
                        'Your campaign join request has been approved',
                        f'Your request to join "{campaign.name}" has been approved for your character {character.name}'
                    )
                current_app.logger.info(f'User {current_user.id} accepted {character.id} to campaign {campaign.id}')
                flash('Character accepted')
            else:
                current_app.logger.info(f'User {current_user.id} denied {character.id} to campaign {campaign.id}')
                db.session.delete(character)
                flash('Character denied')
            db.session.commit()
            return redirect(url_for('.campaign_dm_controls', campaign_id=campaign_id))
        elif form_type == 'name_description':
            campaign.name = request.form['name']
            campaign.description = request.form['description']
            db.session.commit()
            flash('Campaign name/desciption updated')
            current_app.logger.info(f'User {current_user.id} updated campaign {campaign.id} name or description')
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
            current_app.logger.warning(f'Incorrect login for "{email}"')
            flash('Login failed', 'error')
            return redirect(url_for('.profile_login'))
        flash('Login successful')
        login_user(user, remember=True)
        current_app.logger.info(f'User {current_user.id} logged in')
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
        login_user(new_user, remember=True)
        current_app.logger.info(f'User {current_user.id} registered')
        flash('Login successful')
        return redirect(url_for('.campaigns'))
    return render_template('register.jinja2')


@blueprint.route('/profile/reset_password', methods=['GET', 'POST'])
def profile_reset_password():
    if request.method == 'POST':
        email = request.form['email']
        user = User.query.filter_by(email=email).first()
        if not user:
            flash('Password reset link sent')
            return redirect(url_for('.profile_login'))
        if user:
            key = create_password_reset_key(user.email)
            link = url_for('.profile_reset_password_confirm', email=user.email, key=key, _external=True)
            send_email([user.email], 'Password reset link', link)
            flash('Password reset link sent')
            current_app.logger.info(f'User {current_user.id} requested password reset link')
            return redirect(url_for('.profile_login'))
    return render_template('reset_password.jinja2')


@blueprint.route('/profile/reset_password/<email>/<key>', methods=['GET', 'POST'])
def profile_reset_password_confirm(email, key):
    user = User.query.filter_by(email=email).first()
    if not user:
        return redirect(url_for('.profile_login'))
    actual_key = get_password_reset_key(email)
    if not key == actual_key:
        flash('Wrong reset key', 'error')
        return redirect(url_for('.profile_login'))
    if request.method == 'POST':
        if not request.form['new_password'] == request.form['new_password_confirm']:
            flash('New passwords don\'t match', 'error')
            return redirect(url_for('.profile_reset_password_confirm', email=email, key=key))
        if not len(request.form['new_password']) > 5:
            flash('Password must be at least 5 characters long', 'error')
            return redirect(url_for('.profile_reset_password_confirm', email=email, key=key))
        user.set_password(request.form['new_password'])
        db.session.commit()
        clear_password_reset_keys(email)
        current_app.logger.info(f'User {current_user.id} updated password via reset link')
        flash('New password saved, please log in')
        return redirect(url_for('.profile_login'))
    return render_template('reset_password_confirm.jinja2', email=email, key=key)


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
            current_app.logger.info(f'User {current_user.id} created new character with name "{character.name}"')
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
            current_app.logger.info(f'User {current_user.id} deleted character {character.id}')
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
                current_app.logger.info(f'User {current_user.id} set character {character.id} name to "{new_value}"')
                character.name = new_value
            elif form_field == 'tag':
                current_app.logger.info(f'User {current_user.id} set character {character.id} tag to "{new_value}"')
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
            current_app.logger.info(f'User {current_user.id} updated post settings')
            flash('Post settings saved')
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
                current_app.logger.info(f'User {current_user.id} updated email settings')
                flash('Email settings saved')
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
            current_app.logger.info(f'User {current_user.id} updated password')
            flash('New password saved')
            return redirect(url_for('base.profile_settings'))
        elif settings_type == 'email_notifications':
            current_user.email_for_accepted = 'email_for_accepted' in request.form
            current_user.email_for_dm_post = 'email_for_dm_post' in request.form
            current_user.email_for_any_post = 'email_for_any_post' in request.form
            db.session.commit()
            current_app.logger.info(f'User {current_user.id} updated email notification settings')
            flash('Email settings saved')
            return redirect(url_for('base.profile_settings'))
        else:
            flash('Unknown setting value', 'error')
            return redirect(url_for('base.profile_settings'))
    return render_template('profile_settings.jinja2')


@blueprint.route('/profile/logout')
def profile_logout():
    current_app.logger.info(f'User {current_user.id} logged out')
    logout_user()
    return redirect(url_for('.profile_login'))


def send_email(recipients, subject, body):
    current_app.logger.info('Sending email to "{}" with subject "{}"'.format(
        ', '.join(recipients),
        subject
    ))
    return _send_email.apply_async(args=[
        current_app.config['EMAIL_API_KEY'],
        current_app.config['EMAIL_DOMAIN'],
        current_app.config['EMAIL_FROM'],
        recipients,
        subject,
        body
    ])
