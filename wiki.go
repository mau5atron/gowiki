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



func viewHandler(w http.ResponseWriter, r *http.Request, title string){
	// loads page data, using blank identifier to throw out error
	p, err := loadPage(title)
	if err != nil {
		// redirects to edit page to add data of non is found
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string){
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string){
	body := r.FormValue("body")
	// []byte(body) converts the FormValue of string to []byte before fitting in Page struct
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	// error handling, any errors that occur during save() will be reported to user	
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// global variable
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
// Once the program initializes, ParseFiles will be called instead of 
// inefficiently calling it twice with the renderTemplate method
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

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
/* 
 - regexp.MustCompile will parse and compile the regular expression, and 
 	 return a regexp.Regexp
 - MustCompile is distinct from Compile in that it will panic if expression 
 	 compilation fails.
 - Compile only returns an error as second parameter
*/
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		m := validPath.FindStringSubmatch(r.URL.Path)
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


