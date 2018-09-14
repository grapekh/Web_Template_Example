/* Copyright (c) 2018 by Howard I Grapek <howiegrapek@yahoo.com>
 * All rights reserved.
 *
 * License:
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *   - Redistributions of source code must retain the above copyright notice, this
 *     list of conditions and the following disclaimer.
 *
 *   - Redistributions in binary form must reproduce the above copyright notice,
 *     this list of conditions and the following disclaimer in the documentation
 *     and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * Description:
 *   This is a small test program which demonstrates launching of a basic web server on port 8000.
 *   It has a main index page and just a couple of secondary pages.
 *   This test shows how to use templates and variables on each page.
 *
 * 	TODO:
 * 		
 *		Add Dashboard Page
 *		Add CRUD Capability.
 *
 *
 * Version 1.0 - Grapek - 20180904 - Original
 * Version 1.1 - Grapek - 20180905 - Added 404 not found check for index handler and favicon handler
 * Version 1.2 - Grapek - 20180806 - Added Login Page, protected "internal" page,  and secure cookies using gorilla Apps.
 */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"github.com/gorilla/securecookie"
)

// Generic global variables to be used to manage template data for any html page.

var t *template.Template // template data built to be used in the page display
var c config             // global variaable holding the json data from "config.json" file  (see struct below)

// define the global data structures containing all variables which will go into the html pages:
type indexPageData struct {
	PageTitle string
	Weekday   time.Weekday
	Greeting  string
}

type asicPageData struct {
	PageTitle   string
	Weekday     time.Weekday
	Greeting    string
	Username    string
	DeviceModel string
}

type asic2PageData struct {
	PageTitle   string
	Weekday     time.Weekday
	Greeting    string
	Username    string
	DeviceModel string
}

type errorPageData struct {
	PageTitle string
	ErrorMsg  string
}

type internalPageData struct {
	PageTitle string
	LoggedIn bool
	Username  string
}

type loginPageData struct {
	PageTitle string
	Username string
	LoggedIn bool
}

// greeting structure type definition. matches data in config.json
type config struct {
	Greeting    string
	Username    string
	DeviceModel string
}

// Cookie stuff: 
var cookieHandler = securecookie.New(
    securecookie.GenerateRandomKey(64),
    securecookie.GenerateRandomKey(32))

func main() {
	fmt.Println("Test Website for single page with templates... launch http://localhost:8000")

	// load in all of the template information for the webs living in ./www
	loadTemplates()

	// Get the json data from config.json and load it into memory
	loadConfigJson()

	// initializeTemplates()			// really only need to do this outside if pages are not templates.

	// Manage the pages.
	http.HandleFunc("/", handleIndexPage)
	http.HandleFunc("/index.html", handleIndexPage)
	http.HandleFunc("/asic.html", handleAsicPage)
	http.HandleFunc("/asic2.html", handleAsic2Page)

	// Intenal pages (requiring login)  
	http.HandleFunc("/login.html", handleLoginPage)
	http.HandleFunc("/logout.html", handleLogoutPage)
	http.HandleFunc("/internal.html", handleInternalPage)

	http.HandleFunc("/favicon.ico", handleFavicon)

	// Start up the webserver
	log.Fatal(http.ListenAndServe(":8000", nil))
}

// Template functions. 

// load all html templates into memory. 
func loadTemplates() {
	// Read in all the templates. (./www/*html)
	var allFiles []string
	files, err := ioutil.ReadDir("./www")
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		filename := file.Name()
		fmt.Printf("Grabbing file: %s\n", filename)
		if strings.HasSuffix(filename, ".html") {
			allFiles = append(allFiles, "./www/"+filename)
		}
	}

	fmt.Printf("allFiles is: %v\n", allFiles)
	t, err = template.ParseFiles(allFiles...) // parses all .tmpl or ./html files in the 'www' folder
}

// Display the named template
func display(w http.ResponseWriter, tmpl string, data interface{}) {

	// grab the appropriate template (passed in) and display the data accordingly using the template
	err := t.ExecuteTemplate(w, tmpl, data)

	// we should really display an error html page but this is fine - display on terminal and die.
	if err != nil {
		fmt.Println("In display: got an error: =%v\n", err)
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		renderErrorPage(w, err)
	}
}

// Get raw data from the json file called "config.json"
// Put all data into the configuration structure global variable called "c".
func loadConfigJson() {
	contentBytes, _ := ioutil.ReadFile("config.json")
	json.Unmarshal(contentBytes, &c)

	fmt.Printf("Got config data, Greeting: %s, Username: %s, DeviceModel: %s\n", c.Greeting, c.Username, c.DeviceModel)
}

// 
// Secure Cookie / Session Functions
//

func getUserName(request *http.Request) (userName string) {
    if cookie, err := request.Cookie("session"); err == nil {
        cookieValue := make(map[string]string)
        if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
            userName = cookieValue["name"]
        }
    }
    return userName
}

func checkLoggedIn(request *http.Request) (state bool) {
       if _, err := request.Cookie("session"); err == nil {
            fmt.Printf("Cookie is set.... user IS logged in\n")
            state = true
        } else {
            fmt.Printf("Cookie is not set... user NOT logged in\n")
            state = false
        }

        return state
}

func setSession(userName string, response http.ResponseWriter) {
    value := map[string]string{
        "name": userName,
    }
    if encoded, err := cookieHandler.Encode("session", value); err == nil {
        cookie := &http.Cookie{
            Name:  "session",
            Value: encoded,
            Path:  "/",
        }
        http.SetCookie(response, cookie)
    }
}

func clearSession(response http.ResponseWriter) {
    cookie := &http.Cookie{
        Name:   "session",
        Value:  "",
        Path:   "/",
        MaxAge: -1,
    }
    http.SetCookie(response, cookie)
}

// 
// HTML PAGE HANDLERS - One for each page.
//

func handleIndexPage(w http.ResponseWriter, r *http.Request) {
	// fill the structure template with actual values calculated and read in accordingly.

	// This on eis pretty simple -
	// check to see if the page is "/" or "/index.html"
	// Anything else is an error.
	if r.URL.Path != "/" && r.URL.Path != "/index.html" {
		fmt.Println("Found a page for index which was not index... 404 should be sent. ")

		// Build an actual error message...
		message := fmt.Sprintf("Error 404: Page %s not found", r.URL.Path)
		errorMessage := errors.New(message)

		// Call Render Error  Page with error message.
		renderErrorPage(w, errorMessage)

		return
	}

	PageTitle := "Home Page"
	indexPageData := &indexPageData{PageTitle, time.Now().Weekday(), c.Greeting}
	fmt.Printf("Handling index.html page.  indexPageData = %v\n", indexPageData)

	display(w, "index.html", &indexPageData)
}

func handleLoginPage(w http.ResponseWriter, request *http.Request) {
	PageTitle := "Login Page"
	loggedIn := false					// the default value, not logged in yuet. 
	var userName string 				// Placeholder for username in cookie (or typed in)

    // Simple validation here... 
    // Verify username and password. 

    fmt.Printf("\n\nSome debug -- method = %v\n", request.Method)

    
    // Handle this page as GET request
    if (request.Method == "GET") {
        if (checkLoggedIn(request)) {
            fmt.Printf("in handleLoginPage: We are already logged in.\n")
            // since we are already logged in, we can display the login page with the fact that we are logged in already
            // and return (stop processing)
            userName = getUserName(request)
            loggedIn = true
            
        } else {
            fmt.Printf("in handleLoginPage: We are not logged in yet... setting session and continuing. \n")
            // we dont' know if we are not logged in because no cookie or user is incorrect yet. 
        }

        lpd := &loginPageData{PageTitle, userName, loggedIn}
        display(w, "login.html", &lpd)
        return
    }

  
  	// What remains is the logic when the form is submitted (POST)
    fmt.Printf("Got login POST request... checking variables in the form\n")
    name := request.FormValue("name")
    pass := request.FormValue("password")

    // check get or post. ... get returns nothing, post will fill the form variables. 
    fmt.Printf("LoginHandler here... creds entered are: name: (%s), pass: (%s)\n", name, pass)

    //
    // check validity of login... (We can get these from a database or whatever - 
    // for testing, just use howie/123)
    //
    if (name == "howie" && pass == "123") {
    	fmt.Printf("Username/Password entered matched required credentials... set session, etc. continuing\n")

    	// set the session secure cookie and display login success message. 
    	setSession(name, w)
    	loggedIn = true
    } else {
    	fmt.Printf("Username/Pass entered does not match required credentials\n")
    	// display the login failure message and try again. (Logic is on the )
    	// Display the login page... 
    }

    userName = name
    lpd := &loginPageData{PageTitle, userName, loggedIn}
   	display(w, "login.html", &lpd)
    return
}

func handleLogoutPage(response http.ResponseWriter, request *http.Request) {
    fmt.Printf("Got logout request... clearing cookies if they are set and redirecting back to index page\n")

    // Only need to clear the session if we are logged in. 
     if (!checkLoggedIn(request)) {
        fmt.Printf("logoutHandler: We are not logged in, no need to clear cookie - bounce back to index.\n")
        http.Redirect(response, request, "/", 302)
        return
    }

    // we are currently logged in, so clear session and return to home page. 
    clearSession(response)
    http.Redirect(response, request, "/", 302)
}


func handleInternalPage(w http.ResponseWriter, request *http.Request) {
	var loggedIn bool		// This page is one which requires to be logged in first.

	//userName := "How Grap"
	PageTitle := "Internal Page"
	userName := getUserName(request)

	if (!checkLoggedIn(request)) {
        fmt.Printf("internalPageHandler: We are not logged in yet.\n")
        loggedIn = false
    } else {
    	fmt.Printf("internalPageHandler: We are logged in... coookie Username is: (%s) - go and display the internal page\n", userName)
		loggedIn = true
    }

	internalPageData := &internalPageData{PageTitle, loggedIn, userName}
	fmt.Printf("Handling internal.html page.  internalPageData = %v\n", internalPageData)
	display(w, "internal.html", &internalPageData)
}


func handleAsicPage(w http.ResponseWriter, r *http.Request) {

	// fill the structure template with actual values calculated and read in accordingly.
	PageTitle := "Asic Page"
	asicPageData := &asicPageData{PageTitle, time.Now().Weekday(), c.Greeting, c.Username, c.DeviceModel}
	fmt.Printf("Handling asic.html page.  asicPageData = %v\n", asicPageData)

	display(w, "asic.html", &asicPageData)
}

func handleAsic2Page(w http.ResponseWriter, r *http.Request) {

	// fill the structure template with actual values calculated and read in accordingly.
	PageTitle := "Asic2 Page"
	asic2PageData := &asic2PageData{PageTitle, time.Now().Weekday(), c.Greeting, c.Username, c.DeviceModel}
	fmt.Printf("Handling asic2.html page.  asic2PageData = %v\n", asic2PageData)

	display(w, "asic2.html", &asic2PageData)
}

func handleFavicon(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Handling favicon.ico...\n")
	w.Header().Set("Content-Type", "image/x-icon")
	http.ServeFile(w, r, "www/favicon.ico")
}

func renderErrorPage(w http.ResponseWriter, errorMsg error) {
	fmt.Printf("Here we are in renderErrorPage: error message is: %v\n", errorMsg)

	PageTitle := "Error"
	errorPageData := &errorPageData{PageTitle, errorMsg.Error()}

	display(w, "error.html", &errorPageData)
}


