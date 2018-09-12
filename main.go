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
 * 		Add Login Page and cookies using gorilla Apps. 
 *		Add Dashboard Page, with a couple of protected pages
 *		Add CRUD Capability. 
 * 
 */

package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"
	"log"
	"fmt"
	"strings"
)

// Generic global variables to be used to manage template data for any html page. 

var t *template.Template 		// template data built to be used in the page display 
var c config 					// global variaable holding the json data from "config.json" file  (see struct below)

// define the global data structures containing all variables which will go into the html pages:
type indexPageData struct {
	PageTitle string
	Weekday  time.Weekday
	Greeting string
}

type asicPageData struct {
    PageTitle string
	Weekday time.Weekday
	Greeting string
	Username string
	DeviceModel string
}

type asic2PageData struct {
    PageTitle string
	Weekday time.Weekday
	Greeting string
	Username string
	DeviceModel string
}

// greeting structure type definition. matches data in config.json
type config struct {
	Greeting string
	Username string
	DeviceModel string
}

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
	http.HandleFunc("/asic.html",  handleAsicPage)
	http.HandleFunc("/asic2.html",  handleAsic2Page)

	// Start up the webserver 
	log.Fatal(http.ListenAndServe(":8000", nil))
}

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
    t, err = template.ParseFiles(allFiles...) 	// parses all .tmpl or ./html files in the 'www' folder

}

// Get raw data from the json file called "config.json"
// Put all data into the configuration structure global variable called "c".  
func loadConfigJson() {
	contentBytes, _ := ioutil.ReadFile("config.json")
	json.Unmarshal(contentBytes, &c)

	fmt.Printf("Got config data, Greeting: %s, Username: %s, DeviceModel: %s\n", c.Greeting, c.Username, c.DeviceModel)
}

// Display the named template
func display(w http.ResponseWriter, tmpl string, data interface{}) {

	// grab the appropriate template (passed in) and display the data accordingly using the template
	err := t.ExecuteTemplate(w, tmpl, data)

	// we should really display an error html page but this is fine - display on terminal and die. 
	if (err != nil) {
		    log.Println(err)
	}
}


// HANDLERS - One for each page. 

func handleIndexPage(w http.ResponseWriter, r *http.Request) {
	// fill the structure template with actual values calculated and read in accordingly. 
	PageTitle := "Home Page"
	indexPageData := &indexPageData{PageTitle, time.Now().Weekday(), c.Greeting}
	fmt.Printf("Handling index.html page.  indexPageData = %v\n", indexPageData)

	display(w, "index.html", &indexPageData)		
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