package main

import (
	"reflect"
	"testing"
)

type Person struct {
	Name   string
	Friend Friend
}
type Friend struct {
	Name string
	Age  int
}
type Connection map[string]Person

func TestWalk(t *testing.T) {
	cases := []struct {
		Name          string
		Input         interface{}
		ExpectedCalls []string
	}{
		{
			"struct with 1 string field",
			struct {
				Name string
			}{"Chris"},
			[]string{"Chris"},
		},
		{
			"struct with 2 string field",
			struct {
				Name    string
				Surname string
			}{"Chris", "hansel man"},
			[]string{"Chris", "hansel man"},
		},
		{
			"struct with one non-string field",
			Friend{"Chris", 1},
			[]string{"Chris"},
		},
		{
			"struct with nested struct field",
			Person{"Chris", Friend{"matt", 29}},
			[]string{"Chris", "matt"},
		},
		{
			"pointer to struct",
			&Person{"Chris", Friend{"matt", 29}},
			[]string{"Chris", "matt"},
		},
		{
			"Slices",
			[]Friend{
				{"Laura", 33},
				{"Paul", 30},
			},
			[]string{"Laura", "Paul"},
		},
		{
			"Arrays",
			[3]Friend{
				{"Laura", 33},
				{"Paul", 30},
				{"James", 28},
			},
			[]string{"Laura", "Paul", "James"},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			var got []string
			Walk(test.Input, func(input string) {
				got = append(got, input)
			})
			if !reflect.DeepEqual(got, test.ExpectedCalls) {
				t.Errorf("got %v, want %v", got, test.ExpectedCalls)
			}
		})
	}

	t.Run("With Maps", func(t *testing.T) {
		amap := Connection{
			"Chris": Person{"Chris", Friend{"Paul", 30}},
			"Laura": Person{"Laura", Friend{"James", 28}},
		}

		var got []string
		Walk(amap, func(input string) {
			got = append(got, input)
		})

		want := []string{"Chris", "Paul", "Laura", "James"}
		for _, name := range want {
			assertContains(t, got, name)
		}
	})
	t.Run("With Channels", func(t *testing.T) {
		achannel := make(chan Person)

		go func() {
			achannel <- Person{"Chris", Friend{"Paul", 30}}
			achannel <- Person{"Laura", Friend{"James", 33}}
			close(achannel)

		}()

		var got []string
		want := []string{"Chris", "Paul", "Laura", "James"}
		Walk(achannel, func(input string) {
			got = append(got, input)
		})
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("With Func", func(t *testing.T) {
		aFunction := func() (Person, Friend) {
			return Person{"Chris", Friend{"Paul", 30}}, Friend{"James", 33}
		}

		var got []string
		want := []string{"Chris", "Paul", "James"}
		Walk(aFunction, func(input string) {
			got = append(got, input)
		})
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

}

func assertContains(t testing.TB, haystack []string, needle string) {
	t.Helper()
	contains := false
	for _, item := range haystack {
		if item == needle {
			contains = true
		}
	}
	if !contains {
		t.Errorf("expected %+v to contain %q but it didn't", haystack, needle)
	}
}
