package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"html/template"
)

type Page struct {
	// struct is like a class, fields are initializers
	Title string
	Body []byte
}

// This method saves the Page's body to a text file, Title is used 
// as the name of the text file
func (p *Page) save() error {
	// method save, receiver p, pointer Page, 
	// save takes no parameters, returns value of the type "error"
	filename := p.Title + ".txt"
	
	// the return value is of type error bc that is the return type of
	// WriteFile in Go std-lib 
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	// constructs filename from title parameter
	filename := title + ".txt"

	// io.ReadFile returns []byte and error
	body, err := ioutil.ReadFile(filename)

	// if second paramter is nil, page has loaded successfully 
	if err != nil {
		// if not nil, error is returned
		return nil, err
	}
	// loads successfully with fields and nil
	return &Page{Title: title, Body: body}, nil
}
 
// template code to handle parsing view files
// DONT REPEAT YOURSELF

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page){
	t, err := template.ParseFiles(tmpl + ".html")
	// Error handling	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// http.error sends specified HTTP response code(custom response)
}


func viewHandler(w http.ResponseWriter, r *http.Request){
	// extracts page title from the path component of request url
	title := r.URL.Path[len("/view/"):]
	// path is then sliced, /view/ is taken out of title

	// loads page data, using blank identifier to throw out error
	p, _ := loadPage(title)

	if err != nil {
		// redirects to edit page to add data of non is found
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	// Removed and made into method
	// t, _ := template.ParseFiles("view.html")
	// t.Execute(w, p) 

	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request){
	// extracts page titl from the path request url
	// ie if title is = to edit based on path
	title := r.URL.Path[len("/edit/"):]

	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	// else if nil

	// // html stuff from template
	// Removed and made into method
	// t, _ := template.ParseFiles("edit.html")
	// t.Execute(w, p)

	renderTemplate(w, "edit", p)

}

func saveHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	// []byte(body) converts the FormValue of string to []byte before fitting in Page struct
	p := &Page{Title: title, Body: []byte(body)}
	p.save()

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// main function to test methods
// small note: goes below everything
func main(){
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}


