#!/bin/bash
source env/bin/activate

export FLASK_APP=pbp/app.py

if [ "$1" == "shell" ]; then
    flask shell
elif [ "$1" == "tests" ]; then
    pytest --cov=pbp pbp/tests
elif [ "$1" == "prod" ]; then
    gunicorn -w 5 -b 127.0.0.1:5000 pbp:app
elif [ "$1" == "" ]; then
    export FLASK_ENV=development
    flask run --port 5000
else
    echo 'Unknown param'
fi
