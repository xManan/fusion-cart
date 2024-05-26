package merchant

import (
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/xManan/fusion-cart/internal/app/model"
)

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
		Merchant: AMAZON,
		ItemRef: asin,
		Name: title,
		Price: priceF,
		OldPrice: priceF,
		Url: url,
		ImgUrl: imgUrl,
		Rating: ratingF,
	}, nil
}

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


