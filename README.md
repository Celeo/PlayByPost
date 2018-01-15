# Play by post

A full-stack web app for an incredibly simple forum designed for "play by post" RPGs.

Front-end: Vue.js

Back-end: Go (Gin)

## Running

This project is set up to run the backend on [Heroku](https://heroku.com/) and the frontend on [Surge](http://surge.sh/).

When you clone this project, you'll need to create a Heroku app with the [Go buildpack](https://elements.heroku.com/buildpacks/heroku/heroku-buildpack-go).
You'll need the Postgres addon. You can deploy the server with `git subtree push --prefix server heroku master`. This command is more complicated than the
standard git push that you've probably used before for Heroku because the part of this repo that needs to go to Heroku isn't in the root level. As such,
only that level, the `./server` directory, is pushed to Heroku. You'll also want to set the GIN_MODE envrionment variable to "release" with
`heroku config:get GIN_MODE --app [app name]`.

You'll also need to setup a Surge app. Navigate into the `./client` directory, run `yarn build` or `npm run build` to create the inner `./dist` directory, then run
`surge ./dist`. Use the default name or supply your own. After Surge finishes, create a file called `CNAME`, with the name of your app ("whatever.surge.sh").