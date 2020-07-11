# go-metmuseum-openaccess

Go package for working with the [The Metropolitan Museum of Art's Open Access Initiative](https://github.com/metmuseum/openaccess) data.

## Tools

To build binary versions of these tools run the `cli` Makefile target. For example:

```
$> make cli
go build -mod vendor -o bin/emit cmd/emit/main.go
go build -mod vendor -o bin/images cmd/images/main.go
```

### emit

Command line too to emit each line of the Open Access CSV document as a JSON record.

```
> ./bin/emit -h
Usage of ./bin/emit:
  -bucket-uri string
    	A valid GoCloud bucket file:// URI where the MetObjects CSV file is stored.
  -format
    	Format JSON output for each record.
  -images-bucket-uri string
    	A valid GoCloud bucket file:// URI where the images lookup CSV file is stored.
  -images-csv string
    	The path for the images.csv file. (default "images.csv.bz2")
  -images-is-bzip
    	The file defined in -images-csv is a bzip2 compressed file. (default true)
  -json
    	Emit a JSON list.
  -null
    	Emit to /dev/null
  -objects-csv string
    	The path for the MetObjects.csv file. (default "MetObjects.csv")
  -oembed
    	Emit results as OEmbed records
  -oembed-ensure-images
    	Ensure that OEmbed records have an image. (default true)
  -stdout
    	Emit to STDOUT. (default true)
  -with-images
    	Append image URLs for public domain records to output.
```

For example:

```
$> bin/emit \
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

#### JSON

By default the `emit` tool outputs line-separated JSON. If you want to output a well-formed JSON array you would enable the `-json` flag. For example:

```
$> bin/emit \
	-bucket-uri file:///usr/local/openaccess 
	-json \
	
   | jq '.[]["Culture"]' \
   | sort \
   | uniq

"Abau people"
"Abelam people"
"Abenaki (?)"
"Abor, Burmese"
"Abron peoples"
"Acehnese"
"Achaemenid (?)"
"Achaemenid or Seleucid"
"Achaemenid"
"Acheen"
"Acheulean"
"Acoma Pueblo"
"Acoma"
"Acoma, Native American"
"Adjora or Aion"
"Adjora or Kopar"
"Admiralty Islands"
"Aduma peoples"
"Aegean"
"Afghan (Nuristan)"
"Afghan (Pashtun)"
"Afghan (Turkmen)"
"Afghan (possibly Hazaras)"
...
"probably façon de Venise, northern European or Venetian"
"probably northern European (probably German)"
"probably provincial British"
"probably south Lowlands; possibly Bohemia or Saxony"
"probaby Senegalese (Fula or Wolof)"
"saddle plate, Chinese or Tibetan; harness fittings, Tibetan"
"southern German or Tyrolese; cranequin probably German or Swiss"
"spearhead, Chinese or Mongolian; case, Tibetan"
"staff, Indian; banner, Mahdist Sudanese"
"unknown (Italian style)"
"unknown"
```

#### Images

By default the Met OpenAccess data does not contain image URLs. It is possible to append those URLS for public domain records by passing the `-with-images` and `-images-bucket-uri` flags. If present the tool with load a lookup table (produced by the `images` tool discussed below) and append `Main Image` and `Download Image` properties to the JSON output.

For example:

```
$> bin/emit -format \
	-with-images \
	-bucket-uri file:///usr/local/openaccess
	-images-bucket-uri file:///usr/local/go-metmuseum-openaccess/data

{
  "Artist Display Bio": "Mexican, active 1607–70",
  "Artist Nationality": "",
  "Subregion": "",
  "Object Wikidata URL": "https://www.wikidata.org/wiki/Q83560129",
  ...
  "Link Resource": "http://www.metmuseum.org/art/collection/search/9728",
  "Medium": "Tin-glazed earthenware",
  "Region": "",
  "River": "",
  "Artist Gender": "",
  "Object Begin Date": "1660",
  "Main Image": "https://collectionapi.metmuseum.org/api/collection/v1/iiif/9728/24914/main-image",
  "Download Image": "https://images.metmuseum.org/CRDImages/ad/original/DP105071.jpg"
}
```

#### OEmbed

It is also possible to emit OEmbed records by passing the `-oembed` flag. By default all records are included but if you want to ensure that all OEmbed records have an image you would also pass in the `-with-images` and `-oembed-with-images` flags. For example:

```
$> bin/emit \
	-format \
	-oembed -oembed-with-images \
	-with-images \
	-bucket-uri file:///usr/local/openaccess \
	-images-bucket-uri file:///usr/local/go-metmuseum-openaccess/data

{
  "version": "1.0",
  "type": "photo",
  "width": -1,
  "height": -1,
  "title": "Coffee spoon (Gift of Mrs. Samuel P. Avery, 1897)",
  "url": "https://collectionapi.metmuseum.org/api/collection/v1/iiif/188051/390495/main-image",
  "author_name": "Benjamin Brewood II",
  "author_url": "http://www.metmuseum.org/art/collection/search/188051",
  "provider_name": "The Metropolitain Museum of Art",
  "provider_url": "https://metmuseum.org/",
  "object_uri": "metmuseum://o/188051"
}
{
  "version": "1.0",
  "type": "photo",
  "width": -1,
  "height": -1,
  "title": "Sugar spoon (Gift of Mrs. Samuel P. Avery, 1897)",
  "url": "https://collectionapi.metmuseum.org/api/collection/v1/iiif/188052/390496/main-image",
  "author_name": "Pierre-Nicolas Sommé",
  "author_url": "http://www.metmuseum.org/art/collection/search/188052",
  "provider_name": "The Metropolitain Museum of Art",
  "provider_url": "https://metmuseum.org/",
  "object_uri": "metmuseum://o/188052"
}
{
  "version": "1.0",
  "type": "photo",
  "width": -1,
  "height": -1,
  "title": "Sugar spoon (Gift of Mrs. Samuel P. Avery, 1897)",
  "url": "https://collectionapi.metmuseum.org/api/collection/v1/iiif/188053/390497/main-image",
  "author_name": "Claude-Auguste Aubry",
  "author_url": "http://www.metmuseum.org/art/collection/search/188053",
  "provider_name": "The Metropolitain Museum of Art",
  "provider_url": "https://metmuseum.org/",
  "object_uri": "metmuseum://o/188053"
}
```

_The `height` and `width` properties are assigned values of "-1" because images dimensions are unavailable at this time._

#### Inline queries

You can also specify inline queries by passing one or more `-query` parameters in the form of:

```
{PATH}={REGULAR EXPRESSION}
```

Paths follow the dot notation syntax used by the [tidwall/gjson](https://github.com/tidwall/gjson) package and regular expressions are any valid [Go language regular expression](https://golang.org/pkg/regexp/). Successful path lookups will be treated as a list of candidates and each candidate's string value will be tested against the regular expression's [MatchString](https://golang.org/pkg/regexp/#Regexp.MatchString) method.

For example:

```
$> bin/emit \
	-bucket-uri file:///usr/local/openaccess \
	-query 'Credit Line=(?:)Eric G\. Carlson' \
   | jq '.["Title"]'
   
"Portrait of Isidore Pils"
"Miracle of St. Sylvester; the Raising of Two Pagan Priests, after Maso di Banco, S. Croce, Florence"
"Entrance card to the sale exhibition of Edgar Degas's collection at Galerie Georges Petit, Paris, opening March 24, 1918"
"Saint Agatha"
"Man Leaning on a Crutch with Left Arm Raised"
"Studies of Four Englishmen, after James Gillray"
"The Justice of Trajan, plate for l'Album des Salons"
"Sheet of laid paper from the estate of Eugène Delacroix with the watermark: D&C Blauw and HP in a shield"
"Portrait of the photographer Lochard"
"Sheet of laid paper from the estate of Eugène Delacroix with watermark: D&C Blauw and HP in a shield."
"Sheet of laid paper from the estate of Eugène Delacroix with watermark: D&C Blauw and HP in a shield"
"Le Cap Martin"
"Three Children in a Landscape"
"Anatomical Studies of Legs"
"Children Seated on a Hillside with a Racquet"
"Cré chien, Loïse! T'as là une casquette un peu chouette!, plate XVIII from the suite Les Lorettes"
"Alors, si vous permettez, j'aurai l'honneur de vous envoyer ma voiture à onze heures...Ça me botte!, published in L'Éclair, May 22, 1852"
"Portrait of Charles Hayem"
"Unfinished Portrait of Degas"
"Unfinished Portrait of Degas"
"Portrait of Benjamin Mocatta (1802–1865)"
"Self-Portrait as Monsieur Prudhomme"
"Henry Monnier, from Galerie contemporaine, littéraire, artistique"
"M. Henri [sic] Monnier from L'Éclipse, October 25, 1874"
"Self-Portrait as Monsieur Prudhomme"
"Self-Portrait as Monsieur Prudhomme"
"Le Perroquet de M. Cheney: \"Vos injures n'atteindront jamais à la hauteur de mon dédain!\" (Mr. Cheney's Parrot: \"Your insults will never reach the height of my disdain!\")"
"David and Goliath, after a Capital in Vézelay Abbey"
"A Trial Plate for Mallarmé's \"Un coup de dés\" (A Throw of the Dice)"
```

The default query mode is to ensure that all queries match but you can also specify that only one or more queries need to match by passing the `-query-mode ANY` flag. For example:

```
$> bin/emit \
	-bucket-uri file:///usr/local/openaccess \
	-query 'Credit Line=(?:)foundation' \
	-query 'Credit Line=(?:)fund' \
	-query-mode ANY \
   | jq '.["Credit Line"]' \
   | uniq \
   | sort
   
"2013.905, .906: Purchase, Bequest of Stephen V. Grancsay, Rogers Fund, Helmut Nickel Gift, and funds from various donors, by exchange, 2013; 2016.311, .312: Purchase, Michael H. Pourfar Gift, 2016; 2016.403.1–.7: Purchase, Arthur Ochs Sulzberger Gift, 2016; 2016.409: Purchase, Marica F. Vilcek Gift, 2016"
"66.199: Rogers Fund, 1966; 1984.17: Purchase, David and Dorothy Alexander and Mrs. Ridgeley Hunt Gifts, Bequest of Stephen V. Grancsay, by exchange, and funds from various donors, 1984; 1988.170: Gift of Mr. and Mrs. Ronald S. Lauder, 1988"
"Administrative Purchase fund,  1976"
"Arthur Hoppock Hearn Fund and funds from various donors, 1990"
"Bequest of Nellie Kuh, by exchange, Louis V. Bell Fund, and funds from various donors, 1967"
"Brooklyn Museum Costume Collection at The Metropolitan Museum of Art, Gift of the Brooklyn Museum, 2009; Museum Expedition 1923, Purchased with funds given by Frederic B. Pratt and Frank L. Babbott"
"Brooklyn Museum Costume Collection at The Metropolitan Museum of Art, Gift of the Brooklyn Museum, 2009; Museum Expedition 1923, purchased with funds given by Frederic B. Pratt and Frank L. Babbott, 1923"
...and so on
```

### images

Command line tool to generate a CSV document mapping Open Access `Link Resource` URLs to their corresponding "main" and "download" image URLs. It is designed to be used in concert with the `emit` tool and any records marked as `Is Public Domain: false` are excluded.

```
$> ./bin/images -h
Usage of ./bin/images:
  -cookie-name string
    	A valid incap_ses_{SUFFIX} cookie name.
  -cookie-value string
    	A valid incap_ses_{SUFFIX} cookie value.
  -with-archive string
    	The path to an existing CSV file containing image URL mappings. Any URLs listed in this file will be included in the output as is and not retrieved from the metmuseum.org website.
```

_You should not need to use this tool as its output is bundled in the [data/images.csv.bz2](data/) file._

For example:

```
$> bin/emit \
	-bucket-uri file:///usr/local/openaccess/ \

   | bin/images \
	-with-archive data/images.csv \
	-cookie-name {COOKIE_NAME} -cookie-value {COOKIE_VALUE} \

   > images.csv
```

The `-cookie-name` and `-cookie-value` parameters are the name and value of a valid `incap_ses_{SUFFIX}` cookie. I have found the easiest way to deal with this is simply to vist the `metmuseum.org` website in a browser and, using the developer tools, copy and paste the relevant cookie data.

Hopefully future releases of the [openaccess data](https://github.com/metmuseum/openaccess) will include image URL information so that this tool won't be necessary anymore.

## See also

* https://github.com/metmuseum/openaccess