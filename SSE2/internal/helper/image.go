package helper

import (
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"golang.org/x/image/draw"
)

const ErrFUnsupportedFormat = "unsupported format '%v'"

const (
	ImageFormatPng = "png"
)

func GetImage(filePath string) (img image.Image, err error) {
	var file *os.File
	file, err = os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := file.Close()
		if derr != nil {
			err = multierr.Combine(err, derr)
		}
	}()

	var imageFormat string
	img, imageFormat, err = image.Decode(file)
	if err != nil {
		return nil, err
	}

	switch imageFormat {
	case ImageFormatPng:

	default:
		return nil, errors.Errorf(ErrFUnsupportedFormat, imageFormat)
	}

	return img, nil
}

func GetScaleFactorForMaxSide(
	source image.Rectangle,
	maxSideSize int,
) (scaleCoefficient float64) {
	return float64(maxSideSize) / float64(getRectangleMaximumDimensionSize(source))
}

func getRectangleMaximumDimensionSize(rectangle image.Rectangle) (size int) {
	rs := rectangle.Size()

	if rs.X >= rs.Y {
		return rs.X
	}

	return rs.Y
}

func ScaleImage(
	sourceImage image.Image,
	rectangle image.Rectangle,
) (newImage image.Image) {
	const HighQualitySizeThresholdPx = 360

	var scaler draw.Scaler
	if getRectangleMaximumDimensionSize(rectangle) >= HighQualitySizeThresholdPx {
		scaler = draw.CatmullRom
	} else {
		scaler = draw.ApproxBiLinear
	}

	result := image.NewRGBA(rectangle)

	scaler.Scale(result, rectangle, sourceImage, sourceImage.Bounds(), draw.Over, nil)

	return result
}

func SaveImageAsPngFile(
	img image.Image,
	filePath string,
) (err error) {
	var file *os.File
	file, err = os.Create(filePath)
	if err != nil {
		return err
	}

	defer func() {
		derr := file.Close()
		if derr != nil {
			err = multierr.Combine(err, derr)
		}
	}()

	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}

func getFileNameWithoutExtension(fileName string) (fileBaseName string) {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

func AppendSuffixToFileBaseName(
	fileName string,
	suffix string,
) (newFileName string) {
	const (
		SuffixSeparator = "_"
	)

	baseName := getFileNameWithoutExtension(fileName)

	return baseName + SuffixSeparator + suffix + filepath.Ext(fileName)
}
