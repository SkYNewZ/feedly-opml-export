// This pacakge export OPML file from Feddly to stdout
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"time"
)

const feedlyOPMLURL string = "https://cloud.feedly.com/v3/opml"

var accessToken string
var version string = "dev"

type feedlyErrorResponse struct {
	ErrorCode           int    `json:"errorCode"`
	ErrorID             string `json:"errorId"`
	ErrorMessage        string `json:"errorMessage"`
	ParsedErrorResponse string
}

func (f feedlyErrorResponse) Error() string {
	if f.ParsedErrorResponse != "" {
		return fmt.Sprintf("Feedly error: %s", f.ParsedErrorResponse)
	}

	return fmt.Sprintf("Feedly error: %d %s", f.ErrorCode, f.ErrorMessage)
}

func init() {
	flag.Parse()
}

func main() {
	// Print version if asked
	if flag.Arg(0) == "version" {
		fmt.Printf("feedly-opml-export %s compiled with %v on %v/%v\n", version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	// Ensure token is given
	accessToken, ok := os.LookupEnv("FEEDLY_ACCESS_TOKEN")
	if !ok || accessToken == "" {
		handleError(fmt.Errorf("Mising or invalid $FEEDLY_ACCESS_TOKEN"))
	}

	// Make request
	req, err := http.NewRequest("GET", feedlyOPMLURL, nil)
	handleError(err)

	// Add authentication
	req.Header.Add("Authorization", fmt.Sprintf("OAuth %s", accessToken))
	resp, err := http.DefaultClient.Do(req)
	handleError(err)

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		handleFeedlyError(resp)
	}

	// Read body
	data, err := ioutil.ReadAll(resp.Body)
	handleError(err)

	// Print result
	fmt.Printf(string(data))
}

// Parse Feedly response
func handleFeedlyError(resp *http.Response) {
	var e feedlyErrorResponse
	handleError(json.NewDecoder(resp.Body).Decode(&e))

	// Read timestamp if any
	var re = regexp.MustCompile(`\d+`)
	t := re.FindAllString(e.ErrorMessage, -1)
	if len(t) >= 2 {
		// t[0] is expired date as timestape
		// t[1] is the number of secondes since expiration
		expiredDate, err := msToTime(t[0])
		handleError(err)

		secs, err := time.ParseDuration(t[1] + "s")
		handleError(err)

		e.ParsedErrorResponse = fmt.Sprintf("Token expired since the %s (%s)", expiredDate.Format(time.RFC822), secs.String())
	}

	handleError(e)
}

// https://stackoverflow.com/a/13295158
func msToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(0, msInt*int64(time.Millisecond)), nil
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
