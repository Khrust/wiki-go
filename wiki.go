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
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func getTitle(url string) (string, error){
	match := validPath.FindStringSubmatch(url)
	if match == nil{
		return "", errors.New("Invalid page title")
	}
	return  match[2], nil
}

type Page struct {
	Title string
	Body  []byte
}

func renderTemplate(w http.ResponseWriter, templateName string, p *Page){
	err := templates.ExecuteTemplate(w, templateName, p)
	if err != nil{
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

func viewHandler(w http.ResponseWriter, r *http.Request) {
	Title, err := getTitle(r.URL.Path)
	if err != nil{
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	page, err := loadPage(Title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+Title, http.StatusFound)
		return
	}
	renderTemplate(w,"view.html", page)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	Title := r.URL.Path[len("/edit/"):]
	page, err := loadPage(Title)
	if err != nil {
		page = &Page{Title: Title}
	}
	renderTemplate(w,"edit.html", page)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	Title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: Title, Body: []byte(body)}
	err := p.save()
	if err!=nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+Title, http.StatusFound)
}

func main() {
	page1 := &Page{Title: "TestPage", Body: []byte("This is a sample content")}
	page1.save()
	page2, _ := loadPage("TestPage")
	fmt.Println(string(page2.Body))

	/*http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi! I love %s!", r.URL.Path[1:])
	})*/

	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)

}
