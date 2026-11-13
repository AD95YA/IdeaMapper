package ui

import (
	"bytes"
	"container/list"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// LRUThumbnailCache provides bounded in-memory caching of downscaled images
type LRUThumbnailCache struct {
	mu       sync.Mutex
	capacity int
	items    map[string]*list.Element
	evictList *list.List
}

type lruEntry struct {
	key   string
	value fyne.Resource
}

var (
	globalThumbnailCache = NewLRUThumbnailCache(50)
)

// GetThumbnailCache returns the global LRU cache instance
func GetThumbnailCache() *LRUThumbnailCache {
	return globalThumbnailCache
}

// NewLRUThumbnailCache creates an LRU cache with the specified capacity
func NewLRUThumbnailCache(capacity int) *LRUThumbnailCache {
	if capacity <= 0 {
		capacity = 50
	}
	return &LRUThumbnailCache{
		capacity:  capacity,
		items:     make(map[string]*list.Element),
		evictList: list.New(),
	}
}

// Get retrieves a cached thumbnail resource if available
func (c *LRUThumbnailCache) Get(key string) (fyne.Resource, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.items[key]; found {
		c.evictList.MoveToFront(elem)
		return elem.Value.(*lruEntry).value, true
	}
	return nil, false
}

// Put stores a downscaled thumbnail resource into the cache
func (c *LRUThumbnailCache) Put(key string, val fyne.Resource) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.items[key]; found {
		c.evictList.MoveToFront(elem)
		elem.Value.(*lruEntry).value = val
		return
	}

	elem := c.evictList.PushFront(&lruEntry{key: key, value: val})
	c.items[key] = elem

	if c.evictList.Len() > c.capacity {
		oldest := c.evictList.Back()
		if oldest != nil {
			c.evictList.Remove(oldest)
			delete(c.items, oldest.Value.(*lruEntry).key)
		}
	}
}

// Clear flushes all cached thumbnails
func (c *LRUThumbnailCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*list.Element)
	c.evictList.Init()
}

// LoadThumbnailAsync asynchronously loads and downscales an image, returning the resource to callback
func LoadThumbnailAsync(path string, maxDim int, onLoaded func(fyne.Resource)) {
	if path == "" {
		return
	}

	cacheKey := fmt.Sprintf("%s@%d", path, maxDim)
	if res, hit := globalThumbnailCache.Get(cacheKey); hit {
		if onLoaded != nil {
			onLoaded(res)
		}
		return
	}

	go func() {
		res, err := decodeAndDownscale(path, maxDim)
		if err == nil && res != nil {
			globalThumbnailCache.Put(cacheKey, res)
			if onLoaded != nil {
				onLoaded(res)
			}
		}
	}()
}

// CreateAsyncImageWidget creates a canvas.Image that populates its content asynchronously
func CreateAsyncImageWidget(path string, maxDim int) *canvas.Image {
	img := canvas.NewImageFromResource(nil)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(float32(maxDim), float32(maxDim)))

	LoadThumbnailAsync(path, maxDim, func(res fyne.Resource) {
		img.Resource = res
		img.Refresh()
	})

	return img
}

func decodeAndDownscale(path string, maxDim int) (fyne.Resource, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	srcImg, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := srcImg.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	if srcW <= 0 || srcH <= 0 {
		return nil, fmt.Errorf("invalid image bounds")
	}

	// Calculate target dimensions
	targetW := srcW
	targetH := srcH

	if maxDim > 0 && (srcW > maxDim || srcH > maxDim) {
		if srcW > srcH {
			targetW = maxDim
			targetH = (srcH * maxDim) / srcW
		} else {
			targetH = maxDim
			targetW = (srcW * maxDim) / srcH
		}
	}
	if targetW < 1 {
		targetW = 1
	}
	if targetH < 1 {
		targetH = 1
	}

	// Downsample image
	dst := image.NewRGBA(image.Rect(0, 0, targetW, targetH))
	for y := 0; y < targetH; y++ {
		srcY := bounds.Min.Y + (y * srcH) / targetH
		for x := 0; x < targetW; x++ {
			srcX := bounds.Min.X + (x * srcW) / targetW
			dst.Set(x, y, srcImg.At(srcX, srcY))
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, dst); err != nil {
		return nil, err
	}

	return fyne.NewStaticResource(path, buf.Bytes()), nil
}

// GenerateTestDummyImage creates an in-memory test image and saves it to path
func GenerateTestDummyImage(path string, width, height int, col color.Color) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, col)
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return jpeg.Encode(f, img, &jpeg.Options{Quality: 80})
}
