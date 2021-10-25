package utils

import (
	"fmt"
	"github.com/redmask-hb/GoSimplePrint/goPrint"
	"github.com/spf13/cast"
	"time"
)

type LocalBar struct {
	BarCount    int
	Start       int
	Notice      string
	Graph       string
	NoticeColor int
	GoBar       *goPrint.Bar
	SleepTime   int
}

func (l *LocalBar) GenBar() {
	cast.ToInt(If(l.BarCount == 0, 50, l.BarCount))
	l.GoBar = goPrint.NewBar(cast.ToInt(If(l.BarCount == 0, 50, l.BarCount)))
	if l.Notice != "" {
		l.GoBar.SetNotice(l.Notice)
	}

	if l.Graph != "" {
		l.GoBar.SetGraph(l.Graph)
	}

	if l.NoticeColor != 0 {
		l.GoBar.SetNoticeColor(l.NoticeColor)
	}

	if l.SleepTime == 0 {
		l.SleepTime = 50
	}
}

func (l *LocalBar) PrintBar() {
	l.Start++
	l.GoBar.PrintBar(l.Start)
	time.Sleep(time.Duration(l.SleepTime) * time.Millisecond)
}

func (l *LocalBar) EndBar() {
	fmt.Printf(" Done!\n")
}
