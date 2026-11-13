package ui

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"testing"
	"time"

	"fyne.io/fyne/v2"
)

func TestLRUThumbnailCache_CapacityAndEviction(t *testing.T) {
	cache := NewLRUThumbnailCache(3) // Tiny capacity of 3 for testing eviction

	r1 := fyne.NewStaticResource("k1", []byte("data1"))
	r2 := fyne.NewStaticResource("k2", []byte("data2"))
	r3 := fyne.NewStaticResource("k3", []byte("data3"))
	r4 := fyne.NewStaticResource("k4", []byte("data4"))

	cache.Put("k1", r1)
	cache.Put("k2", r2)
	cache.Put("k3", r3)

	// Verify all 3 present
	if _, ok := cache.Get("k1"); !ok {
		t.Fatalf("expected k1 to be present")
	}
	if _, ok := cache.Get("k2"); !ok {
		t.Fatalf("expected k2 to be present")
	}
	if _, ok := cache.Get("k3"); !ok {
		t.Fatalf("expected k3 to be present")
	}

	// Insert 4th item -> k1 was accessed, so oldest unaccessed is k2 or k1
	// Let's touch k1 and k2 so k3 is oldest
	cache.Get("k1")
	cache.Get("k2")
	cache.Put("k4", r4)

	// k3 should be evicted
	if _, ok := cache.Get("k3"); ok {
		t.Fatalf("expected k3 to be evicted after reaching capacity")
	}
	if _, ok := cache.Get("k4"); !ok {
		t.Fatalf("expected k4 to be present")
	}
}

func TestAsyncImageDecodingAndDownscale(t *testing.T) {
	tempDir := t.TempDir()
	testImgPath := filepath.Join(tempDir, "test_thumb.jpg")

	// Generate 400x300 image
	err := GenerateTestDummyImage(testImgPath, 400, 300, color.RGBA{R: 79, G: 70, B: 229, A: 255})
	if err != nil {
		t.Fatalf("failed to generate dummy image: %v", err)
	}
	defer os.Remove(testImgPath)

	done := make(chan fyne.Resource, 1)
	LoadThumbnailAsync(testImgPath, 128, func(res fyne.Resource) {
		done <- res
	})

	select {
	case res := <-done:
		if res == nil {
			t.Fatalf("expected valid downscaled resource")
		}
		if len(res.Content()) == 0 {
			t.Fatalf("expected non-empty thumbnail content")
		}
	case <-time.After(3 * time.Second):
		t.Fatalf("timed out waiting for async thumbnail decoding")
	}

	// Immediate cache hit
	cacheKey := fmt.Sprintf("%s@%d", testImgPath, 128)
	if _, hit := GetThumbnailCache().Get(cacheKey); !hit {
		t.Fatalf("expected image to be cached in LRU cache")
	}
}
