package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox.subh.am/internal/models"
	"strconv"
	"strings"
	"unicode/utf8"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets
	// Pass the data to the render() helper as normal.
	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, r, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "create.tmpl", data)
}
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w,http.StatusBadRequest)
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires ,err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w,http.StatusBadRequest)
	}
	fieldsErrors := make(map[string]string)

	if strings.TrimSpace(title) == ""{
		fieldsErrors["title"] = "this field cannot be blanked"
	}else if utf8.RuneCountInString(title) > 100 {
		fieldsErrors["title"] ="this fields cannot be more than 100 characters long"
	}

	if strings.TrimSpace(content) == ""{
		fieldsErrors["content"] = "this field cannnot be more than 100 chars"
	}

	if expires != 1 && expires != 7 && expires != 345 {
		fieldsErrors["expires"] = "this should be equal 1,7,or 365"
	}
	
	if len(fieldsErrors )> 0 {
		fmt.Fprint(w,fieldsErrors)
		return
	}


	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
