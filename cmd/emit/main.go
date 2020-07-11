package main

import (
	"bufio"
	"bytes"
	"compress/bzip2"
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aaronland/go-metmuseum-openaccess"
	"github.com/tidwall/pretty"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

func assignBoolean(row map[string]interface{}, key string) error {

	bool_v, err := strconv.ParseBool(row[key].(string))

	if err != nil {
		return err
	}

	row[key] = bool_v
	return nil
}

func assignBooleanValues(row map[string]interface{}) error {

	keys := []string{
		"Is Public Domain",
		"Is Timeline Work",
		"Is Highlight",
	}

	for _, k := range keys {

		err := assignBoolean(row, k)

		if err != nil {
			return err
		}
	}

	return nil
}

func main() {

	bucket_uri := flag.String("bucket-uri", "", "A valid GoCloud bucket file:// URI where the MetObjects CSV file is stored.")
	objects_csv := flag.String("objects-csv", "MetObjects.csv", "The path for the MetObjects.csv file.")

	images_bucket_uri := flag.String("images-bucket-uri", "", "A valid GoCloud bucket file:// URI where the images lookup CSV file is stored.")
	images_csv := flag.String("images-csv", "images.csv.bz2", "The path for the images.csv file.")

	format := flag.Bool("format", false, "Format JSON output for each record.")
	stdout := flag.Bool("stdout", true, "Emit to STDOUT.")
	null := flag.Bool("null", false, "Emit to /dev/null")

	with_images := flag.Bool("with-images", false, "Append image URLs for public domain records to output.")
	images_is_bzip := flag.Bool("images-is-bzip", true, "The file defined in -images-csv is a bzip2 compressed file.")

	as_json := flag.Bool("json", false, "Emit a JSON list.")
	as_oembed := flag.Bool("oembed", false, "Emit results as OEmbed records")

	oembed_ensure_images := flag.Bool("oembed-ensure-images", true, "Ensure that OEmbed records have an image.")

	flag.Parse()

	if *oembed_ensure_images && !*with_images {
		log.Fatal("-oembed-ensure-images flag is set but -with-images is not")
	}

	ctx := context.Background()

	bucket, err := blob.OpenBucket(ctx, *bucket_uri)

	if err != nil {
		log.Fatalf("Failed to open bucket, %v", err)
	}

	defer bucket.Close()

	var im_lookup *sync.Map

	if *with_images {

		images_bucket, err := blob.OpenBucket(ctx, *images_bucket_uri)

		if err != nil {
			log.Fatalf("Failed to open images bucket, %v", err)
		}

		defer images_bucket.Close()

		im_reader, err := images_bucket.NewReader(ctx, *images_csv, nil)

		if err != nil {
			log.Fatalf("Failed to open images, %v", err)
		}

		defer im_reader.Close()

		br := bufio.NewReader(im_reader)

		if *images_is_bzip {
			cr := bzip2.NewReader(br)
			br = bufio.NewReader(cr)
		}

		r := csv.NewReader(br)

		_, err = r.Read()

		if err != nil {
			log.Fatal(err)
		}

		im_lookup = new(sync.Map)

		for {

			row, err := r.Read()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatal(err)
			}

			im_lookup.Store(row[0], row)
		}

	}

	writers := make([]io.Writer, 0)

	if *stdout {
		writers = append(writers, os.Stdout)
	}

	if *null {
		writers = append(writers, ioutil.Discard)
	}

	wr := io.MultiWriter(writers...)

	fh, err := bucket.NewReader(ctx, *objects_csv, nil)

	if err != nil {
		log.Fatal(err)
	}

	defer fh.Close()

	r := csv.NewReader(fh)

	header, err := r.Read()

	if err != nil {
		log.Fatal(err)
	}

	if *as_json {
		wr.Write([]byte("["))
	}

	counter := int32(0)

	for {

		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		row := make(map[string]interface{})

		for idx, key := range header {
			row[key] = record[idx]
		}

		err = assignBooleanValues(row)

		if err != nil {
			log.Fatal(err)
		}

		body, err := json.Marshal(row)

		if err != nil {
			log.Fatal(err)
		}

		var rec *openaccess.OpenAccessRecord

		err = json.Unmarshal(body, &rec)

		if err != nil {
			log.Fatal(err)
		}

		if *with_images && rec.IsPublicDomain {

			link := rec.LinkResource
			id := strings.Replace(link, openaccess.LINK_RESOURCE_PREFIX, "", 1)

			v, ok := im_lookup.Load(id)

			if ok {
				row := v.([]string)
				rec.MainImage = fmt.Sprintf("%s%s", openaccess.MAIN_IMAGE_PREFIX, row[1])
				rec.DownloadImage = fmt.Sprintf("%s%s", openaccess.DOWNLOAD_IMAGE_PREFIX, row[2])
			}
		}

		var output interface{}
		output = rec

		if *as_oembed {

			if rec.MainImage == "" {
				continue
			}

			author_name := rec.ArtistDisplayName

			if author_name == "" {
				author_name = rec.Department
			}

			oe_rec := openaccess.OEmbedRecord{
				Version:      "1.0",
				Type:         "photo",
				Width:        -1,
				Height:       -1,
				Title:        fmt.Sprintf("%s (%s)", rec.ObjectName, rec.CreditLine),
				URL:          rec.MainImage,
				AuthorName:   author_name,
				AuthorURL:    rec.LinkResource,
				ProviderName: "The Metropolitain Museum of Art",
				ProviderURL:  "https://metmuseum.org/",
				ObjectURI:    fmt.Sprintf("metmuseum://o/%s", rec.ObjectID),
				// DataURL
			}

			output = oe_rec
		}

		body, err = json.Marshal(output)

		if err != nil {
			log.Fatal(err)
		}

		if *format {
			body = pretty.Pretty(body)
		}

		body = bytes.TrimSpace(body)

		new_counter := atomic.AddInt32(&counter, 1)

		if *as_json && new_counter > 1 {
			wr.Write([]byte(","))
		}

		wr.Write(body)
		wr.Write([]byte("\n"))
	}

	if *as_json {
		wr.Write([]byte("]"))
	}

}
