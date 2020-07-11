# go-metmuseum-openaccess

Go package for working with the The Metropolitan Museum of Art's Open Access Initiative data.

## Tools

### emit

Command line too to emit each line of the Open Access CSV document as a JSON record.

For example:

```
$> go run -mod vendor cmd/emit/main.go \
	-bucket-uri file:///usr/local/aaronland/openaccess/ \
	-format
	
{
  "Artist Display Bio": "American, Delaware County, Pennsylvania 1794–1869 Philadelphia, Pennsylvania",
  "Artist Nationality": "American",
  "Subregion": "",
  "Object Wikidata URL": "",
  "Tags AAT URL": "",
  "Object Name": "Coin",
  "Portfolio": "",
  "Artist Alpha Sort": "Longacre, James Barton",
  "Object End Date": "1853",
  "Excavation": "",
  "Metadata Date": "",
  "Object ID": "1",
  "Department": "The American Wing",
  "Country": "",
  "Locus": "",
  "Period": "",
  "County": "",
  "Culture": "",
  "Dynasty": "",
  "Reign": "",
  "Artist Wikidata URL": "",
  "Locale": "",
  "Is Public Domain": false,
  "Title": "One-dollar Liberty Head Coin",
  "Artist Role": "Maker",
  "Artist End Date": "1869      ",
  "Credit Line": "Gift of Heinz L. Stoppelmann, 1979",
  "Repository": "Metropolitan Museum of Art, New York, NY",
  "Object Number": "",
  "Is Timeline Work": false,
  "Artist Begin Date": "1794      ",
  "Dimensions": "Dimensions unavailable",
  "Geography Type": "",
  "Classification": "Metal",
  "Tags": "",
  "AccessionYear": "1979",
  "Artist Display Name": "James Barton Longacre",
  "Artist Suffix": "",
  "Artist ULAN URL": "http://vocab.getty.edu/page/ulan/500011409",
  "Object Date": "1853",
  "City": "",
  "State": "",
  "Rights and Reproduction": "",
  "Is Highlight": false,
  "Artist Prefix": "",
  "Link Resource": "http://www.metmuseum.org/art/collection/search/1",
  "Medium": "Gold",
  "Region": "",
  "River": "",
  "Artist Gender": "",
  "Object Begin Date": "1853"
}
... and so on
```

### images

Command line tool to generate a CSV document mapping Open Access `Link Resource` URLs to their corresponding "main" and "download" image URLs. It is designed to be used in concert with the `emit` tool and any records marked as `Is Public Domain: false` are excluded.

_You should not need to use this tool as its output is bundled in the [data/images.csv.bz2](data/README.md) file._

For example:

```
$> go run -mod vendor cmd/emit/main.go \
	-bucket-uri file:///usr/local/openaccess/ \
	| \
	go run -mod vendor cmd/images/main.go \
	-with-archive data/images.csv \
	-cookie-name {COOKIE_NAME -cookie-value {COOKIE_VALUE} \
	> images.csv
```

The `-cookie-name` and `-cookie-value` parameters are the name and value of a valid `incap_ses_{SUFFIX}` cookie. I have found the easiest way to deal with this is simply to vist the `metmuseum.org` website in a browser and, using the developer tools, copy and paste the relevant cookie data.

Hopefully future releases of the [openaccess data](https://github.com/metmuseum/openaccess) will include image URL information so that this tool won't be necessary anymore.

## See also

* https://github.com/metmuseum/openaccess