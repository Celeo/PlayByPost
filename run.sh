#!/bin/bash
source env/bin/activate

export FLASK_APP=pbp/app.py
export FLASK_ENV=development

if [ "$1" == "shell" ]; then
    flask shell
else
    flask run --port 5000
fi
