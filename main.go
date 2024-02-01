package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Timer struct {
	Label    string
	Start    time.Time
	Ticker   *time.Ticker
	Finished bool
}

func (t *Timer) StartTimer() {
	t.Start = time.Now()
	t.Finished = false
	t.Ticker = time.NewTicker(1 * time.Second)
	fmt.Printf("Timer '%s' started\n", t.Label)

	go func() {
		for {
			select {
			case <-t.Ticker.C:
				if t.Finished {
					return
				}
				fmt.Printf("\rTimer '%s': %s", t.Label, time.Since(t.Start).Round(time.Second))
			}
		}
	}()
}

func (t *Timer) StopTimer() {
	t.Ticker.Stop()
	t.Finished = true
	fmt.Printf("\nTimer '%s' stopped at %s\n", t.Label, time.Since(t.Start).Round(time.Second))
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	timers := make(map[string]*Timer)

	for {
		fmt.Print("Enter command and timer name (e.g., start Timer1): ")
		scanner.Scan()
		input := scanner.Text()
		args := strings.SplitN(input, " ", 2)

		if len(args) < 2 {
			fmt.Println("Invalid input. Please enter a command and a timer name.")
			continue
		}

		command := args[0]
		label := args[1]

		switch command {
		case "start":
			if timer, exists := timers[label]; exists && !timer.Finished {
				fmt.Println("Timer is already running. Please stop it before starting a new one.")
			} else {
				newTimer := &Timer{Label: label, Finished: true} // タイマーが終了したと仮定し、再起動を許可する
				timers[label] = newTimer
				newTimer.StartTimer()
			}
		case "stop":
			if timer, exists := timers[label]; exists && !timer.Finished {
				timer.StopTimer()
			} else {
				fmt.Println("No active timer to stop.")
			}
		case "exit":
			fmt.Println("Exiting application.")
			return
		default:
			fmt.Println("Invalid command. Please enter 'start', 'stop', or 'exit'.")
		}
	}
}
