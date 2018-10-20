# Play by post

[![pipeline status](https://gitlab.com/Celeo/PlayByPost/badges/master/pipeline.svg)](https://gitlab.com/Celeo/PlayByPost/commits/master)
[![coverage report](https://gitlab.com/Celeo/PlayByPost/badges/master/coverage.svg)](https://gitlab.com/Celeo/PlayByPost/commits/master)

A full-stack web app for a simple forum designed for "play by post" RPGs.

## System requirements

1. Python 3.7
1. Redis

## Getting set up

1. Clone the repo
1. Create a virtualenv
1. Activate the virtualenv
1. Install Python deps
1. Copy 'config.example.json' to 'config.json'
1. Populate the config
1. Navigate to the project root

This project is set up to run with SQLite by default. If you want to run on some other
database, update 'config.json' accordingly and install any required libraries.

## Running

Run './run.sh' for local development and './run.sh prod' for production.

If you're using any of the emailing, you'll need a provider (I'm using Mailgun) and to run the celery
worker alongside your app (both in development and producton).

### Running tests

1. Run './run.sh tests'

## Deploying

Basically the same as the "Getting set up" instructions, but on whichever server you're running the app on. Then:

1. Setup a Nginx (or Apache) proxy
1. Run the app with './run.sh prod'
