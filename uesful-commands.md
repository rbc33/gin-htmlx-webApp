rebuild container -> docker-compose up -d --no-deps --build < service >
goose add migrationfile -> goose create add_page_table sql
goose migrate -> GOOSE_DRIVER="mysql" GOOSE_DBSTRING="root:secret@tcp(192.168.0.100:33060)/gocms" goose up