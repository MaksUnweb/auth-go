package handlers

import (
	"html/template"
	"net/http"
)


func LoginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)	
	}
	
	some_template, err := template.ParseFiles("./templates/login.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	err = some_template.Execute(w, nil)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
}
