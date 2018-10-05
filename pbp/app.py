from flask import Flask, render_template

from .shared import db


app = Flask(__name__)
app.config.from_json('config.json')
db.app = app
db.init_app(app)


@app.route('/')
def index():
    return render_template('index.html')


@app.route('/campaigns')
def campaigns():
    return 'VIEW: campaigns'


@app.route('/search')
def search():
    return 'VIEW: search'


@app.route('/glossary')
def glossary():
    return 'VIEW: glossary'


@app.route('/help')
def help():
    return 'VIEW: help'


@app.route('/profile/login')
def profile_login():
    return 'VIEW: profile/login'


@app.route('/profile/join')
def profile_join():
    return 'VIEW: profile/join'


@app.route('/profile/settings')
def profile_settings():
    return 'VIEW: profile/settings'


@app.route('/profile/logout')
def profile_logout():
    return 'VIEW: profile/logout'