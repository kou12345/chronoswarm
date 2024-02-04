package main

import (
	"context"
	"fmt"
	"log"
	"os"
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
var logView *tview.TextView

type Timer struct {
	Label    string
	Start    time.Time
	Ticker   *time.Ticker
	Finished bool
	Ctx      context.Context
	Cancel   context.CancelFunc
}

func (t *Timer) StartTimer() {
	t.Start = time.Now()
	t.Finished = false
	t.Ticker = time.NewTicker(1 * time.Second)
	t.Ctx, t.Cancel = context.WithCancel(context.Background())

	fmt.Fprintf(textView, "Timer '%s' started\n", t.Label)

	go func() {
		for {
			select {
			case <-t.Ticker.C:
				app.QueueUpdateDraw(func() {
					fmt.Fprintf(textView, "Timer '%s': %s \r", t.Label, time.Since(t.Start).Round(time.Second))
				})
			case <-t.Ctx.Done():
				log.Println("LOG: Timer stopped")
				fmt.Fprintf(logView, "Timer '%s' finished\n", t.Label)
				return
			}
		}
	}()
}

func (t *Timer) StopTimer() {
	log.Println("START: StopTimer()")

	if t.Ticker != nil {
		t.Ticker.Stop()
		log.Println("LOG: Ticker stopped")
		t.Finished = true
		log.Println("LOG: Finished set to true")
		t.Cancel()

		// ! これが実行されるとアプリケーションが停止する
		if app != nil {
			log.Println("LOG: app is not nil")
			app.QueueUpdateDraw(func() {
				log.Println("LOG: QueueUpdateDraw")
				fmt.Fprintf(textView, "\nTimer '%s' stopped at %s\n", t.Label, time.Since(t.Start).Round(time.Second))
			})
		}

	}

	log.Println("END: StopTimer()")
}

func main() {

	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("ERROR: logファイルを開けませんでした: %s", err)
	}
	defer logFile.Close()

	// logの出力先をファイルに変更
	log.SetOutput(logFile)

	app = tview.NewApplication()

	timers := make(map[string]*Timer)

	// textViewにTimerの情報を表示する
	textView = tview.NewTextView()
	textView.SetTitle("textView")
	textView.SetBorder(true)

	// logViewにログを表示する
	logView = tview.NewTextView()
	logView.SetTitle("logView")
	logView.SetBorder(true)

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
					fmt.Fprintf(logView, "Invalid input. Please enter a command and a timer name.\n")
					break
				}

				if timer, exists := timers[args[1]]; exists && !timer.Finished {
					fmt.Fprintf(logView, "Timer is already running. Please stop it before starting a new one.\n")
					break
				}

				newTimer := &Timer{Label: args[1], Finished: false}
				timers[args[1]] = newTimer
				newTimer.StartTimer()

			case "stop":
				if len(args) < 2 {
					fmt.Fprintf(logView, "Invalid input. Please enter a command and a timer name.\n")
					break
				}

				if timer, exists := timers[args[1]]; exists && !timer.Finished {
					timer.StopTimer()
				} else {
					fmt.Fprintf(logView, "No active timer to stop.\n")
				}

			case "exit":
				fmt.Fprintf(logView, "Exiting application.\n")
				app.Stop()
			default:
				fmt.Fprintf(logView, "Invalid command. Please enter 'start', 'stop', or 'exit'.\n")
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
		AddItem(logView, 0, 1, false).
		AddItem(textView, 0, 4, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
