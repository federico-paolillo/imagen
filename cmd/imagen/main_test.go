package main

import (
	"bytes"
	"image/png"
	"os"
	"strings"
	"testing"
)

func TestRunWritesPNG(t *testing.T) {
	t.Chdir(t.TempDir())

	var stdout, stderr bytes.Buffer

	if code := run([]string{"-w", "320", "-h", "200"}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit = %d, stderr: %s", code, stderr.String())
	}

	if got, want := strings.TrimSpace(stdout.String()), "320x200.png"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}

	f, err := os.Open("320x200.png")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()

	cfg, err := png.DecodeConfig(f)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Width != 320 || cfg.Height != 200 {
		t.Fatalf("size = %dx%d, want 320x200", cfg.Width, cfg.Height)
	}
}

func TestRunErrors(t *testing.T) {
	cases := map[string]struct {
		args []string
		code int
	}{
		"missing flags":   {nil, 2},
		"missing height":  {[]string{"-w", "10"}, 2},
		"unknown flag":    {[]string{"-x"}, 2},
		"positional args": {[]string{"-w", "10", "-h", "10", "extra"}, 2},
		"non-numeric":     {[]string{"-w", "abc", "-h", "10"}, 2},
		"negative":        {[]string{"-w", "-5", "-h", "10"}, 1},
		"too large":       {[]string{"-w", "100000", "-h", "10"}, 1},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Chdir(t.TempDir())

			var stdout, stderr bytes.Buffer

			if code := run(tc.args, &stdout, &stderr); code != tc.code {
				t.Fatalf("exit = %d, want %d (stderr: %s)", code, tc.code, stderr.String())
			}

			if stderr.Len() == 0 {
				t.Error("expected message on stderr")
			}

			if entries, _ := os.ReadDir("."); len(entries) != 0 {
				t.Errorf("unexpected files written: %v", entries)
			}
		})
	}
}

func TestRunHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer

	if code := run([]string{"-help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}

	if !strings.Contains(stderr.String(), "Usage: imagen") {
		t.Errorf("usage missing: %q", stderr.String())
	}
}
