package assets

import (
	"fmt"

	"github.com/gobuffalo/packr"
)

// const (
// 	GITHUB_ICON_NORMAL = "./"
// )

var box = packr.NewBox("icons")

var (
	GitHubIconNormal      = extract("icon-normal.png")
	GitHubIconGreen       = extract("icon-green.png")
	GitHubIconTenderGreen = extract("icon-tender-green.png")
	GitHubIconYellow      = extract("icon-yellow.png")
	GitHubIconOrange      = extract("icon-orange.png")
	GitHubIconRed         = extract("icon-red.png")
)

func extract(name string) []byte {
	return unwrap(box.Find(name))
}

func unwrap(data []byte, err error) []byte {
	if err != nil {
		panic(fmt.Sprintf("Failed to extract embeded assets: %s", err))
	}
	return data
}
