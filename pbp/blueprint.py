from flask import Blueprint, render_template


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


@blueprint.route('/profile/join')
def profile_join():
    return 'VIEW: profile/join'


@blueprint.route('/profile/settings')
def profile_settings():
    return 'VIEW: profile/settings'


@blueprint.route('/profile/logout')
def profile_logout():
    return 'VIEW: profile/logout'
