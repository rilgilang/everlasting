package validator

import (
	"fmt"
	"net/http"
)

var AllowedMIMEs = map[string]string{
	"image/jpeg": "jpg",
	"image/png":  "png",
	//"image/webp": "webp",
}

// DetectAndValidateImage checks magic bytes and ensures the MIME is allowed
func DetectAndValidateImage(b []byte) (mime, ext string, err error) {
	mime = http.DetectContentType(b)
	if mime == "image/jpg" {
		mime = "image/jpeg"
	}

	ext, ok := AllowedMIMEs[mime]
	if !ok {
		return "", "", fmt.Errorf("unsupported image type: %s", mime)
	}
	return mime, ext, nil
}
