package main

// quando variavel for do tipo "const" éla é imutavél ou seja o valor atribuido não muda.!
const a = "Hello World and Happy New Year!"

// podemos criar um novo "tipo"
type ID int

var (

	// podemos declarar e atribuir um valor
	b bool = true

	// ID agora é um novo "tipo"
	c ID = 99
)

func main() {

	println(a)
	println(b)
	println(c)
}
