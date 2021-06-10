package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"tilank/utils/logger"
	"tilank/utils/mjwt"
	"tilank/utils/rest_err"
)

const (
	jpgExtension  = ".jpg"
	pngExtension  = ".png"
	jpegExtension = ".jpeg"
)

// saveImage return path to save in db
func saveImage(c *fiber.Ctx, claims mjwt.CustomClaim, folder string, imageName string, needThumbnail bool) (string, resterr.APIError) {
	file, err := c.FormFile("image")
	if err != nil {
		apiErr := resterr.NewAPIError("File gagal di upload", http.StatusBadRequest, "bad_request", []interface{}{err.Error()})
		logger.Info(fmt.Sprintf("u: %s | formfile | %s", claims.Name, err.Error()))
		return "", apiErr
	}

	fileName := file.Filename
	fileExtension := strings.ToLower(filepath.Ext(fileName))
	if !(fileExtension == jpgExtension || fileExtension == pngExtension || fileExtension == jpegExtension) {
		apiErr := resterr.NewBadRequestError("Ektensi file tidak di support")
		logger.Info(fmt.Sprintf("u: %s | validate | %s", claims.Name, apiErr.Error()))
		return "", apiErr
	}

	if file.Size > 2*1024*1024 { // 1 MB
		apiErr := resterr.NewBadRequestError("Ukuran file tidak dapat melebihi 2MB")
		logger.Info(fmt.Sprintf("u: %s | validate | %s", claims.Name, apiErr.Error()))
		return "", apiErr
	}

	// rename image
	// path := filepath.Join("static", "image", folder, imageName + fileExtension)
	// pathInDB := filepath.Join("image", folder, imageName + fileExtension)
	path := fmt.Sprintf("static/image/%s/%s", folder, imageName+fileExtension)
	pathInDB := fmt.Sprintf("image/%s/%s", folder, imageName+fileExtension)

	err = c.SaveFile(file, path)
	if err != nil {
		logger.Error(fmt.Sprintf("%s gagal mengupload file", claims.Name), err)
		apiErr := resterr.NewInternalServerError("File gagal di upload", err)
		return "", apiErr
	}

	// generate thumbnail
	if needThumbnail {
		go func() {
			err = generateThumbnail(path, imageName, fileExtension, folder)
			if err != nil {
				logger.Error(fmt.Sprintf("%s gagal menggenerate thumbnail file", claims.Name), err)
			}
		}()
	}

	return pathInDB, nil
}

// generateThumbnail cukup lambat sehingga saya taruh di goroutine
func generateThumbnail(path string, fileName string, fileExtension, folder string) error {
	// open image
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	var img image.Image
	var err2 error

	// decode file into image.Image
	if fileExtension == pngExtension {
		img, err2 = png.Decode(file)
	} else {
		img, err2 = jpeg.Decode(file)
	}

	if err2 != nil {
		return err2
	}

	// resize to width 200 using Lanczos resampling
	m := resize.Thumbnail(300, 300, img, resize.Lanczos3)

	thumbnailPath := fmt.Sprintf("static/image/%s/%s", folder, "thumb_"+fileName+fileExtension)
	out, err := os.Create(thumbnailPath)
	if err != nil {
		return err
	}

	defer out.Close()

	// write new image to file
	if fileExtension == pngExtension {
		err2 = png.Encode(out, m)
	} else {
		err2 = jpeg.Encode(out, m, nil)
	}

	if err2 != nil {
		return err2
	}

	return nil
}

// merubah string masukan ke int , jika error mereturn 0
func stringToInt(queryString string) int {
	number, err := strconv.Atoi(queryString)
	if err != nil {
		return 0
	}
	return number
}
