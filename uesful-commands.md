rebuild container -> docker-compose up -d --no-deps --build < service >
goose add migrationfile -> goose create add_page_table sql
goose migrate -> GOOSE_DRIVER="mysql" GOOSE_DBSTRING="root:secret@tcp(192.168.0.100:33060)/gocms" goose up

make-release -> 
    git commit -m ""
    git tag v0.1.3_alpha
    git push origin v0.1.3_alpha
    GOPROXY=proxy.golang.org go list -m github.com/rbc33/gocms@v0.1.3_alpha