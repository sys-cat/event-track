export GOPATH="~/.go"
export GOROOT="/usr/local/go"
export PATH="$GOROOT/bin:$PATH"
$(go get github.com/gin-gonic/gin)
$(go get github.com/BurntSushi/toml)
$(go get github.com/go-sql-driver/mysql)
$(mysql -u root < ./event.sql)
