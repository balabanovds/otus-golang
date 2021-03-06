package memorystorage

import (
	"testing"
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

const (
	layoutISO = "2006-01-02"
)

func TestStorage(t *testing.T) {
	t.Run("create in already busy time", func(t *testing.T) {
		st := NewTestStorage(time.Now(), 3)

		_, err := st.Events().Create(nil, storage.NewTestEvent(time.Now().Add(10*time.Second)))

		require.EqualError(t, err, storage.ErrEventExists.Error())
	})

	t.Run("update not found", func(t *testing.T) {
		st := NewTestStorage(time.Now(), 0)

		err := st.Events().Update(nil, 999, storage.NewTestEvent(time.Now().Add(10*time.Second)))

		require.EqualError(t, err, storage.ErrEvent404.Error())
	})

	t.Run("delete", func(t *testing.T) {
		st := NewTestStorage(time.Now(), 0)
		ev := storage.NewTestEvent(time.Now())
		ev, err := st.Events().Create(nil, ev)
		require.NoError(t, err)
		require.Len(t, st.Events().ListForDay(nil, time.Now()).List, 1)

		st.Events().Delete(nil, ev.ID)
		list := st.Events().ListForDay(nil, time.Now())

		require.Len(t, list.List, 0)
	})

	t.Run("list events for day", func(t *testing.T) {
		st := NewTestStorage(time.Now(), 5)

		got := st.Events().ListForDay(nil, time.Now())

		require.Len(t, got.List, 1)
	})

	t.Run("list events for week", func(t *testing.T) {
		pt, err := time.Parse(layoutISO, "2020-08-31")
		require.NoError(t, err)
		st := NewTestStorage(pt, 50)

		got := st.Events().ListForWeek(nil, pt.Add(24*time.Hour))

		require.Len(t, got.List, 7)
	})

	t.Run("list events for month", func(t *testing.T) {
		pt, err := time.Parse(layoutISO, "2020-12-25")
		require.NoError(t, err)
		st := NewTestStorage(pt, 50)

		got := st.Events().ListForMonth(nil, pt)

		require.Len(t, got.List, 7)
	})

	t.Run("list events before date", func(t *testing.T) {
		pt, err := time.Parse(layoutISO, "2020-08-31")
		require.NoError(t, err)
		st := NewTestStorage(pt, 50)

		got := st.Events().ListBeforeDate(nil, pt.AddDate(0, 0, 3))

		require.Len(t, got, 3)
	})
}
