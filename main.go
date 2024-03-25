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
	label     string
	textView  *tview.TextView
	startTime time.Time     // タイマーの開始時刻
	elapsed   time.Duration // タイマーの経過時間
	stopChan  chan struct{} // タイマーを停止するためのチャンネル
	isRunning bool          // タイマーが実行中かどうかを示すフラグ
}

func (timer *Timer) updateTime() {
	timer.isRunning = true
	for timer.isRunning {
		select {
		case <-timer.stopChan: // タイマーを停止
			timer.isRunning = false
			return
		case <-time.After(refreshInterval):
			app.QueueUpdateDraw(func() {
				now := time.Now()
				elapsed := now.Sub(timer.startTime)
				hours := int(elapsed.Hours())
				minutes := int(elapsed.Minutes()) % 60
				seconds := int(elapsed.Seconds()) % 60
				timer.textView.SetText(fmt.Sprintf("Timer '%s': %02d:%02d:%02d", timer.label, hours, minutes, seconds))
			})
		}
	}
}

func startTimer(timerName string, timerView *tview.Flex) {
	if _, ok := timers[timerName]; ok {
		// 既に存在する場合は何もしない
		return
	}

	timer := &Timer{
		label:     timerName,
		textView:  tview.NewTextView(),
		startTime: time.Now(),
		stopChan:  make(chan struct{}),
		isRunning: true,
	}

	timers[timerName] = timer

	go timer.updateTime()

	// timerViewにtextViewを追加
	timerView.AddItem(timer.textView, 1, 1, false)
}

func stopTimer(timerName string) {
	if timer, ok := timers[timerName]; ok && timer.isRunning {
		timer.isRunning = false
		close(timer.stopChan)                        // タイマーの停止制御用チャンネルを閉じる
		timer.elapsed += time.Since(timer.startTime) // 経過時間を更新
	}
}

func restartTimer(timerName string) {
	// タイマーが存在し、かつ停止している場合は再スタート
	if timer, ok := timers[timerName]; ok && !timer.isRunning {
		timer.startTime = time.Now().Add(-timer.elapsed) // 過去に開始したことにする
		timer.stopChan = make(chan struct{})
		timer.isRunning = true
		go timer.updateTime()
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

		// コマンドを最初のスペースで分割
		commandArgs := strings.SplitN(command, " ", 2)

		if len(commandArgs) != 2 {
			return event
		}

		cmd, timerName := commandArgs[0], commandArgs[1]

		switch cmd {
		case "start":
			startTimer(timerName, timerView)

		case "stop":
			stopTimer(timerName)

		case "restart":
			restartTimer(timerName)
		}

		// 入力欄をクリア
		commandInputField.SetText("")

		return event
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(commandInputField, 0, 1, true).
		AddItem(timerView, 0, 9, false)

	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}
