package main

// importando bibliotecas nativa do GO
//import "fmt"

// para importar mais de uma biblioteca em um unico import
import (
	"fmt"
	//"XPTO"
)

// quando variavel for do tipo "const" éla é imutavél ou seja o valor atribuido não muda.!
const a = "Hello World and Happy New Year!"

// podemos criar um novo "tipo"
type ID int

var (
	aa int = 10

	// ID agora é um novo "tipo"
	c ID = 99
)

func main() {

	// biblioteca "fmt" imprimir o tipo da variavel
	fmt.Printf("O tipo de a é %T", aa)

	// biblioteca "fmt" imprimir o valor da variavel
	fmt.Printf("O valor de a é %v", aa)

}
