package openaccess

// maybe URI templates?

const LINK_RESOURCE_PREFIX string = "http://www.metmuseum.org/art/collection/search/"
const MAIN_IMAGE_PREFIX string = "https://collectionapi.metmuseum.org/api/collection/v1/iiif/"
const DOWNLOAD_IMAGE_PREFIX string = "https://images.metmuseum.org/CRDImages/"

type OpenAccessRecord struct {
	ArtistDisplayBio      string `json:"Artist Display Bio"`
	ArtistNationality     string `json:"Artist Nationality"`
	SubRegion             string `json:"Subregion"`
	ObjectWikidataURL     string `json:"Object Wikidata URL"`
	TagsAATURL            string `json:"Tags AAT URL"`
	ObjectName            string `json:"Object Name"`
	Portfolio             string `json:"Portfolio"`
	ArtistAlphaSort       string `json:"Artist Alpha Sort"`
	ObjectEndDate         string `json:"Object End Date"`
	Excavation            string `json:"Excavation"`
	MetadataDate          string `json:"Metadata Date"`
	ObjectID              string `json:"Object ID"`
	Department            string `json:"Department"`
	Country               string `json:"Country"`
	Locus                 string `json:"Locus"`
	Period                string `json:"Period"`
	County                string `json:"County"`
	Culture               string `json:"Culture"`
	Dynasty               string `json:"Dynasty"`
	Reign                 string `json:"Reign"`
	ArtistWikidataURL     string `json:"Artist Wikidata URL"`
	Locale                string `json:"Locale"`
	IsPublicDomain        bool   `json:"Is Public Domain"`
	Title                 string `json:"Title"`
	ArtistRole            string `json:"Artist Role"`
	ArtistEndDate         string `json:"Artist End Date"`
	CreditLine            string `json:"Credit Line"`
	Repository            string `json:"Repository"`
	ObjectNumber          string `json:"Object Number"`
	IsTimelineWork        bool   `json:"Is Timeline Work"`
	ArtistBeginDate       string `json:"Artist Begin Date"`
	Dimensions            string `json:"Dimensions"`
	GeographyType         string `json:"Geography Type"`
	Classification        string `json:"Classification"`
	Tags                  string `json:"Tags"`
	AccessionYear         string `json:"AccessionYear"`
	ArtistDisplayName     string `json:"Artist Display Name"`
	ArtistSuffic          string `json:"Artist Suffix"`
	ArtistsULANURL        string `json:"Artist ULAN URL"`
	ObjectDate            string `json:"Object Date"`
	City                  string `json:"City"`
	State                 string `json:"State"`
	RightsAndReproduction string `json:"Rights and Reproduction"`
	IsHighlight           bool   `json:"Is Highlight"`
	ArtistPrefix          string `json:"Artist Prefix"`
	LinkResource          string `json:"Link Resource"`
	Medium                string `json:"Medium"`
	Region                string `json:"Region"`
	River                 string `json:"River"`
	ArtistGender          string `json:"Artist Gender"`
	ObjectBeginDate       string `json:"Object Begin Date"`
	MainImage             string `json:"Main Image,omitempty"`
	DownloadImage         string `json:"Download Image,omitempty"`
}
