package reflection

import (
	"reflect"
	"testing"
)

type Person struct {
	Name string
	Att  Profile
}

type Profile struct {
	Age     int
	Address string
}

func TestWalk(t *testing.T) {
	cases := []struct {
		Name          string
		Input         interface{}
		ExpectedCalls []string
	}{
		{"struct with one string field",
			struct {
				Name string
			}{"Tung"},
			[]string{"Tung"},
		},
		{
			"struct with two string field",
			struct {
				FirstName string
				LastName  string
			}{"Tung", "Dang"},
			[]string{"Tung", "Dang"},
		},
		{
			"struct with non-string field",
			struct {
				Name string
				Age  int
			}{"Tung", 32},
			[]string{"Tung"},
		},
		{
			"struct with nested struct field",
			struct {
				Name string
				Att  Profile
			}{"Tung", Profile{32, "Hamburg"}},
			[]string{"Tung", "Hamburg"},
		},
		{
			"Pointer to things",
			&Person{"Tung", Profile{32, "Hamburg"}},
			[]string{"Tung", "Hamburg"},
		},
		{
			"slice of struct",
			[]Profile{
				{18, "Hanoi"},
				{29, "Hamburg"},
			},
			[]string{"Hanoi", "Hamburg"},
		},
		{
			"arrays",
			[2]Profile{
				{18, "Hanoi"},
				{29, "Hamburg"},
			},
			[]string{"Hanoi", "Hamburg"},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			var got = []string{}
			walk(test.Input, func(s string) {
				got = append(got, s)
			})
			if !reflect.DeepEqual(got, test.ExpectedCalls) {
				t.Fatalf("got %#v, expected %#v", got, test.ExpectedCalls)
			}
		})
	}

	t.Run("with maps", func(t *testing.T) {
		test := map[string]string{
			"key":  "value",
			"key1": "value1",
		}
		expected := []string{"value", "value1"}
		got := []string{}

		walk(test, func(s string) {
			got = append(got, s)
		})

		assertContains := func(t *testing.T, arr []string, value string) {
			t.Helper()
			contains := false
			for i := 0; i < len(arr); i++ {
				if arr[i] == value {
					contains = true
					break
				}
			}
			if !contains {
				t.Fatalf("%v does not contain value %s", arr, value)
			}
		}
		for _, val := range got {
			assertContains(t, expected, val)
		}
	})

	t.Run("with channel", func(t *testing.T) {
		aChannel := make(chan Profile)

		go func() {
			aChannel <- Profile{33, "Hamburg"}
			aChannel <- Profile{34, "Berlin"}
			close(aChannel)
		}()

		got := []string{}

		walk(aChannel, func(s string) {
			got = append(got, s)
		})

		expected := []string{"Hamburg", "Berlin"}

		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("Wrong result, expected %v, got %v", expected, got)
		}
	})

	t.Run("with func", func(t *testing.T) {
		func_name := func() []Profile {
			return []Profile{{32, "Hello"}, {20, "World"}}
		}
		got := []string{}

		walk(func_name, func(s string) {
			got = append(got, s)
		})

		expected := []string{"Hello", "World"}

		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("Wrong result, expected %v, got %v", expected, got)
		}
	})
}
