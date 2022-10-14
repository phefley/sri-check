package sricheck

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Permitted algorithms are: sha256, sha384, and sha512

func getAlgorithms() []string {
	return []string{"sha256", "sha384", "sha512"}
}

func checkAlgorithm(algorithm string) bool {
	for _, a := range getAlgorithms() {
		if a == algorithm {
			return true
		}
	}
	return false
}

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

func generateSriMapFromData(data []byte) map[string]string {
	outvar := make(map[string]string)
	outvar["sha256"] = generateBase64HashOf(data, "sha256")
	outvar["sha384"] = generateBase64HashOf(data, "sha384")
	outvar["sha512"] = generateBase64HashOf(data, "sha512")
	return outvar
}

func generateSriMap(jsSrc string) map[string]string {
	var outvar map[string]string
	// Get the resource at jsSrc
	resp, err := http.Get(jsSrc)
	if err != nil {
		log.Fatal(err)
		return outvar
	}
	if resp.StatusCode != 200 {
		log.Fatal("[*] Non 200 status code while obtaining ", jsSrc)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}
	return generateSriMapFromData(data)
}

func CheckIntegrity(jsSrc string, integrity string) bool {
	sriMap := generateSriMap(jsSrc)
	intSplit := strings.Split(integrity, "-")
	if checkAlgorithm(intSplit[0]) {
		return sriMap[intSplit[0]] == intSplit[1]
	} else {
		return false
	}
}
