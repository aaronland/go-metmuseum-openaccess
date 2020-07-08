package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	_ "fmt"
	"net/http"
	"os"
	"github.com/aaronland/go-metmuseum-openaccess"
	"github.com/aaronland/go-metmuseum-openaccess/html"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
)

func main() {

	cookie_name := flag.String("cookie-name", "", "...")
	cookie_value := flag.String("cookie-value", "", "...")
	
	flag.Parse()

	ctx := context.Background()
	
	ck := &http.Cookie{
		Name:   *cookie_name,
		Value:  *cookie_value,
		Path:   "/",
		Domain: "metmuseum.org",
	}

	reader := bufio.NewReader(os.Stdin)
	writer := csv.NewWriter(os.Stdout)
	
	for {

		select {
		case <-ctx.Done():
			break
		default:
			// pass
		}

		body, err := reader.ReadBytes('\n')

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Failed to read bytes, %v", err)
		}

		body = bytes.TrimSpace(body)

		var rec *openaccess.OpenAccessRecord

		err = json.Unmarshal(body, &rec)

		if err != nil {
			log.Println("Failed to unmarshal OpenAccess record, %v", err)
			continue
		}

		if !rec.IsPublicDomain {
			continue
		}

		link := rec.LinkResource
		
		if link == "" {
			continue
		}

		im, err := html.ExtractImageURLsFromLink(ctx, link, ck)

		if err != nil {
			log.Println(link, err)
			continue
		}

		row := []string{
			link,
			"",
			"",
		}

		if im.Main != nil {
			row[1] = im.Main.String()
		}

		if im.Download != nil {
			row[2] = im.Download.String()
		}
		
		writer.Write(row)
	}

	writer.Flush()

	err := writer.Error()

	if err != nil {
		log.Fatal(err)
	}
}
