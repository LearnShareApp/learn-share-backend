package imgutils

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
)

func GetImageDimensions(data []byte) (width, height int, err error) {
	reader := bytes.NewReader(data)

	// Read only Metadata
	config, _, err := image.DecodeConfig(reader)
	if err != nil {
		return 0, 0, err
	}

	return config.Width, config.Height, nil
}
