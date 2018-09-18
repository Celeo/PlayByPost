from flask import Flask

from .shared import db
from .models import *


app = Flask(__name__)
app.config.from_json('config.json')
db.app = app
db.init_app(app)
db.create_all()


@app.route('/')
def index():
    return 'Index'
