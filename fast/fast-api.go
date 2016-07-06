package fast

import "fmt"
import "bytes"
import "net/http"
import "io"
import "regexp"

// UseHTTPS sets if HTTPS is used
var UseHTTPS = true

// GetDlUrls returns a list of urls to the fast api downloads
func GetDlUrls(urlCount uint64) (urls []string) {
	token := getFastToken()
	// fmt.Printf("token=%s\n", token)

	httpProtocol := "https"
	if !UseHTTPS {
		httpProtocol = "http"
	}

	url := fmt.Sprintf("%s://api.fast.com/netflix/speedtest?https=%t&token=%s&urlCount=%d",
		httpProtocol, UseHTTPS, token, urlCount)
	// fmt.Printf("url=%s\n", url)

	jsonData, _ := getPage(url)

	re := regexp.MustCompile("(?U)\"url\":\"(.*)\"")
	reUrls := re.FindAllStringSubmatch(jsonData, -1)

	for _, arr := range reUrls {
		urls = append(urls, arr[1])
	}

	return
}

// GetDefaultURL returns the fallback download URL
func GetDefaultURL() (url string) {
	httpProtocol := "https"
	if !UseHTTPS {
		httpProtocol = "http"
	}
	url = fmt.Sprintf("%s://api.fast.com/netflix/speedtest", httpProtocol)
	return
}

func getFastToken() (token string) {
	baseURL := "https://fast.com/"
	if !UseHTTPS {
		baseURL = "http://fast.com"
	}
	fastBody, _ := getPage(baseURL)

	// Extract the app script url
	re := regexp.MustCompile("app-.*\\.js")
	scriptNames := re.FindAllString(fastBody, 1)

	scriptURL := fmt.Sprintf("%s/%s", baseURL, scriptNames[0])
	// fmt.Printf("scriptUrl=%s\n", scriptURL)

	// Extract the token
	scriptBody, _ := getPage(scriptURL)
	// fmt.Printf("got:\n-----\n%s\n-----\n", scriptBody)

	re = regexp.MustCompile("token:\"[[:alpha:]]*\"")
	tokens := re.FindAllString(scriptBody, 1)

	if len(tokens) > 0 {
		token = tokens[0][7 : len(tokens[0])-1]
	} else {
		fmt.Printf("no token found\n")
	}

	return
}

func getPage(url string) (contents string, err error) {
	// Create the string buffer
	buffer := bytes.NewBuffer(nil)

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return contents, err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(buffer, resp.Body)
	if err != nil {
		return contents, err
	}
	contents = buffer.String()

	return
}