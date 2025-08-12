package main

import (
	"flag"
	"fmt"
	"github.com/moverq1337/wbOrder/internal/app/db"
)

var migrate = flag.Bool("migrate", false, "Run database migration")

func main() {
	flag.Parse()
	db.Connection()
	if *migrate {
		db.Migration()
	} else {
		fmt.Println("Migration skipped, because u run it with out flag (-migration)")
	}

}
