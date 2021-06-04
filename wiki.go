package main

import (
	// "errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))

type Page struct{
	Title string
	Body []byte
}

// The Page struct describes how page data will be stored in memory. But what about persistent storage? We can address that by creating a save method on Page:
func (p *Page) save() error{
	filename := p.Title +".txt"
	return ioutil.WriteFile("data/"+filename, p.Body, 0600)
}

// func getTitle(w http.ResponseWriter, r *http.Request)(string, error){
// 	m := validPath.FindStringSubmatch(r.URL.Path)
// 	if m == nil {
// 		http.NotFound(w, r)
// 		return "", errors.New("invalid Page Title")
// 	}
// 	return m[2], nil // The title is the second subexpression.
// }

func loadPage(title string) (*Page, error){
	filename := title + ".txt"
	body, err := ioutil.ReadFile("data/"+filename)
	if err != nil{
		return nil, err
	}
	return &Page{Title: title, Body:body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page){
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	// t, err := template.ParseFiles(tmpl+ ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		// return
	}
	// t.Execute(w, p)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string){
	// title, err := getTitle(w,r)
	// if err != nil {
	// 	return
	// }
	// title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
	renderTemplate(w,"view",p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string){
	// title := r.URL.Path[len("/edit/"):]
	// title, err := getTitle(w, r)
    // if err != nil {
    //     return
    // }
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w,"edit",p)
	// t, _ := template.ParseFiles("edit.html")
	// t.Execute(w, p)

	// fmt.Fprintf(w, "<h1>Editing %s</h1>"+
	// 	"<form action=\"/save/%s\" method=\"POST\">" +
	// 	"<textarea name=\"body\">%s</textarea><br>"+
	// 	"<input type=\"submit\" value=\"Save\">"+
    //     "</form>",
	// 	p.Title, p.Title, p.Body)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string){
	// title := r.URL.Path[len("/save/"):]
	// title, err := getTitle(w, r)
    // if err != nil {
    //     return
    // }
	body := r.FormValue("body")
	p := &Page{Title:title, Body: []byte(body)}
	err := p.save()
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+ title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request)  {
		// here we will extrat the page title from the request,
		// and call the provided hander 'fn'
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main(){
	// p1 := &Page{Title: "TestPage", Body: []byte("This is a sample page")}
	// p1.save()
	// p2, _ := loadPage("TestPage")
	// fmt.Println(string(p2.Body))
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}