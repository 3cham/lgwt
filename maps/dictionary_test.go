package maps

import "testing"

func TestSearch(t *testing.T) {
	t.Run("Map lookup value for key", func(t *testing.T) {
		dict := Dictionary{"test": "this is just a test"}

		got, err := dict.Search("test")
		expected := "this is just a test"

		if err != nil {
			t.Fatal("Error not expected")
		}
		assertCorrectMessage(t, got, expected)
	})

	t.Run("Map lookup value for unknown key", func(t *testing.T) {
		dict := Dictionary{"test": "this is just a test"}

		_, err := dict.Search("unknown test")
		if err == nil {
			t.Fatal("Error expected")
		}

		expected := "could not find the key you provide"
		assertCorrectMessage(t, err.Error(), expected)
	})

	t.Run("Map lookup value for new added key", func(t *testing.T) {
		dict := Dictionary{}
		dict.Add("test", "this is just a test")

		got, err := dict.Search("test")
		expected := "this is just a test"

		if err != nil {
			t.Fatal("Error not expected")
		}
		assertCorrectMessage(t, got, expected)
	})

	t.Run("Map should not overwrite existing key", func(t *testing.T) {
		dict := Dictionary{}
		dict.Add("test", "this is just a test")
		err := dict.Add("test", "this is just another test")

		if err == nil {
			t.Fatal("Error expected")
		}
		assertCorrectMessage(t, err.Error(), ErrKeyAlreadyExists.Error())

		expected := "this is just a test"
		got, _ := dict.Search("test")
		assertCorrectMessage(t, got, expected)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Update should change key value", func(t *testing.T) {
		dict := Dictionary{}
		dict.Add("key", "value")
		dict.Update("key", "new value")

		got, _ := dict.Search("key")
		expected := "new value"

		assertCorrectMessage(t, got, expected)
	})

	t.Run("Update existing keys should throw error", func(t *testing.T) {
		dict := Dictionary{}
		err := dict.Update("key", "value")

		if err == nil {
			t.Fatalf("update non existing key, error expected")
		}

		assertCorrectMessage(t, err.Error(), ErrUpdateNonExistingKey.Error())
	})
}

func TestDelete(t *testing.T) {
	t.Run("Delete existing key", func(t *testing.T) {
		dict := Dictionary{}
		dict.Add("key", "value")
		dict.Delete("key")

		_, err := dict.Search("key")
		if err == nil {
			t.Fatalf("Delete should remove key from dictionary")
		}
	})

	t.Run("Delete non existing key should throw error", func(t *testing.T) {
		dict := Dictionary{}
		err := dict.Delete("key")
		if err == nil {
			t.Fatalf("Delete non existing key, error expected")
		}
		assertCorrectMessage(t, err.Error(), ErrDeleteNonExistingKey.Error())
	})
}

func assertCorrectMessage(t *testing.T, got string, expected string) {
	t.Helper()
	if got != expected {
		t.Errorf("expected %q but got %q", expected, got)
	}
}
