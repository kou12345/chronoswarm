package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

/*
TODO

start Timer1してから
stop Timer1をすると、アプリケーションが停止する

*/

var app *tview.Application
var textView *tview.TextView

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
	fmt.Fprintf(textView, "Timer '%s' started\n", t.Label)

	go func() {
		for {
			select {
			case <-t.Ticker.C:
				if t.Finished {
					return
				}
				app.QueueUpdateDraw(func() {
					fmt.Fprintf(textView, "\rTimer '%s': %s", t.Label, time.Since(t.Start).Round(time.Second))
				})
			}
		}
	}()
}

func (t *Timer) StopTimer() {
	if t.Ticker != nil {
		t.Ticker.Stop()
		t.Finished = true
		app.QueueUpdateDraw(func() {
			fmt.Fprintf(textView, "\nTimer '%s' stopped at %s\n", t.Label, time.Since(t.Start).Round(time.Second))
		})
	}
}

func main() {

	app = tview.NewApplication()

	timers := make(map[string]*Timer)

	// textViewにTimerの情報を表示する
	textView = tview.NewTextView()
	textView.SetTitle("textView")
	textView.SetBorder(true)

	// inputFieldにコマンドを入力する
	inputField := tview.NewInputField()
	inputField.SetLabel("input: ")
	inputField.SetTitle("inputField").
		SetBorder(true)

	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			// ここで入力されたコマンドを処理する
			// 入力された文字列を取得する
			inputText := inputField.GetText()

			// inputTextをスペースで分割する
			// args[0]にコマンド、args[1]にタイマー名が入る
			// 例: "start Timer1" -> args[0] = "start", args[1] = "Timer1"
			args := strings.SplitN(inputText, " ", 2)

			switch args[0] {
			case "start":
				if len(args) < 2 {
					fmt.Fprintf(textView, "Invalid input. Please enter a command and a timer name.\n")
					break
				}

				if timer, exists := timers[args[1]]; exists && !timer.Finished {
					fmt.Fprintf(textView, "Timer is already running. Please stop it before starting a new one.\n")
					break
				}

				newTimer := &Timer{Label: args[1], Finished: true} // タイマーが終了したと仮定し、再起動を許可する
				timers[args[1]] = newTimer
				newTimer.StartTimer()

			case "stop":
				if len(args) < 2 {
					fmt.Fprintf(textView, "Invalid input. Please enter a command and a timer name.\n")
					break
				}

				if timer, exists := timers[args[1]]; exists && !timer.Finished {
					timer.StopTimer()
				} else {
					fmt.Fprintf(textView, "No active timer to stop.\n")
				}

			case "exit":
				fmt.Fprintf(textView, "Exiting application.\n")
				app.Stop()
			default:
				fmt.Fprintf(textView, "Invalid command. Please enter 'start', 'stop', or 'exit'.\n")
			}

			// 入力フィールドをクリアする
			inputField.SetText("")

			return nil
		}
		return event
	})

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow).
		AddItem(inputField, 3, 0, true).
		AddItem(textView, 0, 1, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
