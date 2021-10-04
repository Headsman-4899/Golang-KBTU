package main

import (
	"fmt"      // вывод результата
	"math"     // нахождение корня и возведения в степень
	"net/http" // запуск сервера
	"strconv"  // для конвертирования из float64 в string
)

type Gipotenusa struct { // обычная структура с 2 параметрами
	a float64
	b float64
}

func FindGipotenusa(a float64, b float64) float64 { // функиция находит гипотенузу
	result := math.Sqrt(math.Pow(a, 2) + math.Pow(b, 2))
	return result
}

func Home(w http.ResponseWriter, r *http.Request) { // наша главная страница
	g := Gipotenusa{a: 1, b: 1}
	fmt.Fprintf(w, "The result is: "+strconv.FormatFloat(FindGipotenusa(g.a, g.b), 'g', 17, 64))
}

func StartServer() { // отвечает за запуск сервера
	http.HandleFunc("/", Home)
	http.ListenAndServe(":8080", nil)
}

func main() {
	StartServer() // запуск сервера
}
