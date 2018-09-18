#!/bin/bash
source env/bin/activate
FLASK_APP=pbp/app.py FLASK_ENV=development flask run --port 5000
