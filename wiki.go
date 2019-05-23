package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseFiles("view.html", "edit.html"))
var validPath = regexp.MustCompile("^/(edit|save|view|[a-zA-Z0-9]+)/*([a-zA-Z0-9]+)*$")

func getTitle(url string) (string, error) {
	match := validPath.FindStringSubmatch(url)
	if match == nil {
		return "", errors.New("Invalid page title")
	}
	title := match[2]
	if len(match[2]) == 0 {
		title = match[1]
	}
	return title, nil
}

type Page struct {
	Title string
	Body  []byte
}

func renderTemplate(w http.ResponseWriter, templateName string, p *Page) {
	err := templates.ExecuteTemplate(w, templateName, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(Title string) (*Page, error) {
	filename := Title + ".txt"
	Body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: Title, Body: Body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view.html", page)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}
	renderTemplate(w, "edit.html", page)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func rootHandler(w http.ResponseWriter, r *http.Request, title string) {
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Title, err := getTitle(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		fn(w, r, Title)
	}
}

func main() {
	page1 := &Page{Title: "TestPage", Body: []byte("This is a sample content")}
	page1.save()
	page2, _ := loadPage("TestPage")
	fmt.Println(string(page2.Body))

	/*http.HandleFunc("/", makeHandler(func(w http.ResponseWriter, r *http.Request, title string) {
		http.Redirect(w, r, "/view/"+title, http.StatusFound)
	}))*/

	http.HandleFunc("/", makeHandler(rootHandler))
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":8080", nil)
	fmt.Println("Pezdec")

}
