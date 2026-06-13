package handler

import (
	"image"
	"image/color"
	"strings"
	"testing"
)

func TestRandomString(t *testing.T) {
	s := randomString(12)
	if len(s) != 12 {
		t.Fatalf("expected length 12, got %d", len(s))
	}
	// Two calls should (almost certainly) differ
	if randomString(12) == s {
		t.Fatal("two random strings were identical")
	}
}

func TestRandomString_AllowedChars(t *testing.T) {
	const allowed = "abcdefghijklmnopqrstuvwxyz0123456789"
	s := randomString(100)
	for _, r := range s {
		if !strings.ContainsRune(allowed, r) {
			t.Fatalf("unexpected char %q in %q", r, s)
		}
	}
}

func TestExtractDominantColor_Solid(t *testing.T) {
	// 4x4 solid red image
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}
	got := extractDominantColor(img)
	if got != "#ff0000" {
		t.Fatalf("expected #ff0000, got %q", got)
	}
}

func TestExtractDominantColor_ZeroSize(t *testing.T) {
	// Degenerate image with zero-width bounds → count 0 → fallback
	img := image.NewRGBA(image.Rect(0, 0, 0, 0))
	got := extractDominantColor(img)
	if got != "#888888" {
		t.Fatalf("expected fallback #888888, got %q", got)
	}
}

func TestExtractDominantColor_Format(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			img.Set(x, y, color.RGBA{R: 0, G: 0, B: 128, A: 255})
		}
	}
	got := extractDominantColor(img)
	if !strings.HasPrefix(got, "#") || len(got) != 7 {
		t.Fatalf("expected #rrggbb format, got %q", got)
	}
	if got != "#000080" {
		t.Fatalf("expected #000080, got %q", got)
	}
}
