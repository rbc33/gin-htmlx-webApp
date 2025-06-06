GOOSE_UP = """GOOSE_DRIVER="mysql" GOOSE_DBSTRING="root:secret@tcp(192.168.0.100:33060)/gocms" goose up"""
MY_SQL_COMMAND= mysql -h 192.168.0.100 -P 33060 -u root -p
go_lint= golangci-lint run --disable-all -E errcheck -E staticcheck -E unused -E gosimple -E gofmt
docker= docker run -it -v /Users/ric/code/go/gocms:/gocms --entrypoint bash rbenthem/gocms ,
docker run -it -v -p 8080:8080  /Users/ric/code/go/gocms:/gocms --entrypoint sh rbenthem/gocms