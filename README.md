# Play by post

A full-stack web app for a simple forum designed for "play by post" RPGs.

## Getting set up

1. Clone the repo
1. Create a virtualenv
1. Activate the virtualenv
1. Install Python deps
1. Navigate into './widgets'
1. Install Node deps
1. Run the 'build_prod' package.json target
1. Navigate to '../pbp'
1. Copy 'config.example.json' to 'config.json'
1. Populate the config
1. Navigate to the project root

This project is set up to run with SQLite by default. If you want to run on some other
database, update 'config.json' accordingly and install any required libraries.

## Running

1. Run './run.sh' for local development and './run.sh prod' for production.

### Running tests

1. Run './run.sh tests'

## Deploying

TBD
