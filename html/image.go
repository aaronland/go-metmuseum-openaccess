package html

import (
	"bytes"
	"context"
	gohtml "golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

/*

https://www.metmuseum.org/art/collection/search/7843

<img id="artwork__image" class="artwork__image" src="https://collectionapi.metmuseum.org/api/collection/v1/iiif/7843/48863/main-image" alt="" itemprop="contentUrl" style="max-height: 453px">

<a href="https://images.metmuseum.org/CRDImages/ad/original/DP255720.jpg" class="gtm__download__image" title="Download" download="">

*/

const DOWNLOAD_IMAGE string = "gtm__download__image"
const MAIN_IMAGE string = "artwork__image"

type ImageURLs struct {
	Main     *url.URL
	Download *url.URL
}

func ExtractImageURLsFromLink(ctx context.Context, link string, cookie *http.Cookie) (*ImageURLs, error) {

	select {
	case <-ctx.Done():
		return nil, nil
	default:
		// pass
	}

	u, err := url.Parse(link)

	if err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(nil)

	if err != nil {
		return nil, err
	}

	if cookie != nil {

		var cookies []*http.Cookie
		cookies = append(cookies, cookie)
		jar.SetCookies(u, cookies)
	}

	client := &http.Client{
		Jar: jar,
	}

	rsp, err := client.Get(link)

	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)

	if err != nil {
		return nil, err
	}

	// log.Println(string(body))

	br := bytes.NewReader(body)

	return ExtractImageURLs(ctx, br)
}

func ExtractImageURLs(ctx context.Context, fh io.Reader) (*ImageURLs, error) {

	select {
	case <-ctx.Done():
		return nil, nil
	default:
		// pass
	}

	doc, err := gohtml.Parse(fh)

	if err != nil {
		return nil, err
	}

	im := new(ImageURLs)

	var f func(*gohtml.Node)

	f = func(n *gohtml.Node) {

		if n.Type == gohtml.ElementNode {

			switch n.Data {
			case "a":

				u, err := extractDownloadImageFromAnchorNode(ctx, n)

				if err != nil {
					log.Println(n, err)
				}

				if u != nil {
					im.Download = u
				}

			case "img":

				u, err := extractMainImageFromImgNode(ctx, n)

				if err != nil {
					log.Println(n, err)
				}

				if u != nil {
					im.Main = u
				}

			default:
				// pass
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	return im, nil
}

func extractMainImageFromImgNode(ctx context.Context, n *gohtml.Node) (*url.URL, error) {

	attr_map := make(map[string]string)

	for _, a := range n.Attr {
		attr_map[a.Key] = a.Val
	}

	html_id, ok := attr_map["id"]

	if !ok {
		return nil, nil
	}

	if html_id != MAIN_IMAGE {
		return nil, nil
	}

	src, ok := attr_map["src"]

	if !ok {
		return nil, nil
	}

	return url.Parse(src)
}

func extractDownloadImageFromAnchorNode(ctx context.Context, n *gohtml.Node) (*url.URL, error) {

	attr_map := make(map[string]string)

	for _, a := range n.Attr {
		attr_map[a.Key] = a.Val
	}

	html_class, ok := attr_map["class"]

	if !ok {
		return nil, nil
	}

	if html_class != DOWNLOAD_IMAGE {
		return nil, nil
	}

	href, ok := attr_map["href"]

	if !ok {
		return nil, nil
	}

	return url.Parse(href)
}
