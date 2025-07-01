// queue.go
// Ejemplo de patrón productor-consumidor usando un canal como cola de trabajos.

package main

import (
	"fmt"
	"sync"
	"time"
)

// ------------------------------------------------------------
// producer: genera un número fijo de items y los encola.
// ------------------------------------------------------------
func producer(name string, count int, queue chan<- string) {
	for i := 1; i <= count; i++ {
		item := fmt.Sprintf("%s-item-%d", name, i)
		fmt.Printf("%s produjo %s\n", name, item)
		queue <- item                                          // encola el item (bloquea si la cola está llena)
		time.Sleep(time.Duration(100+i*20) * time.Millisecond) // simula variabilidad
	}
}

// ------------------------------------------------------------
// consumer: lee items de la cola hasta que se cierre el canal.
// ------------------------------------------------------------
func consumer(id int, queue <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()           // marca este consumer como terminado al salir
	for item := range queue { // itera mientras el canal esté abierto
		fmt.Printf("Consumidor %d procesando %s\n", id, item)
		time.Sleep(200 * time.Millisecond) // simula tiempo de procesamiento
	}
	fmt.Printf("Consumidor %d ha terminado\n", id)
}

func main() {
	// 1) Creamos un canal con buffer para usarlo como cola
	queue := make(chan string, 5)

	var wg sync.WaitGroup

	// 2) Lanzamos una goroutine que actúa como agrupador de productores
	go func() {
		// Producimos 5 items con dos productores distintos
		producer("ProdA", 5, queue)
		producer("ProdB", 5, queue)
		close(queue) // importante: cierra la cola para que los consumers salgan del range
	}()

	// 3) Arrancamos 3 consumidores en paralelo
	numConsumers := 3
	wg.Add(numConsumers)
	for i := 1; i <= numConsumers; i++ {
		go consumer(i, queue, &wg)
	}

	// 4) Esperamos a que todos los consumidores terminen
	wg.Wait()
	fmt.Println("Todos los consumidores han terminado. Saliendo.")
}
