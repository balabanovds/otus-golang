// apckage that represents custom progress bar for CLI
package pb

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/superhawk610/terminal"
)

type Bar struct {
	total      int64
	current    int64
	ctrLen     int
	tickStream chan int64
	done       chan struct{}
	output     io.Writer
	doneStr    string
	headStr    string
	emptyStr   string
	barLen     int
}

func New(total int64) *Bar {
	return &Bar{
		total:      total,
		ctrLen:     len(strconv.Itoa(int(total))),
		tickStream: make(chan int64),
		done:       make(chan struct{}),
		output:     os.Stdout,
		doneStr:    "-",
		headStr:    ">",
		emptyStr:   " ",
		barLen:     50,
	}
}

func (b *Bar) SetOutput(w io.Writer) *Bar {
	b.output = w
	return b
}

func (b *Bar) Start() {
	go func(in chan int64) {
		for {
			select {
			case <-b.done:
				return
			case i, ok := <-in:
				if !ok {
					return
				}
				b.current += i
				b.render(b.current)
			}
		}
	}(b.tickStream)
}

func (b *Bar) Add(i int64) {
	b.tickStream <- i
}

func (b *Bar) Finish() {
	b.done <- struct{}{}
	fmt.Println()
}

func (b *Bar) render(current int64) {
	terminal.ClearLine()
	percent := b.calcPercent(current)
	terminal.Overwritef("%10d / %-10d\t%s %3d%%",
		current, b.total, b.renderBar(percent), percent)
}

func (b *Bar) renderBar(percent int) string {
	doneLen := percent / 2.0
	var sb strings.Builder
	sb.WriteRune('[')
	if doneLen > 1 {
		sb.WriteString(strings.Repeat(b.doneStr, doneLen-1))
	}
	sb.WriteString(b.headStr)
	sb.WriteString(strings.Repeat(b.emptyStr, b.barLen-doneLen))
	sb.WriteRune(']')

	return sb.String()
}

func (b *Bar) calcPercent(current int64) int {
	return int(float64(current) / float64(b.total) * 100.0)
}
