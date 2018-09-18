<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.1.3/css/bootstrap.min.css" integrity="sha384-MCw98/SFnGE8fJT3GXwEOngsV7Zt27NXFoaoApmYm81iuXoPkFOJwJ8ERdknLPMO" crossorigin="anonymous">
    <link rel="stylesheet" href="static/style.css">
    <title>Play by Post{{if .campaign}}- {{.campaign}}{{end}}</title>
  </head>
  <body>
    <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
      <div class="container-fluid">
        <span class="navbar-brand">Play by Post</span>
        <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarSupportedContent" aria-controls="navbarSupportedContent">
          <span class="anvbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarSupportedContent">
          <ul class="navbar-nav flex-row ml-md-auto d-none d-md-flex">
            <li class="nav-item"><a href="/campaigns" class="nav-link">Campaign selector</a></li>
            <li class="nav-item dropdown">
              <a href="#" class="nav-item nav-link dropdown-toggle" data-toggle="dropdown">Tools</a>
              <div class="dropdown-menu bg-dark">
                <a class="dropdown-item text-light" href="/search">Search</a>
                <a class="dropdown-item text-light" href="/glossary">Glossary</a>
                <a class="dropdown-item text-light" href="/help">Help</a>
              </div>
            </li>
            <li class="nav-item dropdown">
              <a href="#" class="nav-item nav-link dropdown-toggle" data-toggle="dropdown">Profile</a>
              <div class="dropdown-menu bg-dark">
                <a class="dropdown-item text-light" href="/profile/login">Log in</a>
                <a class="dropdown-item text-light" href="/profile/join">Join</a>
                <a class="dropdown-item text-light" href="/profile/settings">Settings</a>
                <a class="dropdown-item text-light" href="/profile/logout">Log out</a>
              </div>
            </li>
          </ul>
        </div>
      </div>
    </nav>
    <main class="py-md-4 pl-md-5 pr-md-5" role="main">
      <div class="container-fluid">
        {{template "content" .}}
      </div>
    </main>
    </div>
    <script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.3/umd/popper.min.js" integrity="sha384-ZMP7rVo3mIykV+2+9J3UJ46jBk0WLaUAdn689aCwoqbBJiSnjAK/l8WvCWPIPm49" crossorigin="anonymous"></script>
    <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.1.3/js/bootstrap.min.js" integrity="sha384-ChfqqxuZUCnJSK3+MXmPNIyE6ZbWh2IMqE241rYiqJxyMiZ6OW/JmZQ5stwEULTy" crossorigin="anonymous"></script>
    <script src="https://vuejs.org/js/vue.js"></script>
  </body>
</html>
