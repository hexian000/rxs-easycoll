/* This is a go script. Usage: `go run gen_list.go` */

package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
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

type ParsedRx struct {
	Repo, User string
}

type List []ParsedRx

func (l List) Len() int { return len(l) }
func (l List) Less(i, j int) bool {
	c := strings.Compare(
		strings.ToLower(l[i].User),
		strings.ToLower(l[j].User),
	)
	if c == 0 {
		c = strings.Compare(
			l[i].User,
			l[j].User,
		)
	}
	if c == 0 {
		c = strings.Compare(
			strings.ToLower(l[i].Repo),
			strings.ToLower(l[j].Repo),
		)
	}
	if c == 0 {
		c = strings.Compare(
			l[i].Repo,
			l[j].Repo,
		)
	}
	return c < 0
}
func (l List) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
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

	prx := make([]ParsedRx, 0)
	for _, rx := range rxs.Rx {
		v := strings.Split(rx.Link[1:], "/")
		prx = append(prx, ParsedRx{v[1], v[0]})
	}
	sort.Stable(List(prx))

	buf := bytes.NewBuffer(b[:0])
	buf.WriteString("# 致谢（字典序）\n\n")
	buf.WriteString("| 名称 | 作者 |\n")
	buf.WriteString("| --- | --- |\n")

	for _, rx := range prx {
		buf.WriteString(fmt.Sprintf("| [%[2]s](https://github.com/%[1]s/%[2]s) | [%[1]s](https://github.com/%[1]s) |\n",
			rx.User, rx.Repo))
	}

	err = ioutil.WriteFile("CREDITS.md", buf.Bytes(), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
