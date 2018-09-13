# Play by post

A full-stack web app for a simple forum designed for "play by post" RPGs.

* Front-end: Vue.js
* Back-end: Go (Gin)

## Getting set up

Front-end:

```bash
cd client
yarn
```

You can use `npm i(nstall)` instead of `yarn` if you don't use yarn.

Back-end:

```bash
for p in \
        github.com/gin-gonic/gin \
        github.com/itsjamie/gin-cors \
        github.com/jmoiron/sqlx \
        github.com/nu7hatch/gouuid \
        golang.org/x/crypto/bcrypt \
        github.com/mattn/go-sqlite3 \
        github.com/stretchr/testify \
        golang.org/x/text/search \
        github.com/foolin/gin-template
    do
    go get $p
done
```

## Running

Run the front-end with `yarn dev` or `npm run dev` (inside the client/ directory).

Run the back-end with `make` (inside the server/ directory).

## Deploying

If you want to deploy like how I'm deploying, you'll need a server with an SSH-based login called "playbypost" in your SSH config,
and the app to be running at /srv/ on the server with the appropriate Nginx configuration and systemd script.

Otherwise disregard the Makefile at the root level and deploy however you prefer. The server binary is built with `make build` (i.e.
`go build` in server/), and the minified client files are built with `yarn dist` or `npm run dist` in client/.
