package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type html string

func betweenTextFirstCut(s, begin, end string, cutL, cutR int) (out string, err error) {
	var indBegin, indEnd int
	indBegin = strings.Index(s, begin)
	if indBegin < indEnd {
		err = fmt.Errorf("Не найдено начало: %v", begin)
		return
	}
	indBegin += len(begin)
	indEnd = strings.Index(s[indBegin:], end) + indBegin
	if indEnd-indBegin < 1 {
		err = fmt.Errorf("Не найден конец: %v, или длинна найденной строки - 0", end)
		return
	}
	out = s[indBegin+cutL : indEnd+cutR]
	return
}

func betweenTextAllCut(s, begin, end string, cutL, cutR int) (out []string) {
	var indBegin, indEnd int
	for {
		indBegin = strings.Index(s[indEnd:], begin) + indEnd
		if indBegin < indEnd {
			break
		}
		indBegin += len(begin)
		indEnd = strings.Index(s[indBegin:], end) + indBegin
		if indEnd-indBegin <= 1 {
			break
		}
		out = append(out, s[indBegin+cutL:indEnd+cutR])
	}
	return
}

func (h *html) HtmlBase() (out string, err error) {
	out, err = betweenTextFirstCut(string(*h), "<base", ">", 0, 0)
	if err != nil {
		err = fmt.Errorf("Тег <base> не найден")
		return
	}
	return
}

func (h html) HtmlBaseHref() (out string, err error) {
	var s string
	s, err = h.HtmlBase()
	if err != nil {
		return
	}
	out, err = betweenTextFirstCut(s, "href=\"", "\"", 0, 0)
	if err != nil {
		err = fmt.Errorf("Атрибут href не найден")
		return
	}
	return
}

func (h html) HtmlA() (out []string) {
	out = betweenTextAllCut(string(h), "<a", "></a>", 0, 0)
	return
}

func (h html) HtmlAHref() (out []string) {
	for _, v := range h.HtmlA() {
		out = append(out, betweenTextAllCut(v, "href=\"", "\"", 0, 0)...)
	}
	return
}

func Crawl(host string) (out []string) {
	c := make(chan string)
	go func() {
		for v := range c {
			out = append(out, v)
		}
	}()
	intoCrawl(host, &out, c)
	close(c)
	return
}

func intoCrawl(host string, inout *[]string, c chan string) {
	var (
		newHost string
		tagA    []string
	)
	resp, err := http.Get(host)
	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		return
	}

	s1 := fmt.Sprint(resp.Request.URL)[len(host):]
	if len(s1) != 0 {
		for _, v := range *inout {
			if v == s1 {
				return
			}
		}
		c <- s1
		newHost = host
	} else {
		newHost = "http://" + resp.Request.Host
		s2 := fmt.Sprint(resp.Request.URL)[len(newHost):]
		for _, v := range *inout {
			if v == s2 {
				return
			}
		}
		c <- s2
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	sbody := string(body)
	hbody := html(sbody)
	tagBase, errTB := hbody.HtmlBaseHref()
	tagA = hbody.HtmlAHref()

	if errTB == nil {
		for i, v := range tagA {
			if v[:1] != "/" && v[:1] != "h" {
				tagA[i] = tagBase + v
			}
		}
	}
	for _, v := range tagA {
		if v[:1] == "/" {
			intoCrawl(newHost+v, inout, c)
		}
	}
	resp.Body.Close()
	return
}

func main() {
	fmt.Println("Hi!")
}
