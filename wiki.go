package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type Page struct {
	Title string
	Body  []byte
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
	Title := r.URL.Path[len("/view/"):]
	p, err := loadPage(Title)
	if err != nil {
		return
	}
	t, _ := template.ParseFiles("view.html")
	t.Execute(w, p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	Title := r.URL.Path[len("/edit/"):]
	page, err := loadPage(Title)
	if err != nil {
		page = &Page{Title: Title}
	}
	t, _ := template.ParseFiles("edit.html")
	t.Execute(w, page)

}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	Title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: Title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+Title, http.StatusFound)
}

func main() {
	page1 := &Page{Title: "TestPage", Body: []byte("This is a sample content")}
	page1.save()
	page2, _ := loadPage("TestPage")
	fmt.Println(string(page2.Body))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi! I love %s!", r.URL.Path[1:])
	})

	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)

}

func renderTemplate(w http.ResponseWriter, r *http.Request) {
	t, _ := 
}
