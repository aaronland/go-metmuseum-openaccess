package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	_ "fmt"
	"github.com/aaronland/go-metmuseum-openaccess"
	"github.com/aaronland/go-metmuseum-openaccess/html"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {

	cookie_name := flag.String("cookie-name", "", "...")
	cookie_value := flag.String("cookie-value", "", "...")

	with_archive := flag.String("with-archive", "", "...")

	flag.Parse()

	ctx := context.Background()

	seen := new(sync.Map)
	writer := csv.NewWriter(os.Stdout)

	if *with_archive != "" {

		arch_fh, err := os.Open(*with_archive)

		if err != nil {
			log.Fatal(err)
		}

		defer arch_fh.Close()

		arch_reader := csv.NewReader(arch_fh)

		for {
			row, err := arch_reader.Read()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatal(err)
			}

			writer.Write(row)

			seen.Store(row[0], true)
		}

		writer.Flush()
	}

	ck := &http.Cookie{
		Name:   *cookie_name,
		Value:  *cookie_value,
		Path:   "/",
		Domain: "metmuseum.org",
	}

	reader := bufio.NewReader(os.Stdin)

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

		_, ok := seen.Load(link)

		if ok {
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
