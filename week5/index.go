package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
)

type Gipotenusa struct {
	a float64
	b float64
}

func FindGipotenusa(a float64, b float64) float64 {
	result := math.Sqrt(math.Pow(a, 2) + math.Pow(b, 2))
	return result
}

func Home(w http.ResponseWriter, r *http.Request) {
	g := Gipotenusa{a: 1, b: 1}
	fmt.Fprintf(w, "The result is: "+strconv.FormatFloat(FindGipotenusa(g.a, g.b), 'g', 17, 64))
}

func StartServer() {
	http.HandleFunc("/", Home)
	http.ListenAndServe(":8080", nil)
}

func main() {
	// var a float64 = 3
	// var b float64 = 4
	// fmt.Sprintf("The result is: %d", FindGipotenusa(a, b))
	StartServer()
}
