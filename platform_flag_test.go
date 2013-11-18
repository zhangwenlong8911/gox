package main

import (
	"flag"
	"reflect"
	"testing"
)

func TestPlatformFlagPlatforms(t *testing.T) {
	cases := []struct {
		OS        []string
		Arch      []string
		Supported []Platform
		Result    []Platform
	}{
		// Building a new list of platforms
		{
			[]string{"foo", "bar"},
			[]string{"baz"},
			[]Platform{
				{"foo", "baz"},
				{"bar", "baz"},
				{"boo", "bop"},
			},
			[]Platform{
				{"foo", "baz"},
				{"bar", "baz"},
			},
		},

		// Skipping platforms
		{
			[]string{"!foo"},
			[]string{},
			[]Platform{
				{"foo", "bar"},
				{"foo", "baz"},
				{"bar", "bar"},
			},
			[]Platform{
				{"bar", "bar"},
			},
		},

		// Specifying only an OS
		{
			[]string{"foo"},
			[]string{},
			[]Platform{
				{"foo", "bar"},
				{"foo", "baz"},
				{"bar", "bar"},
			},
			[]Platform{
				{"foo", "bar"},
				{"foo", "baz"},
			},
		},

		// Building a new list, but with some skips
		{
			[]string{"foo", "bar", "!foo"},
			[]string{"baz"},
			[]Platform{
				{"foo", "bar"},
				{"foo", "baz"},
				{"bar", "baz"},
				{"baz", "bar"},
			},
			[]Platform{
				{"bar", "baz"},
			},
		},

		// Unsupported pairs
		{
			[]string{"foo", "bar"},
			[]string{"baz"},
			[]Platform{
				{"foo", "baz"},
			},
			[]Platform{
				{"foo", "baz"},
			},
		},
	}

	for _, tc := range cases {
		f := PlatformFlag{
			OS:   tc.OS,
			Arch: tc.Arch,
		}

		result := f.Platforms(tc.Supported)
		if !reflect.DeepEqual(result, tc.Result) {
			t.Errorf("input: %#v\nresult: %#v", f, result)
		}
	}
}

func TestPlatformFlagArchFlagValue(t *testing.T) {
	var f PlatformFlag
	val := f.ArchFlagValue()
	if err := val.Set("foo bar"); err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := []string{"foo", "bar"}
	if !reflect.DeepEqual(f.Arch, expected) {
		t.Fatalf("bad: %#v", f.Arch)
	}
}

func TestPlatformFlagOSFlagValue(t *testing.T) {
	var f PlatformFlag
	val := f.OSFlagValue()
	if err := val.Set("foo bar"); err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := []string{"foo", "bar"}
	if !reflect.DeepEqual(f.OS, expected) {
		t.Fatalf("bad: %#v", f.OS)
	}
}

func TestAppendPlatformValue_impl(t *testing.T) {
	var _ flag.Value = new(appendPlatformValue)
}

func TestAppendPlatformValue(t *testing.T) {
	var value appendPlatformValue

	if err := value.Set(""); err != nil {
		t.Fatalf("err: %s", err)
	}

	if len(value) > 0 {
		t.Fatalf("bad: %#v", value)
	}

	if err := value.Set("windows LINUX"); err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := []string{"windows", "linux"}
	if !reflect.DeepEqual([]string(value), expected) {
		t.Fatalf("bad: %#v", value)
	}

	if err := value.Set("darwin"); err != nil {
		t.Fatalf("err: %s", err)
	}

	expected = []string{"windows", "linux", "darwin"}
	if !reflect.DeepEqual([]string(value), expected) {
		t.Fatalf("bad: %#v", value)
	}

	if err := value.Set("darwin"); err != nil {
		t.Fatalf("err: %s", err)
	}

	expected = []string{"windows", "linux", "darwin"}
	if !reflect.DeepEqual([]string(value), expected) {
		t.Fatalf("bad: %#v", value)
	}
}
