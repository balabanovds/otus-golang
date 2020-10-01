// +build integration

package integration

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	cfg "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/cmd/config"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/app"
	internalhttp "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/server/http"
	sqlstorage "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage/sql"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/pkg/utils"
)

var s *suite

type config struct {
	Storage cfg.Storage `koanf:"storage"`
	HTTP    cfg.HTTP    `koanf:"http"`
}

func TestMain(m *testing.M) {
	var c config

	// we use empty filename to get all parameters from environment
	err := cfg.New("").Unmarshal(&c)
	if err != nil {
		log.Fatalln(err)
	}

	strg := sqlstorage.New(c.Storage)
	defer utils.Close(strg)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = strg.Connect(ctx); err != nil {
		log.Fatalln(err)
	}

	calendar := app.New(strg)

	server := internalhttp.New(calendar, c.HTTP)
	defer utils.Close(server)

	go func() {
		if err := server.Start(); err != nil {
			log.Fatalln(err)
		}
	}()

	s = newSuite(c, strg)
	defer utils.Close(s)

	os.Exit(m.Run())
}
