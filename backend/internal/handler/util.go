package handler

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": message,
		"code":  status,
	})
}

// saveAndProcessImage saves an uploaded file and extracts the dominant color.
func saveAndProcessImage(file multipart.File, header *multipart.FileHeader, entityID string, subdir string) (string, string, error) {
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	ext = strings.ToLower(ext)

	filename := fmt.Sprintf("%s%s", entityID+"_"+randomString(8), ext)
	dir := filepath.Join("uploads", subdir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", "", fmt.Errorf("create upload dir: %w", err)
	}

	fullPath := filepath.Join(dir, filename)
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", "", fmt.Errorf("create file: %w", err)
	}
	defer dst.Close()

	// Reset file reader
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", "", fmt.Errorf("seek file: %w", err)
	}

	if _, err := io.Copy(dst, file); err != nil {
		return "", "", fmt.Errorf("save file: %w", err)
	}

	// Extract dominant color
	dominantColor := "#888888"
	if _, err := file.Seek(0, io.SeekStart); err == nil {
		img, _, err := image.Decode(file)
		if err == nil {
			dominantColor = extractDominantColor(img)
		}
	}

	url := fmt.Sprintf("/uploads/%s/%s", subdir, filename)
	return url, dominantColor, nil
}

// extractDominantColor returns the average color of the center region of an image.
func extractDominantColor(img image.Image) string {
	bounds := img.Bounds()
	// Sample center 50% of the image
	x0 := bounds.Min.X + bounds.Dx()/4
	x1 := bounds.Min.X + 3*bounds.Dx()/4
	y0 := bounds.Min.Y + bounds.Dy()/4
	y1 := bounds.Min.Y + 3*bounds.Dy()/4

	var rSum, gSum, bSum uint32
	var count uint32

	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rSum += r >> 8
			gSum += g >> 8
			bSum += b >> 8
			count++
		}
	}

	if count == 0 {
		return "#888888"
	}

	r := rSum / count
	g := gSum / count
	b := bSum / count

	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	rand.Read(b)
	for i := range b {
		b[i] = letters[b[i]%byte(len(letters))]
	}
	return string(b)
}
