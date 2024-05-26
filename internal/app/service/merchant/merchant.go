package merchant

import (
	"errors"
	"net/url"
	"os"
	"strings"

	"github.com/xManan/fusion-cart/internal/app/model"
)

const (
	AMAZON   = "amazon"
	FLIPKART = "flipkart"
)

var ErrMerchantNotSupported = errors.New("merchant not supported")
var ErrMerchantNotActive = errors.New("merchant not active")
var ErrItemRefNotFound = errors.New("item reference not found")


func GetItemFromUrl(urlStruct url.URL) (model.Item, error) {
	merchant := GetMerchantFromUrl(urlStruct)
	if !IsActive(merchant) {
		return model.Item{}, ErrMerchantNotActive
	}

	itemRef, err := GetItemRefFromUrl(urlStruct, merchant)
	if err != nil {
		return model.Item{}, err
	}

	var item model.Item

	switch merchant {
	case AMAZON:
		item, err = GetItemFromAmazon(itemRef)
	case FLIPKART:
		item, err = GetItemFromFlipkart(itemRef)
	default:
		return model.Item{}, ErrMerchantNotSupported
	}

	return item, err
}

func GetMerchantFromUrl(urlStruct url.URL) string {
	host := urlStruct.Hostname()
	host = strings.TrimPrefix(host, "www.")
	hostSplit := strings.Split(host, ".")
	if len(hostSplit) > 0 {
		return hostSplit[0]
	}
	return ""
}

func GetItemRefFromUrl(urlStruct url.URL, merchant string) (string, error) {
	switch merchant {
	case AMAZON:
		return GetItemRefFromAmazonUrl(urlStruct)
	case FLIPKART:
		return GetItemRefFromFlipkartUrl(urlStruct)
	default:
		return "", ErrMerchantNotSupported
	}
}

func IsSupported(vendor string) bool {
	return vendor == AMAZON || vendor == FLIPKART
}

func IsActive(vendor string) bool {
	switch(vendor) {
	case AMAZON:
		return os.Getenv("AMAZON_ENABLED") == "true"
	case FLIPKART:
		return os.Getenv("FLIPKART_ENABLED") == "true"
	default:
		return false
	}
}

