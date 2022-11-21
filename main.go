package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/PuerkitoBio/goquery"
)

// helper functions
func Prompt(question string) (ans string) {
    fmt.Print(question)
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan()
    ans = scanner.Text()
    return
}

func Mkdir(dir string) {
    if e := os.Mkdir(dir, os.ModePerm); e != nil {
        return
    }
}

func WriteFile(text string, flag int, src string) {
    f, e := os.OpenFile(src, flag, 0644)
    if e != nil {
        panic(e)
    }
    _, e = fmt.Fprint(f, text)
    if e != nil {
        panic(e)
    }
}

type Data struct {
    Title string
    Head string
    Divs string
}

func main() {
    // get information
    rawUrl := Prompt("Enter 'url': ")
    elementsString := Prompt("Enter 'element/s': ")
    elements := strings.Fields(elementsString)
    outDir := Prompt("Enter 'output directory': ")

    // get html of page
    res, err := http.Get(rawUrl)
    if err != nil {panic(err)}
    body, err := io.ReadAll(res.Body)
    if err != nil {panic(err)}
    html := string(body)

    // make goquery document
    reader := strings.NewReader(html)
    doc, err := goquery.NewDocumentFromReader(reader)
    if err != nil {panic(err)}

    URL, _ := url.Parse(rawUrl)
    Mkdir(outDir)

    var (
        title = doc.Find("title").Text()
        head = ""
        divs = ""
    )

    // edit link tag
    doc.Find("link").Each(func(i int, s *goquery.Selection) {
        rel, rex := s.Attr("rel")
        t, tyex := s.Attr("type")
        href, hex := s.Attr("href")
        title, tiex := s.Attr("title")
        out := "<link"
        if rex {
            out += " rel=\"" + rel + "\""
        }
        if hex {
            if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
                out += " href=\"" + href + "\""
            } else {
                out += " href=\"" + URL.Scheme + "://" + URL.Host + href + "\""
            }
        }
        if tiex {
            out += " title=\"" + title + "\""
        }
        if tyex {
            out += " type=\"" + t + "\""
        }
        out += " >"
        head += out
    })

    // edit style tag
    doc.Find("style").Each(func(i int, s *goquery.Selection) {
        t, _ := s.Attr("type")
        head += "<style type=\"" + t + "\">" + s.Text() + "</style>"
    })

    // edit script tag
    doc.Find("script").Each(func(i int, s *goquery.Selection) {
        t, tex := s.Attr("type")
        src, srcex := s.Attr("src")
        _, aex := s.Attr("async")
        chset, chex := s.Attr("charset")
        out := "<script"
        if tex {
            out +=" type=\"" + t + "\""
        }
        if srcex {
            if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
                out += " src=\"" + src + "\""
            } else {
                out += " src=\"" + URL.Scheme + "://" + URL.Host + src + "\""
            }
        }
        if aex {
            out += " async"
        }
        if chex {
            out += " charset=\"" + chset + "\""
        }
        if s.Text() != "" {
            out += " >" + s.Text() + "</script>"
        } else {
            out += "></script>"
        }
        head += out
    })

    // edit img tag
    doc.Find("img").Each(func(i int, s *goquery.Selection) {
        src, _ := s.Attr("src")
        if (!strings.HasPrefix(src, "http://") || !strings.HasPrefix(src, "https://")) && !strings.HasPrefix(src, "//") {
            s.SetAttr("src", URL.Scheme + "://" + URL.Host + src)
        }
    })

    // edit a tag
    doc.Find("a").Each(func(i int, s *goquery.Selection) {
        href, ex := s.Attr("href")
        if !strings.HasPrefix(href, "#") && !strings.HasPrefix(href, "http://") && !strings.HasPrefix(href, "https://") && ex {
            s.SetAttr("href", URL.Scheme + "://" + URL.Host + href)
        }
    })

    // look for entered elements
    for _, e := range elements {
        doc.Find(e).Each(func(i int, s *goquery.Selection) {
            class, _ := s.Attr("class")
            id, _ := s.Attr("id")
            ret, err := s.Html()
            if err != nil {
                panic(err)
            }
            divs += "<div class=\"" + class + "\" id=\"" + id +"\" >" + ret + "</div>"
        })
    }

    // html template for output page
    tmpl := template.Must(template.New("").Parse(`<!DOCTYPE html>
<html lang="en">
    <head>
        {{.Head}}
        <title>{{.Title}}</title>
    </head>
    <body>
        {{.Divs}}
    </body>
</html>
`))

    // create file
    f, err := os.Create(outDir + title + ".html")
    if err != nil {panic(err)}
    defer func() {f.Close()}()
    tmpl.Execute(f, Data{title, head, divs})

    // open output file in browser
    path, _ := filepath.Abs("./" + outDir + title + ".html")
    var args []string
    switch runtime.GOOS {
    case "darwin":
        args = []string{"open", path}
    case "windows":
        args = []string{"cmd", "/c", "start", "chrome", path}
    default:
        args = []string{"xdg-open", path}
    }
    cmd := exec.Command(args[0], args[1:]...)
    e := cmd.Run()
    if e != nil {
        log.Printf("openinbrowser: %v\n", err)
    }
}