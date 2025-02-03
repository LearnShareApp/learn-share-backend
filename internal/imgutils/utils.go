package imgutils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strings"
)

const MaxSize = 3 << 20

func DecodeImage(base64Img string) ([]byte, error) {
	imgBytes, err := base64.StdEncoding.DecodeString(base64Img)
	if err != nil {
		return nil, fmt.Errorf("invalid image format")
	}

	// check for weight
	if len(imgBytes) > MaxSize {
		return nil, fmt.Errorf("file too large")
	}

	// check is MIME-type (image)
	mimeType := http.DetectContentType(imgBytes)
	if !strings.HasPrefix(mimeType, "image/") {
		return nil, fmt.Errorf("file is not an image")
	}

	return imgBytes, nil
}

func GetImageDimensions(data []byte) (width, height int, err error) {
	reader := bytes.NewReader(data)

	// Read only Metadata
	config, _, err := image.DecodeConfig(reader)
	if err != nil {
		return 0, 0, err
	}

	return config.Width, config.Height, nil
}

func CheckDimension(xRatio, yRatio, width, height int) error {
	ratio := float64(width) / float64(height)

	if ratio != (float64(xRatio) / float64(yRatio)) {
		return fmt.Errorf("invalid image ratiom, must be %dx%d ", xRatio, yRatio)
	}
	return nil
}
