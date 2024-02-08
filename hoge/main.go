package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// テキストビューを作成します。
	textView := tview.NewTextView()
	textView.SetDynamicColors(true).
		SetWrap(true).
		SetBorder(true).
		SetTitle("Current Time")

	var ticker *time.Ticker
	var stopChan chan bool // タイマーの停止制御用

	startTicker := func() {
		if ticker != nil {
			ticker.Stop() // 既存のTickerがあれば停止します。
		}
		ticker = time.NewTicker(1 * time.Second)
		stopChan = make(chan bool) // ストップチャネルを再作成

		go func() {
			for {
				select {
				case t := <-ticker.C:
					app.QueueUpdateDraw(func() {
						fmt.Fprintf(textView, "Current time: %s\n", t.Format(time.RFC1123))
					})
				case <-stopChan:
					return // ストップチャネルがクローズされたらループを終了
				}
			}
		}()
	}

	stopTicker := func() {
		if ticker != nil {
			ticker.Stop()   // Tickerを停止します。
			close(stopChan) // ストップチャネルをクローズしてゴルーチンを終了させる
		}
	}

	// 最初にタイマーをスタートします。
	startTicker()

	// アプリケーションを終了するためのショートカットキーを設定します。
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 's': // 's'キーでタイマーを停止します。
			stopTicker()
		case 'r': // 'r'キーでタイマーを再スタートします。
			startTicker()
		}
		if event.Key() == tcell.KeyEscape {
			stopTicker() // Tickerを停止します。
			app.Stop()   // アプリケーションを終了します。
		}
		return event
	})

	// テキストビューをルートとしてアプリケーションを実行します。
	if err := app.SetRoot(textView, true).Run(); err != nil {
		panic(err)
	}
}
