package filterit

import (
	"aws-scalable-image-filter/internal/pkg/util"
	"io"
	"io/ioutil"

	"gopkg.in/gographics/imagick.v2/imagick"
)

func applyFilter(imageReader io.ReadCloser, filter int) (blob []byte, err error) {
	// Convert reader to byte[]
	blob, err = ioutil.ReadAll(imageReader)

	if err != nil {
		return nil, err
	}

	imagick.Initialize()
	defer imagick.Terminate() // As per docs

	mw := imagick.NewMagickWand()

	err = mw.ReadImageBlob(blob)

	if err != nil {
		return nil, err
	}

	// Manipulate image based on supplied filter
	switch filter {
	case util.GREYSCALE:
		err = mw.SetColorspace(imagick.COLORSPACE_GRAY)
	case util.INVERT:
		err = mw.NegateImage(false)
	case util.SEPIA:
		err = mw.SepiaToneImage(80)
	}

	if err != nil {
		return nil, err
	}

	err = mw.SetImageFormat("JPEG")

	if err != nil {
		return nil, err
	}

	return mw.GetImageBlob(), nil
}
