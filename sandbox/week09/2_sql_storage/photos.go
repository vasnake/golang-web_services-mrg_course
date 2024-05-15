package sql_storage

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image/jpeg"
	"io"
	"math/rand"
	"os"

	"github.com/nfnt/resize"
)

// same shit

var (
	sizes       = []uint{80, 160, 320}
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func SaveFile(in io.Reader) (string, error) {
	tmpName := RandStringRunes(32)

	tmpFile := "./images/" + tmpName + ".jpg"
	newFile, err := os.Create(tmpFile)
	if err != nil {
		return "", err
	}

	hasher := md5.New()
	_, err = io.Copy(newFile, io.TeeReader(in, hasher))
	if err != nil {
		return "", err
	}
	newFile.Sync()
	newFile.Close()

	md5Sum := hex.EncodeToString(hasher.Sum(nil))

	realFile := "./images/" + md5Sum + ".jpg"
	err = os.Rename(tmpFile, realFile)
	if err != nil {
		return "", err
	}

	return md5Sum, nil
}

func MakeThumbnails(realFile, md5Sum string) error {
	for _, size := range sizes {
		resizedPath := fmt.Sprintf("./images/%s_%d.jpg", md5Sum, size)
		err := ResizeImage(realFile, resizedPath, size)
		if err != nil {
			return err
		}
	}
	return nil
}

// не очень эффективно - каждый раз вычитываем файл
func ResizeImage(originalPath string, resizedPath string, size uint) error {
	file, err := os.Open(originalPath)
	if err != nil {
		return fmt.Errorf("cant open file %s: %s", originalPath, err)
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		return fmt.Errorf("cant jpeg decode file %s", err)
	}
	file.Close()

	resizeImage := resize.Resize(size, 0, img, resize.Lanczos3)

	out, err := os.Create(resizedPath)
	if err != nil {
		return fmt.Errorf("cant create file %s: %s", resizedPath, err)
	}
	defer out.Close()

	jpeg.Encode(out, resizeImage, nil)

	return nil
}
