package photos

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	// "image/jpeg"
	"bytes"
	"io"
	"math/rand"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
)

var (
	sizes = []int{
		32,
		600,
	}
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

type Putter interface {
	Put(io.ReadSeeker, string, string, uint32) error
}

func MakeThumbnails(storage Putter, source io.ReadSeeker, objectName string, userID uint32) error {
	dst := bytes.NewBuffer(make([]byte, 0, 250*1024*1024))
	for _, size := range sizes {
		source.Seek(0, io.SeekStart)
		dst.Reset()
		err := ResizeImageV2(source, dst, size)
		if err != nil {
			return err
		}
		resizedImg := bytes.NewReader(dst.Bytes())
		err = storage.Put(resizedImg,
			objectName+"_"+strconv.Itoa(size)+".jpg", "image/jpeg",
			userID)
		if err != nil {
			return err
		}
	}
	return nil
}

func ResizeImageV2(source io.Reader, dst io.Writer, size int) error {
	srcImage, err := imaging.Decode(source)
	if err != nil {
		return fmt.Errorf("failed to open image: %w\n", err)
	}

	dstImageFill := imaging.Fill(srcImage, size, size, imaging.Center, imaging.Lanczos)
	err = imaging.Encode(dst, dstImageFill, imaging.JPEG)
	if err != nil {
		return fmt.Errorf("failed to write message: %w\n", err)
	}

	return nil
}

// -----

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

func MakeThumbnails0(realFile, md5Sum string) error {
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

	// imaging.JPEG
	/*
		func Save(img image.Image, filename string, opts ...EncodeOption) (err error) {
			f, err := FormatFromFilename(filename)
			if err != nil {
				return err
			}
			file, err := fs.Create(filename)
			if err != nil {
				return err
			}
			err = Encode(file, img, f, opts...)

	*/
	err = imaging.Save(dstImageFill, resizedPath)
	if err != nil {
		return fmt.Errorf("failed to save image: %v\n", err)
	}

	return nil
}
