# Synopsis

### use in lib

```golang
package mylib

import (
    "github.com/raohwork/envexist"
    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// read db config from env var and create *sql.DB
func init() {
    m, data := envexist.New("mylib", mysetup)
    m.Need("DB_CONN", "connect string for mysql", "user:pass@tcp(mysql)/mydb")
    m.May("DB_PARAM", "additional param for constr", "parseTime=true")
}

func mysetup(env map[string]string) {
    var err error
    constr := env["DB_CONN"] 
    if param := env["DB_PARAM"]; param != "" {
        constr += "?" + param
    }
    if DB, err = sql.Open("mysql", constr); err != nil {
        log.Fatal(err)
    }
}
```

### in application entry

```golang
package main

import (
    "github.com/raohwork/envexist"
    "./mylib"
)

func main() {
    m, ch := envexist.Main("myprog")
    m.Want("DEBUG", "any non-empty value enables debug mode", "")
    if !envexist.Parse() {
        envexist.PrintEnvList()
        os.Exist(1)
    }
    
    // codes below will always run after mylib.mysetup()
    data := <- ch
    if data["DEBUG"] != "" {
        // enable debug mode here
    }
    
    // other application codes
}
```

# License

MPL 2.0
