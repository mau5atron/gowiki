package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"html/template"
	"regexp"
)

type Page struct {
	// struct is like a class, fields are initializers
	Title string
	Body []byte
}

// global variable
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
// Once the program initializes, ParseFiles will be called instead of 
// inefficiently calling it twice with the renderTemplate method

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
/* 
 - regexp.MustCompile will parse and compile the regular expression, and 
 	 return a regexp.Regexp
 - MustCompile is distinct from Compile in that it will panic if expression 
 	 compilation fails.
 - Compile only returns an error as second parameter
*/

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
	// t, err := template.ParseFiles(tmpl + ".html") dont need this

	// calls the templates.ExecuteTemplate method with name of appropriate template
	err := templates.ExecuteTemplate(w, tmpl+".html", p)

	// Error handling	
	// removed some error handling stuff

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// http.error sends specified HTTP response code(custom response)
}


func viewHandler(w http.ResponseWriter, r *http.Request, title string){
	// extracts page title from the path component of request url
	title, err := getTitle(w, r)
	// also validates page title
	if err != nil {
		return
	}

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

func editHandler(w http.ResponseWriter, r *http.Request, title string){
	// extracts page titl from the path request url
	// ie if title is = to edit based on path
	// title := r.URL.Path[len("/edit/"):]
	title, err := getTitle(w, r)

	if err != nil {
		return
	}

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

func saveHandler(w http.ResponseWriter, r *http.Request, title string){
	// title := r.URL.Path[len("/save/"):]

	title, err := getTitle(w, r)
	if err != nil {
		return
	}

	body := r.FormValue("body")
	// []byte(body) converts the FormValue of string to []byte before fitting in Page struct
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	// error handling, any errors that occur during save() will be reported to user	
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// getTitle validates path with validPath expression to extract page title
// func getTitle(w http.ResponseWriter, r *http.Request)(string, error){
// 	m := validPath.FindStringSubmatch(r.URL.Path)

// 	if m == nil {
// 		http.NotFound(w, r)
// 		return "", errors.New("Invalid Page Title")
// 	}
// 	return m[2], nil // title is second subexpression
// }

func makeHandler(fn func(http,ResponseWriter, *http.Request, string)) http.HandleFunc {
	return func(w http.ResponseWriter, r *http.Request){
		m := validPath.FindStringSubmatch(r,URL.Request)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}


// main function to test methods
// small note: goes below everything
func main(){
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}


