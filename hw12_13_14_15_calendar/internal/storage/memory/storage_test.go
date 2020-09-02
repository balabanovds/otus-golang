package memorystorage

import (
	"testing"
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func newTestStorage(start time.Time, num int) storage.Repo {
	return &Storage{
		data: generateTestData(start, num),
	}
}

const (
	layoutISO = "2006-01-02"
)

func TestStorage(t *testing.T) {
	t.Run("create in already busy time", func(t *testing.T) {
		st := newTestStorage(time.Now(), 3)

		_, err := st.Events().Create(newTestEvent(time.Now().Add(10 * time.Second)))

		require.EqualError(t, err, storage.ErrEventExists.Error())
	})

	t.Run("update not found", func(t *testing.T) {
		st := newTestStorage(time.Now(), 0)

		err := st.Events().Update("wrongID", newTestEvent(time.Now().Add(10*time.Second)))

		require.EqualError(t, err, storage.ErrEvent404.Error())
	})

	t.Run("delete", func(t *testing.T) {
		st := newTestStorage(time.Now(), 0)
		ev := newTestEvent(time.Now())
		_, err := st.Events().Create(ev)
		require.NoError(t, err)
		require.Len(t, st.Events().ListForDay(time.Now()), 1)

		st.Events().Delete(ev.ID)
		require.Len(t, st.Events().ListForDay(time.Now()), 0)
	})

	t.Run("list events for day", func(t *testing.T) {
		st := newTestStorage(time.Now(), 5)

		got := st.Events().ListForDay(time.Now())

		require.Len(t, got, 1)
	})

	t.Run("list events for week", func(t *testing.T) {
		pt, err := time.Parse(layoutISO, "2020-08-31")
		require.NoError(t, err)
		st := newTestStorage(pt, 50)

		got := st.Events().ListForWeek(pt.Add(24 * time.Hour))

		require.Len(t, got, 7)
	})

	t.Run("list events for month", func(t *testing.T) {
		pt, err := time.Parse(layoutISO, "2020-12-25")
		require.NoError(t, err)
		st := newTestStorage(pt, 50)

		got := st.Events().ListForMonth(pt)

		require.Len(t, got, 7)
	})
}
