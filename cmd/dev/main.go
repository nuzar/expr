package main

import (
	"log"

	"github.com/nuzar/expr"
)

func main() {
	// fmt.Println("now")
	// expr.Run("print(now())")

	log.Println("gnow")
	res, err := expr.Run("print(gnow())")
	log.Printf("res: %+v, err: %s", res, err)
}
