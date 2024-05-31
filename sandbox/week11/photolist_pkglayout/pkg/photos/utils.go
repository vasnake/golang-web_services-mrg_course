package photos

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	// "image/jpeg"
	"io"
	"math/rand"
	"os"

	"github.com/disintegration/imaging"
)

var (
	sizes = []int{
		32,
		600,
	}
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

//  не очень эффективно - каждый раз вычитываем файл
func ResizeImage(originalPath string, resizedPath string, size int) error {
	srcImage, err := imaging.Open(originalPath)
	if err != nil {
		return fmt.Errorf("failed to open image: %v\n", err)
	}

	dstImageFill := imaging.Fill(srcImage, size, size, imaging.Center, imaging.Lanczos)

	err = imaging.Save(dstImageFill, resizedPath)
	if err != nil {
		return fmt.Errorf("failed to save image: %v\n", err)
	}

	return nil
}
