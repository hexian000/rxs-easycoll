/* This is a go script. Usage: `go run gen_list.go` */

package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type Prescriptions struct {
	XMLName xml.Name       `xml:"prescriptions"`
	Rx      []Prescription `xml:"prescription"`
}

type Prescription struct {
	XMLName xml.Name `xml:"prescription"`
	Link    string   `xml:"link,attr"`
}

func main() {
	var rxs Prescriptions
	b, err := ioutil.ReadFile("rxs-easycoll.xml")
	if err != nil {
		log.Fatalln(err)
	}
	err = xml.Unmarshal(b, &rxs)
	if err != nil {
		log.Fatalln(err)
	}

	buf := bytes.NewBuffer(b[:0])
	buf.WriteString("| 名称 | 作者 | 链接 |\n")
	buf.WriteString("|---|---|---|\n")

	for _, rx := range rxs.Rx {
		v := strings.Split(rx.Link[1:], "/")
		user, repo := v[0], v[1]
		buf.WriteString(fmt.Sprintf("| %s | %s | [Link](https://github.com%s) |\n",
			user, repo, rx.Link))
	}

	err = ioutil.WriteFile("list.md", buf.Bytes(), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
