from flask_sqlalchemy import SQLAlchemy
from flask_redis import Redis
from flask_wtf.csrf import CSRFProtect


db = SQLAlchemy()
redis = Redis()
csrf = CSRFProtect()
