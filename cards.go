package main

import "fmt"
import "math/rand"
import "time"
import "sync"


func getDeck() []int {
    rand.Seed(time.Now().UTC().UnixNano())
    deck := rand.Perm(52)
    for i, _ := range deck {
        deck[i]++
    }
    return deck
}

func sum_channel(ch chan int) int {
	total := 0
	for val := range ch {
        total += val
	}
	return total
}


func isMatchFound(input chan int, wg *sync.WaitGroup) {


    deck1 := getDeck()
    deck2 := getDeck()
    count := 0
    for i:=0; i < 52; i++ {
        if deck1[i] == deck2[i] {
            count = 1
            break
        }
    }
    input <- count
    wg.Done()

}


func main() {

    iterations := 10000000
    matches := make(chan int, iterations)
    matches_count := 0
    var wg sync.WaitGroup


    for i:=0; i < iterations; i++ {
        wg.Add(1)
        go isMatchFound(matches, &wg)
    }

    wg.Wait()
    close(matches)

    fmt.Printf("finished")
    matches_count = sum_channel(matches)
    percentage := (float64(matches_count) / float64(iterations)) * 100
    fmt.Printf("found %d matches, which is %v%%:", matches_count, percentage)
}
