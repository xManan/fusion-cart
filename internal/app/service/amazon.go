package service

import (
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/xManan/fusion-cart/internal/app/model"
)

const (
	AMAZON   = "amazon"
	FLIPKART = "flipkart"
)

var ErrVendorNotSupported = errors.New("vendor not supported")
var ErrItemRefNotFound = errors.New("item reference not found in url")

func GetItemPageFromAmazon(itemRef string) (string, error) {
	url := ConstructAmazonItemUrl(itemRef)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:123.0) Gecko/20100101 Firefox/123.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("DNT", "1")
	req.Header.Set("Sec-GPC", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")

	client := http.DefaultClient
	res, err := client.Do(req)

	if err != nil {
		return "", err
	}

	reader := res.Body
	if res.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(res.Body)
		if err != nil {
			return "", err
		}
	}
	defer reader.Close()

	body, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(body), nil
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

func GetVendorFromUrl(urlStruct url.URL) string {
	host := urlStruct.Hostname()
	host = strings.TrimPrefix(host, "www.")
	hostSplit := strings.Split(host, ".")
	if len(hostSplit) > 0 {
		return hostSplit[0]
	}
	return ""
}

func GetItemRefFromUrl(urlStruct url.URL) (string, error) {
	vendor := GetVendorFromUrl(urlStruct)
	switch vendor {
	case AMAZON:
		return GetItemRefFromAmazonUrl(urlStruct)
	case FLIPKART:
		return GetItemRefFromFlipkartUrl(urlStruct)
	default:
		return "", ErrVendorNotSupported
	}
}

func GetItemRefFromAmazonUrl(urlStruct url.URL) (string, error) {
	path := urlStruct.EscapedPath()
	re := regexp.MustCompile("dp/([A-Z0-9]{10})")
	matches := re.FindStringSubmatch(path)
	if len(matches) > 1 {
		return matches[1], nil
	}
	re = regexp.MustCompile("gp/product/([A-Z0-9]{10})")
	matches = re.FindStringSubmatch(path)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", ErrItemRefNotFound
}

func GetItemRefFromFlipkartUrl(urlStruct url.URL) (string, error) {
	path := urlStruct.EscapedPath()
	re := regexp.MustCompile("p/(itm[a-z0-9]{13})")
	matches := re.FindStringSubmatch(path)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", ErrItemRefNotFound
}

func GetItemFromUrl(urlStruct url.URL) (model.Item, error) {
	vendor := GetVendorFromUrl(urlStruct)

	itemRef, err := GetItemRefFromUrl(urlStruct)
	if err != nil {
		return model.Item{}, err
	}

	var item model.Item
	switch vendor {
	case AMAZON:
		item, err = GetItemFromAmazon(itemRef)
	case FLIPKART:
		item, err = GetItemFromFlipkart(itemRef)
	default:
		return model.Item{}, ErrVendorNotSupported
	}
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	return item, err
}

func ConstructAmazonItemUrl(itemRef string) string {
	return "https://amazon.in/dp/" + itemRef
}

func GetItemFromAmazon(itemRef string) (model.Item, error) {
	url := ConstructAmazonItemUrl(itemRef)
	htmlPage, err := GetItemPageFromAmazon(itemRef)
	if err != nil {
		return model.Item{}, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlPage))
	if err != nil {
		return model.Item{}, err
	}
	title := strings.TrimSpace(doc.Find("#productTitle").First().Text())
	asin, exits := doc.Find("#ASIN").First().Attr("value")
	if !exits {
		return model.Item{}, errors.New("ASIN not found!")
	}
	asin = strings.TrimSpace(asin)
	var priceF float64
	price, exits := doc.Find("#priceValue").First().Attr("value")
	if exits {
		priceF, err = strconv.ParseFloat(price, 64)
	}
	imgUrl, exits := doc.Find("#imgTagWrapperId img").First().Attr("src")
	imgUrl = strings.TrimSpace(imgUrl)

	rating := strings.TrimSpace(doc.Find("#acrPopover a span").First().Text())
	ratingF, _ := strconv.ParseFloat(rating, 64)

	return model.Item{ 
		Vendor: AMAZON,
		ItemRef: asin,
		Name: title,
		Price: priceF,
		OldPrice: priceF,
		Url: url,
		ImgUrl: imgUrl,
		Rating: ratingF,
	}, nil
}


func GetItemFromFlipkart(itemRef string) (model.Item, error) {
	return model.Item{}, nil
}

