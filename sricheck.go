package sricheck

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

type PageScript struct {
	Url                string
	IntegrityAttribute string
}

/*
Create a slice of algorithms available for integrity attributes.
Permitted algorithms are: sha256, sha384, and sha512
*/
func getAlgorithms() []string {
	return []string{"sha256", "sha384", "sha512"}
}

/*
Check to see if the given string:algorithm is in the list of approved algorithms.
Return a boolean result TRUE if the provided string is a valid algorithm
*/
func checkAlgorithm(algorithm string) bool {
	for _, a := range getAlgorithms() {
		if a == algorithm {
			return true
		}
	}
	return false
}

/*
Given some data (in a slice of bytes), generate the base64 encoded
version of it's hash using the supplied algorithm
If the supplied algorithm isn't valid, return the empty string
*/
func generateBase64HashOf(data []byte, algorithm string) string {
	var retval string
	if checkAlgorithm(algorithm) {
		switch algorithm {
		case "sha256":
			interim256HashBytes := sha256.Sum256(data)
			retval = base64.StdEncoding.EncodeToString(interim256HashBytes[:])
		case "sha384":
			interim384HashBytes := sha512.Sum384(data)
			retval = base64.StdEncoding.EncodeToString(interim384HashBytes[:])
		case "sha512":
			interim512HashBytes := sha512.Sum512(data)
			retval = base64.StdEncoding.EncodeToString(interim512HashBytes[:])
		}
	}
	return retval
}

/*
Given a slice of data (in bytes), generate a map of hashes for that data
The index, or key, of the map will be a string of the hash algorithm
The value of the map will be the base64 hash for the data
*/
func generateSriMapFromData(data []byte) map[string]string {
	outvar := make(map[string]string)
	outvar["sha256"] = generateBase64HashOf(data, "sha256")
	outvar["sha384"] = generateBase64HashOf(data, "sha384")
	outvar["sha512"] = generateBase64HashOf(data, "sha512")
	return outvar
}

/*
Given a string URL to a page, request that page (via HTTP GET) and note the script blocks with SRC attributes
Return a slice of PageScript structs, which include the Url and the Integrity attribute for any noted JS includes
*/
func getPageScripts(pageUrl string) ([]PageScript, error) {
	var outvar []PageScript
	resp, err := http.Get(pageUrl)
	if err != nil {
		return outvar, err
	}
	if resp.StatusCode != 200 {
		return outvar, errors.New("Received a non-200 status code for " + pageUrl)
	} else {
		defer resp.Body.Close()
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return outvar, err
		}
		doc.Find("script").Each(func(i int, s *goquery.Selection) {
			src, hasSrc := s.Attr("src")
			if hasSrc {
				hash, hasHash := s.Attr("integrity")
				if hasHash {
					outvar = append(outvar, PageScript{Url: src, IntegrityAttribute: hash})
				} else {
					outvar = append(outvar, PageScript{Url: src})
				}
			}
		})
	}
	return outvar, nil
}

// PUBLIC FUNCTIONS

/*
Given a string URL to a file (generally a JS or CSS file), generate a map of hashes for that data
The index, or key, of the map will be a string of the hash algorithm
The value of the map will be the base64 hash for the data
*/
func GenerateSriMap(jsSrc string) (map[string]string, error) {
	var outvar map[string]string
	// Get the resource at jsSrc
	resp, err := http.Get(jsSrc)
	if err != nil {
		return outvar, err
	}
	if resp.StatusCode != 200 {
		return outvar, errors.New("Non 200 status code while obtaining " + jsSrc)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return outvar, err
	}
	return generateSriMapFromData(data), nil
}

/*
Given a string URL to a page, request that page (via HTTP GET) and note the script blocks with SRC attributes
Check the Integrity attribute for any noted JS includes
Return a boolean value based on the integrity checks - true when all integrity attributes are valid (or not present)
and false when an integrity attribute does not check correctly.
*/
func CheckPageIntegrity(pageUrl string) (bool, error) {
	scripts, err := getPageScripts(pageUrl)
	if err != nil {
		return false, err
	}
	for _, ps := range scripts {
		intCheck, err := CheckIntegrity(ps.Url, ps.IntegrityAttribute)
		if err != nil {
			return false, err
		}
		if !intCheck {
			return false, nil
		}
	}
	return true, nil
}

/*
Given a string URL to a page, request that page (via HTTP GET) and note the script blocks with SRC attributes
Check the Integrity attribute for any noted JS includes
Print a table summarizing all of the JS includes, their integrity attributes, and the validity of integrity attributes
*/
func PrintPageIntegrityCheckTable(pageUrl string) {
	scripts, err := getPageScripts(pageUrl)
	if err != nil {
		log.Fatal(err)
	}
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("ID", "SRC", "Integrity", "Valid?")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	for counter, ps := range scripts {
		integrityCheck, err := CheckIntegrity(ps.Url, ps.IntegrityAttribute)
		if err != nil {
			log.Fatal(err)
		}
		tbl.AddRow(counter, ps.Url, ps.IntegrityAttribute, strconv.FormatBool(integrityCheck))
	}
	tbl.Print()
}

/*
Given a string URL to a file (generally a JS or CSS file) and the string of
the integrity attribute, determine if the integrity attribute is valid
Returns true for valid and false for invalid
*/
func CheckIntegrity(jsSrc string, integrity string) (bool, error) {
	sriMap, err := GenerateSriMap(jsSrc)
	if err != nil {
		return false, err
	}
	intSplit := strings.Split(integrity, "-")
	if checkAlgorithm(intSplit[0]) {
		return sriMap[intSplit[0]] == intSplit[1], nil
	} else {
		return false, nil
	}
}
