package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

type moment struct {
	current time.Time
	exact   time.Time
}

func newMoment() moment {
	t, err := ntp.Time("ntp2.stratum2.ru")
	if err != nil {
		log.Fatal(err)
	}

	return moment{
		current: time.Now(),
		exact:   t,
	}
}

func (m moment) String() string {
	return fmt.Sprintf("current time: %v\nexact time: %v",
		m.current.Round(time.Second), m.exact.Round(time.Second))
}

func main() {
	m := newMoment()
	fmt.Println(m)
}
