package merchant

import (
	"net/url"
	"regexp"

	"github.com/xManan/fusion-cart/internal/app/models"
)

func ExtractItemRefFromFlipkartUrl(urlStruct *url.URL) (string, error) {
	path := urlStruct.EscapedPath()
	re := regexp.MustCompile("p/(itm[a-z0-9]{13})")
	matches := re.FindStringSubmatch(path)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", ErrItemRefNotFound
}

func FetchItemFromFlipkart(itemRef string) (models.Item, error) {
	return models.Item{}, nil
}
