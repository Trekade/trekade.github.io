package handler

import (
	"encoding/json"
	"net/http"

	"github.com/keighl/postmark"
)

//Feedback is feedback struct
type Feedback struct {
	Name    string
	Email   string
	Msg 	string
}

//Handler is the default handler
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/sendmail" || r.Method != "POST" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	var fb Feedback
	err := json.NewDecoder(r.Body).Decode(&fb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if fb.Name == "" || fb.Email == "" || fb.Msg == "" {
		http.Error(w, "Not all fields filled out", http.StatusBadRequest)
		return
	}

	res, body, err := SendMail(fb)
	if err != nil {
		println("Error sending Email: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(res)
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(body))
	return
}

func SendMail(f Feedback) (res int, out string, err error) {
	client := postmark.NewClient("e92bacac-f491-4a18-aafb-652cc31b0790", "")
	email := postmark.Email{
		From: "no-reply@trekade.com",
		To: "daniel@trekade.com",
		Subject: "[trekade.com] Contact Form",
		HtmlBody: f.Msg+"<br>"+f.Name+"<br>"+f.Email,
	    TextBody: f.Msg+"\n"+f.Name+"\n"+f.Email,
	}

	_, err = client.SendEmail(email)
	if err != nil {
		return 400, "", err
	}

	return 200, "", nil
}