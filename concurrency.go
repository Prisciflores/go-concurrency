// concurrency.go
// Ejemplo de goroutines, channels y select en Go.

package main

import (
	"fmt"
	"math/rand"
	"time"
)

// ------------------------------------------------------------
// worker: simula trabajo y envía un mensaje por el channel.
// ------------------------------------------------------------
// id  = identificador del worker (solo para impresión).
// ch  = canal de solo envío (chan<- string).
func worker(id int, ch chan<- string) {
	// 1) Calculamos un delay aleatorio entre 0 y 1000 ms
	delay := time.Duration(rand.Intn(1000)) * time.Millisecond
	// 2) Simulamos que el worker “trabaja” durmiendo ese tiempo
	time.Sleep(delay)
	// 3) Cuando termina, envía un mensaje al canal
	ch <- fmt.Sprintf("Worker %d finalizó en %v", id, delay)
}

// ------------------------------------------------------------
// fanIn: combina dos canales en un solo canal de salida.
// ------------------------------------------------------------
// ch1, ch2 = canales de solo recepción (<-chan string).
// out      = canal que reunirá mensajes de ambos.
func fanIn(ch1, ch2 <-chan string) <-chan string {
	out := make(chan string)
	// Goroutine que redirige de ch1 a out
	go func() {
		for msg := range ch1 {
			out <- msg
		}
	}()
	// Goroutine que redirige de ch2 a out
	go func() {
		for msg := range ch2 {
			out <- msg
		}
	}()
	return out
}

func main() {
	// Semilla para rand para que los delays varíen cada ejecución
	rand.Seed(time.Now().UnixNano())

	// 1) Creamos dos canales sin buffer para dos workers
	ch1 := make(chan string)
	ch2 := make(chan string)

	// 2) Lanzamos dos goroutines en paralelo
	go worker(1, ch1)
	go worker(2, ch2)

	// 3) Mezclamos ambos canales en uno
	merged := fanIn(ch1, ch2)

	// 4) Esperamos 2 mensajes, con timeout de 500ms cada vez
	for i := 0; i < 2; i++ {
		select {
		case msg := <-merged:
			// Caso normal: recibimos mensaje de algún worker
			fmt.Println(msg)
		case <-time.After(500 * time.Millisecond):
			// Timeout: si tarda más de 500ms, salimos del select
			fmt.Println("Timeout esperando respuesta")
		}
	}

	fmt.Println("Main finaliza")
}
