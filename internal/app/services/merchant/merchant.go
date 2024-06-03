package merchant

import (
	"errors"
	"net/url"
	"strings"

	"github.com/xManan/fusion-cart/internal/app/models"
)

const (
	AMAZON   = "amazon"
	FLIPKART = "flipkart"
)

var ErrMerchantNotSupported = errors.New("merchant not supported")
var ErrItemRefNotFound = errors.New("item reference not found")

type ProductLink struct {
	Url      *url.URL
	Merchant string
	ItemRef  string
}

func NewProductLink(u *url.URL) (ProductLink, error) {
	var pl ProductLink
	pl.Url = u
	host := u.Hostname()
	host = strings.TrimPrefix(host, "www.")
	hostSplit := strings.Split(host, ".")
	if len(hostSplit) > 0 {
		pl.Merchant = hostSplit[0]
	}
	itemRef, err := ExtractItemRef(u, pl.Merchant)
	if err != nil {
		return pl, err
	}
	pl.ItemRef = itemRef
	return pl, nil
}

func ExtractItemRef(urlStruct *url.URL, merchant string) (string, error) {
	switch merchant {
	case AMAZON:
		return ExtractItemRefFromAmazonUrl(urlStruct)
	case FLIPKART:
		return ExtractItemRefFromFlipkartUrl(urlStruct)
	default:
		return "", ErrMerchantNotSupported
	}
}

func FetchItem(link *ProductLink) (models.Item, error) {
	var item models.Item
	var err error

	switch link.Merchant {
	case AMAZON:
		item, err = FetchItemFromAmazon(link.ItemRef)
	case FLIPKART:
		item, err = FetchItemFromFlipkart(link.ItemRef)
	default:
		return models.Item{}, ErrMerchantNotSupported
	}

	return item, err
}

