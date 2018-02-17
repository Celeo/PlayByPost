# Play by post

A full-stack web app for an incredibly simple forum designed for "play by post" RPGs.

* Front-end: Vue.js
* Back-end: Go (Gin)

## Running

```bash
for p in \
        github.com/gin-gonic/gin \
        github.com/itsjamie/gin-cors \
        github.com/jmoiron/sqlx \
        github.com/nu7hatch/gouuid \
        golang.org/x/crypto/bcrypt \
        github.com/mattn/go-sqlite3 \
        github.com/stretchr/testify
    do
    go get $p
done
```

TBD
