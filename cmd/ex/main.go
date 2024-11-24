package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
	"net/http"
)

func main() {

	a := app.New()

	w := a.NewWindow("Find Email of organization employee...")

	w.Resize(fyne.NewSize(600, 400))

	title1 := widget.NewLabel("Email Finder Verification App")

	title1.Alignment = fyne.TextAlignCenter

	title1.TextStyle = fyne.TextStyle{Bold: true}

	// Get url for user

	long_url := widget.NewEntry()

	long_url.SetPlaceHolder("Enter URL here...")

	long_url.Text = "intercom.io"

	result_data := widget.NewMultiLineEntry()

	result_data.MultiLine = true

	result_data.Resize(fyne.NewSize(300, 400))

	result_data.Move(fyne.NewPos(0, 0))

	result_data.SetPlaceHolder("First name &amp; Email will displayed here...")

	result_container := container.NewWithoutLayout(result_data)

	result_container.Resize(fyne.NewSize(300, 400))

	email_place_holder := container.NewVBox(widget.NewLabel("List of Emails:"))

	email_place_holder2 := container.NewHScroll(email_place_holder)

	btn := widget.NewButton("Find", func() {

		API_KEY := "d7888147f3ec5a63f1f6831f7357b229b90448c5"

		// Raw URL

		// https://api.hunter.io/v2/domain-search?domain=intercom.io&amp;api_key=API_KEY

		// url of api

		URL := fmt.Sprintf("https://api.hunter.io/v2/domain-search?domain=%s&amp;api_key=%s", long_url.Text, API_KEY)

		//URL := "http://127.0.0.1:3000/email"

		//URL := "http://127.0.0.1:3000/email2"

		// http response

		resp, _ := http.Get(URL)

		defer resp.Body.Close()

		// read the response body

		body, _ := ioutil.ReadAll(resp.Body)

		// parse json data

		emailFinder, _ := UnmarshalEmailFinder(body)

		// add the suffix received to main url

		fmt.Print(len(emailFinder.Data.Emails))

		var emailSlice string

		for i := 0; i < len(emailFinder.Data.Emails); i++ {

			var v *string

			if emailFinder.Data.Emails[i].FirstName != nil {

				v = emailFinder.Data.Emails[i].FirstName

				fmt.Println(string(*v))

				f := fmt.Sprintf("\n firstname:%s\n email %d: %s\n", *v, i+1, emailFinder.Data.Emails[i].Value)

				email_ := fmt.Sprintf("%d# %s\nemail : %s", i+1, *v, emailFinder.Data.Emails[i].Value)

				email_place_holder.Add(widget.NewLabel(email_))

				emailSlice += f

			}

		}

		result_data.Text = "Organization Name: " + emailFinder.Data.Organization + "\n" + emailSlice

		result_data.Refresh()

	})

	c := container.NewVBox(title1, long_url, btn, result_container)

	split1 := container.NewHSplit(c, email_place_holder2)

	w.SetContent(split1)

	w.ShowAndRun()

}

// This file was generated from JSON Schema using quicktype, do not modify it directly.

// To parse and unparse this JSON data, add this code to your project and do:

//

//    emailFinder, err := UnmarshalEmailFinder(bytes)

//    bytes, err = emailFinder.Marshal()

func UnmarshalEmailFinder(data []byte) (EmailFinder, error) {

	var r EmailFinder

	err := json.Unmarshal(data, &r)

	return r, err

}

func (r *EmailFinder) Marshal() ([]byte, error) {

	return json.Marshal(r)

}

type EmailFinder struct {
	Data Data `json:"data"`

	Meta Meta `json:"meta"`
}

type Data struct {
	Domain string `json:"domain"`

	Disposable bool `json:"disposable"`

	Webmail bool `json:"webmail"`

	AcceptAll bool `json:"accept_all"`

	Pattern string `json:"pattern"`

	Organization string `json:"organization"`

	Country interface{} `json:"country"`

	State interface{} `json:"state"`

	Emails []Email `json:"emails"`
}

type Email struct {
	Value string `json:"value"`

	Type Type `json:"type"`

	Confidence int64 `json:"confidence"`

	Sources []Source `json:"sources"`

	FirstName *string `json:"first_name"`

	LastName *string `json:"last_name"`

	Position *string `json:"position"`

	Seniority *string `json:"seniority"`

	Department *string `json:"department"`

	Linkedin interface{} `json:"linkedin"`

	Twitter *string `json:"twitter"`

	PhoneNumber *string `json:"phone_number"`

	Verification Verification `json:"verification"`
}

type Source struct {
	Domain string `json:"domain"`

	URI string `json:"uri"`

	ExtractedOn string `json:"extracted_on"`

	LastSeenOn string `json:"last_seen_on"`

	StillOnPage bool `json:"still_on_page"`
}

type Verification struct {
	Date *string `json:"date"`

	Status *string `json:"status"`
}

type Meta struct {
	Results int64 `json:"results"`

	Limit int64 `json:"limit"`

	Offset int64 `json:"offset"`

	Params Params `json:"params"`
}

type Params struct {
	Domain string `json:"domain"`

	Company interface{} `json:"company"`

	Type interface{} `json:"type"`

	Seniority interface{} `json:"seniority"`

	Department interface{} `json:"department"`
}

type Type string

const (
	Generic Type = "generic"

	Personal Type = "personal"
)
