export GOPATH=~/.go
export GOROOT="/usr/local/go"
export PATH="$GOROOT/bin:$PATH"
cd ~/event_track
$(go get github.com/gin-gonic/gin)
$(go get github.com/BurntSushi/toml)
$(go get github.com/go-sql-driver/mysql)
$(mysql -u root < ./event.sql)
$(go run ./main.go)
