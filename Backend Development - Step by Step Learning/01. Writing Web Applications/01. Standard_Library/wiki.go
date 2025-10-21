package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

/*
A wiki consists of a series of interconnected pages, each of which has a title and a body (the page content).
Here, we define Page as a struct with two fields representing the title and body.
*/
type Page struct {
	Title string
	Body []byte
}

/*
The Page struct describes how page data will be stored in memory. But what about persistent storage? 
We can address that by creating a save method on Page:
*/
func (p *Page) save() error { // The save method returns an error value because that is the return type of WriteFile
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
	/* The octal integer literal 0600, passed as the third parameter to WriteFile, 
		indicates that the file should be created with read-write permissions for the current user only. */
}

/*
The function loadPage constructs the file name from the title parameter, 
reads the file's contents into a new variable body, and returns a pointer 
to a Page literal constructed with the proper title and body values.
*/
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename) // The standard library function os.ReadFile returns []byte and error.
	
	if err != nil {
		return nil, err
	}
	
	return &Page {
		Title: title,
		Body: body,
	}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/" + title, http.StatusFound)
		return
	}

	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{
			Title: title,
		}
	}
	renderTemplate(w, "edit", p)
}

/*
The function saveHandler will handle the submission of forms located on the edit pages. 
After uncommenting the related line in main
*/
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	/* The page title (provided in the URL) and the form's only field, 
	   Body, are stored in a new Page. The save() method is then called to write the data to a file, 
	   and the client is redirected to the /view/ page.
	*/
	//title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")

	/* The value returned by FormValue is of type string. 
	   We must convert that value to []byte before it will fit into the Page struct.
	   We use []byte(body) to perform the conversion.
	*/
	p := &Page{
		Title: title,
		Body: []byte(body),
	}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/" + title, http.StatusFound)
}


/* 
The function template.ParseFiles will read the contents of edit.html,  
view.html and return a *template.Template.
-
The ParseFiles function takes any number of string arguments that identify our template files, 
and parses those files into templates that are named after the base file name.
*/
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w,tmpl + ".html", p) 

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}	
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

/* 
Catching the error condition in each handler introduces a lot of repeated code. 
What if we could wrap each of the handlers in a function that does this validation and error checking?
*/
func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	/*
	The function handler is of the type http.HandlerFunc. It takes an http.ResponseWriter and an http.Request as its arguments.

	An http.ResponseWriter value assembles the HTTP server's response; by writing to it, we send data to the HTTP client.

	An http.Request is a data structure that represents the client HTTP request. r.URL.Path is the path component of the request URL. 
	The trailing [1:] means "create a sub-slice of Path from the 1st character to the end." This drops the leading "/" from the path name.
	*/

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}



