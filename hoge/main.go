package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const refreshInterval = 500 * time.Millisecond

var app *tview.Application
var timers map[string]*Timer

type Timer struct {
	Label     string
	TextView  *tview.TextView
	StartTime time.Time // タイマーの開始時刻
}

func (timer *Timer) currentTimeString() string {
	t := time.Now()
	return fmt.Sprintf(t.Format("current time: 15:04:05"))
}

func (timer *Timer) updateTime() {
	for {
		time.Sleep(refreshInterval)
		app.QueueUpdateDraw(func() {
			now := time.Now()
			elapsed := now.Sub(timer.StartTime)
			hours := int(elapsed.Hours())
			minutes := int(elapsed.Minutes()) % 60
			seconds := int(elapsed.Seconds()) % 60
			timer.TextView.SetText(fmt.Sprintf("Timer '%s': %02d:%02d:%02d", timer.Label, hours, minutes, seconds))
		})

	}
}

func main() {
	app = tview.NewApplication()
	timers = make(map[string]*Timer)

	// コマンド入力欄
	commandInputField := tview.NewInputField().SetLabel("Command: ")
	// タイマー表示欄
	timerView := tview.NewFlex().SetDirection(tview.FlexRow)

	// commandInputFieldのイベントハンドラを追加する
	commandInputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() != tcell.KeyEnter {
			return event
		}

		// 入力されたコマンドを取得
		command := commandInputField.GetText()
		commandArgs := strings.Fields(command)

		if len(commandArgs) != 2 {
			return event
		}

		cmd, timerName := commandArgs[0], commandArgs[1]

		switch cmd {
		case "start":
			// timerNameが既に存在する場合は何もしない
			if _, ok := timers[timerName]; ok {
				return event
			}

			// timer構造体を作成
			timer := &Timer{
				Label:     timerName,
				TextView:  tview.NewTextView(),
				StartTime: time.Now(),
			}

			timers[timerName] = timer

			go timer.updateTime()

			// timerViewにtextViewを追加
			timerView.AddItem(timer.TextView, 1, 1, false)

			// 入力欄をクリア
			commandInputField.SetText("")
		}

		return event
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(commandInputField, 0, 1, true).
		AddItem(timerView, 0, 9, false)

	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}
