// Latihan 4: Queue[T any] dengan Enqueue dan Dequeue() (T, bool).
package main

import "fmt"

type Queue[T any] struct{ items []T }

func (q *Queue[T]) Enqueue(x T) { q.items = append(q.items, x) }

func (q *Queue[T]) Dequeue() (T, bool) {
	var nol T
	if len(q.items) == 0 {
		return nol, false
	}
	x := q.items[0]
	q.items[0] = nol // lepas referensi agar GC bisa membebaskan (penting untuk T pointer)
	q.items = q.items[1:]
	return x, true
}

func (q *Queue[T]) Len() int { return len(q.items) }

func main() {
	var q Queue[int]
	q.Enqueue(10)
	q.Enqueue(0) // nilai 0 SAH sebagai data
	q.Enqueue(30)

	for {
		v, ok := q.Dequeue()
		if !ok {
			fmt.Println("antrean kosong")
			break
		}
		fmt.Println("dequeue:", v)
	}

	// jika hanya mengembalikan zero value, 0 yang sah tidak bisa dibedakan dari "kosong"
	q.Enqueue(0)
	v, ok := q.Dequeue()
	fmt.Println("elemen 0 yang sah -> v:", v, "ok:", ok)
	v, ok = q.Dequeue()
	fmt.Println("antrean kosong     -> v:", v, "ok:", ok)
}
