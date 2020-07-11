## images.csv.bz2

A CSV file listing the "main" and "download" URL associated with Met Museum Open Access records that have been marked as public domain.

In order to reduce the size of this file common prefixes have been removed. For example:

```
$> bunzip2 -c -d images.csv.bz2 | less
url,main_image,download_image
34,34/20602/main-image,ad/original/204788.jpg
37,37/46007/main-image,ad/original/DP247752.jpg
38,38/46008/main-image,ad/original/DP247753.jpg
39,39/4360/main-image,ad/original/37808.jpg
40,40/17983/main-image,ad/original/174118.jpg
...
```

The columns contained in the file are:

| index | column | prefix (stripped) |
| --- | --- | --- |
| 0 | url | http://www.metmuseum.org/art/collection/search/ |
| 1 | main_image | https://collectionapi.metmuseum.org/api/collection/v1/iiif/ |
| 2 | download_image | https://images.metmuseum.org/CRDImages/ |

The first column (`url`) maps to the Open Access `Link Resource` property.