package main

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var lotteryNumbers [6]int
var lotteryDrawIndex int = 0

const port = "8090"

// Writes the current lottery number into the passed ResponseWrite object as plain text
func getLotteryNumbers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<p>Lottery numbers: %d</p>", lotteryNumbers)
	if lotteryDrawIndex <= 5 {
		fmt.Fprintf(w, "<p>%d numbers drawn</p>", lotteryDrawIndex)
		fmt.Fprintf(w, "<p><a href=\"http://localhost:%s\">reload page</a></p>", port)
	} else {
		fmt.Fprintf(w, "<p><a href=\"http://localhost:%s/restartLottery\">restart lottery</a></p>", port)
	}
}

func restartLottery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	lotteryDrawIndex = 0
	lotteryNumbers = [6]int{0, 0, 0, 0, 0, 0}
	fmt.Fprintf(w, "Lottery restarted")
	fmt.Fprintf(w, "<p><a href=\"http://localhost:%s\">show lottery numbers</a></p>", port)
}

func main() {
	// http listener
	// listen on main path "/"
	// this http listener returns a list of drawn numbers
	http.HandleFunc("/", getLotteryNumbers)
	http.HandleFunc("/restartLottery", restartLottery)

	// receive signals to stop the server
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// generate six lottery numbers
	go func() {
		for {
			// generate a random number ever 10th second
			select {
			case <-time.NewTicker(time.Second).C:
				if lotteryDrawIndex <= 5 {
					addRandomLotteryNumber()
				}
			case <-sigs:
				fmt.Println("server received signal. Shutting down...")
				break
			}
		}
	}()

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("An error occured while listring to incoming http requests on", port, "error:", err)
	}
}

func addRandomLotteryNumber() {
	randomNumber := rand.IntN(48) + 1 // a number between 1 and 49
	lotteryNumbers[lotteryDrawIndex] = randomNumber
	lotteryDrawIndex += 1
}
