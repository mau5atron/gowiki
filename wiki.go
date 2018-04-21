package main

import (
	"fmt"
	"io/ioutil"
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



// main function to test methods

func main(){
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample page.")}
	p1.save()

	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))
}


