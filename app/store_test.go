package main

import (
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func createTempFile(t *testing.T, data string) (io.ReadWriteSeeker, func()) {
	t.Helper()
	tmpfile, err := ioutil.TempFile("","db")

	if err != nil {
		t.Fatalf("Cannot write temp file at %v", tmpfile)
	}

	tmpfile.Write([]byte(data))

	removefile := func(){
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}
	return tmpfile, removefile
}

func TestFileSystemStore(t *testing.T) {
	initdata := `[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`

	database, cleanDatabase := createTempFile(t, initdata)
	defer cleanDatabase()

	store := NewFileSystemPlayerStore(database)

	t.Run("/league from a Reader", func(t *testing.T) {
		_ = store.GetPlayers()
		var got = store.GetPlayers()
		expected := League{
			{"Cleo", 10},
			{"Chris", 33},
		}

		if !reflect.DeepEqual(expected, got) {
			t.Fatalf("Wrong data returned from FileSystemPlayerStore, expected %v, got %v", expected, got)
		}
	})

	t.Run("test get score from a Reader", func(t *testing.T) {
		player := "Cleo"
		got, _ := store.GetPlayerScore(player)
		expected := 10

		assertScoreEqual(t, expected, got, player)
	})

	t.Run("Update score for non existing player", func(t *testing.T) {
		player := "Andy"
		store.UpdatePlayerScore(player)

		got, _ := store.GetPlayerScore(player)
		expected := 1

		assertScoreEqual(t, expected, got, player)
	})

	t.Run("Update score for existing player", func(t *testing.T) {
		player := "Andy"
		store.UpdatePlayerScore(player)
		store.UpdatePlayerScore(player)

		got, _ := store.GetPlayerScore(player)
		expected := 3

		assertScoreEqual(t, expected, got, player)
	})
}

func assertScoreEqual(t *testing.T, expected int, got int, player string) {
	t.Helper()
	if expected != got {
		t.Fatalf("Wrong score from %s, expected %d, got %d", player, expected, got)
	}
}
