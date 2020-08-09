package main

import (
	"go-web-console/routes"
)

func main() {
	g := routes.GetGin()
	g.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
