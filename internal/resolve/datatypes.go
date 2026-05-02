package resolve

type AURResult struct {
	Name         string
	Version      string
	Description  string
	URL          string
	LastModified int64
	NumVotes     int
}

type PacmanResult struct {
	Name        string
	Version     string
	Description string
	Repository  string
}
