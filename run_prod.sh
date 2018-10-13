#!/bin/bash
source env/bin/activate

gunicorn -w 5 -b 127.0.0.1:5000 pbp:app
