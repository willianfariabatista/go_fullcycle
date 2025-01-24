package main

// quando variavel for do tipo "const" éla é imutavél ou seja o valor atribuido não muda.!
const a = "Hello World and Happy New Year!"

/*
// poder ser declarada individualmente:

//declaração da variavel tipo bollean.
 var b bool

// declaração da variavel tipo inteiro.
 var c int

// declaração da variavel tipo string
 var d string

*/

// ou pode ser declarado em um grupo de variaveis
var (

	// declaração da variavel tipo bollean.
	//b bool

	// podemos declarar e atribuir um valor
	b bool = true

	// declaração da variavel tipo inteiro.
	//c int

	// podemos declarar e atribuir um valor
	c int = 10

	// declaração da variavel tipo string
	//d string

	// podemos declarar e atribuir um valor
	d string = "Willian"

	// declaração da variavel tipo float
	//e float64

	// podemos declarar e atribuir um valor
	e float64 = 9.9
)

func main() {

	// atribuindo valor a variavel b
	b = true

	// declaração variavel de escopo LOCAL (outras funções não consegue acessar o valor de f se declarada dentro desta função)
	//var f string

	// podemos declarar e atribuir um valor
	//var f = "teste"

	// podemos fazer a declação + atribuição na primeira vez que a variavel é decladara desta outra forma tambem:
	f := "teste"

	// Na segunda vez apenas atribuimos pois ele ja foi declarada antes
	f = "teste2"

	println(a)
	println(b)
	println(c)
	println(d)
	println(e)
	println(f)
}
