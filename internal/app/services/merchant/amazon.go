package merchant

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

	"github.com/PuerkitoBio/goquery"
	"github.com/xManan/fusion-cart/internal/app/models"
)

const (
	AMAZON_BASE_URL = "https://amazon.in/dp/"
)

func ExtractItemRefFromAmazonUrl(urlStruct *url.URL) (string, error) {
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

func FetchItemFromAmazon(itemRef string) (models.Item, error) {
	url := AMAZON_BASE_URL + itemRef
	htmlPage, err := FetchItemPageFromAmazon(itemRef)

	/*----------------------------------------------------*/
	fn := AMAZON + "_" + itemRef + ".html"                //
	_, err = os.Stat(fn)                                  //
	if os.IsExist(err) {                                  //
		os.Remove(fn)                                     //
	}                                                     //
	fi, _ := os.Create(AMAZON + "_" + itemRef + ".html")  //
	fi.WriteString(htmlPage)                              //
	//----------------------------------------------------//

	if err != nil {
		return models.Item{}, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlPage))
	if err != nil {
		return models.Item{}, err
	}
	title := strings.TrimSpace(doc.Find("#productTitle").First().Text())
	asin, exits := doc.Find("#ASIN").First().Attr("value")
	if !exits {
		return models.Item{}, errors.New("ASIN not found!")
	}
	asin = strings.TrimSpace(asin)
	var priceF float64
	price := strings.ReplaceAll(doc.Find("#corePriceDisplay_desktop_feature_div .a-price-whole").First().Text(), ",", "")
	if exits {
		priceF, err = strconv.ParseFloat(price, 64)
	}
	imgUrl, exits := doc.Find("#imgTagWrapperId img").First().Attr("src")
	imgUrl = strings.TrimSpace(imgUrl)

	rating := strings.TrimSpace(doc.Find("#acrPopover a span").First().Text())
	ratingF, _ := strconv.ParseFloat(rating, 64)

	return models.Item{ 
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

func FetchItemPageFromAmazon(itemRef string) (string, error) {
	url := AMAZON_BASE_URL + itemRef
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("authority", "www.amazon.in")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("accept-language", "en-GB,en;q=0.6")
	req.Header.Set("sec-ch-ua", `"Chromium";v="122", "Not(A:Brand";v="24", "Brave";v="122"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("sec-gpc", "1")
	req.Header.Set("upgrade-insecure-requests", "1")
	// req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.3")

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

