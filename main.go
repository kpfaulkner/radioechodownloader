package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

func downloadURLToFile( url string, filepath string) error {

	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}
	defer response.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, response.Body)
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}

	return nil
}

func getTitle( contents string) string {

	title := ""

	return title
}

func getLinks( contents string ) ([]string, error) {
	urls := make([]string, 4)
	mp3Regex := regexp.MustCompile(`dl_mp3folder=(.*?)&dl_file=(.*?)&`)

	urlMap := make(map[string]string)

	//res := mp3Regex.FindAllString(contents,-1)
    res := mp3Regex.FindAllStringSubmatch(contents, -1)
	if res != nil {
		fmt.Printf("number of hits %d\n", len(res))

		for _,i := range res {
			directory := i[1]
			name := i[2]

			fmt.Printf("dir %s : %s\n", directory, name)


			url := fmt.Sprintf("https://s3.amazonaws.com/RE-Warehouse/%s/%s", strings.ToLower(directory), name)

			urlMap[url] = url
//			urls = append(urls, url)
		}

	}

    for k,_ := range urlMap {
    	urls = append(urls, k)
	}


	return urls, nil
}

func getDownloadPath( rawUrl string, downloadPath string ) string {

	filename := path.Base(rawUrl)
	fullPath := fmt.Sprintf("%s%s", downloadPath, filename )
	return fullPath
}

func getFilename( rawUrl string ) string {
	b := path.Base(rawUrl)
	return b
}

func main() {

	if len(os.Args) != 3 {
		fmt.Printf("radioechodownloader <url> <download dir>.\n")
		fmt.Println("eg. radioechodownloader \"http://www.radioechoes.com/?page=series&genre=OTR-Comedy&series=Im%20Sorry%20I%20Havent%20A%20Clue\" c:\\temp\\radioecho\\")
		os.Exit(1)
	}

	initialURL := os.Args[1]
	downloadDir := os.Args[2]
	response, err := http.Get(initialURL)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", string(contents))

	strContents := string(contents)
	title := getTitle(strContents)
	fmt.Printf("title is %s\n", title)

	urlLinks,err := getLinks(strContents)
	if err != nil {
		fmt.Printf("ERROR trying to get links: %v\n", err)
		os.Exit(1)
	}

	for _, url := range(urlLinks) {
		fmt.Printf("download %s\n", url)

		downloadFilename := getDownloadPath(url, downloadDir)
		err := downloadURLToFile( url, downloadFilename)
		if err != nil {
			fmt.Printf("Unable to download %s\n", url)
		}
	}

}
