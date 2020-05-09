package main

import (
	"log"

	"github.com/unqnown/smp/app"
	"github.com/unqnown/smp/pkg/check"
)

func main() {
	check.Fatal(app.Run())
}

func init() {
	log.SetFlags(0)
}
