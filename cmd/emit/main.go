package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	_ "fmt"
	"github.com/aaronland/go-metmuseum-openaccess"
	"github.com/tidwall/pretty"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
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

	bucket_uri := flag.String("bucket-uri", "", "A valid GoCloud bucket file:// URI.")
	objects_csv := flag.String("objects-csv", "MetObjects.csv", "The path the MetObjects.csv file.")

	format := flag.Bool("format", false, "Format JSON output for each record.")
	stdout := flag.Bool("stdout", true, "Emit to STDOUT.")
	null := flag.Bool("null", false, "Emit to /dev/null")

	as_json := flag.Bool("json", false, "Emit a JSON list.")
	// as_oembed := flag.Bool("oembed", false, "Emit results as OEmbed records")

	flag.Parse()

	ctx := context.Background()

	bucket, err := blob.OpenBucket(ctx, *bucket_uri)

	if err != nil {
		log.Fatalf("Failed to open bucket, %v", err)
	}

	defer bucket.Close()

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

		body, err = json.Marshal(rec)

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
