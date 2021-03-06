# Web_Template_Example
I have always built small websites in Apache/PHP and was intrigued at the notion that we could get a standard full function web server built in very few lines of code.   This code is portable and works on any computer (Windows, Mac, Unix veriants), and depends only on GO Language.  

This is a small test program written entirely in Go Language which demonstrates launching of a basic web server, and common functionality. I wanted to build a small site which deals with all the major pieces needed for a functional user site.  It handles Templates and also deals with reading and parsing json files.  For testing purposes, I'm not logging data to a file, just lots of debug information to the console.  I suppose i could use my logging library, but I didn't think it was necessary here. 

Specifically this prorotype does the following: 

 *   Launches a basic web server on port 8000.
 *   It has multiple pages, including login/logout functionality using cookies
 *   The Login/Logout functionality is basic for the sake of this exercise.  It uses local variables to check login credentials but could easily be modified to add the users in a database - but that is a different exercise.
 *   This prototype utilizes templates and logic is built into the template pages. 
 *   The "internal.html" page requires to be logged in prior to displaying data from it. 
 *   This prototype uses Gorilla libraries for secure cookies
 *   The prototype addresses 404 not found pages properly, and handles "/" and "/index.html" as the home page. 
 *   The prototype is coded to deal with favicon.ico when running as http server
 *   This prototype reads a json file and uses the data on the "asic" pages

TODO: 
 *  Add Dashboard Page
 *  Add Crud Capability (be logged in to delete/update data, but view from everyone)

# Requirements and Building
````
go get github.com/gorilla/securecookie
go build main.go
````

# Basic Testing and sample outputs
The output on the console is fairly verbose -I added lots of debug line showing what is done as each transation is hit
This was stictly a learning exercise. 

To view this, bring up a web browser as follows: 
```` 
http://localhost:8000
````

# Login Credentials
The credentials are coded in the main.go code, but are as follows
```` 
username: howie
password: 123
````

# Sample Console Output
This is what the testing on a simple windows box looks like.  I clicked on each page from the website, logged in and tried a bad login combination.  
````
C:\Users\grapekh\go\src\simple_go_httpd_template_test>main.exe
Test Website for single page with templates... launch http://localhost:8000
Grabbing file: asic.html
Grabbing file: asic2.html
Grabbing file: error.html
Grabbing file: favicon.ico
Grabbing file: favicon2.ico
Grabbing file: footer.html
Grabbing file: header.html
Grabbing file: index.html
Grabbing file: internal.html
Grabbing file: login.html
allFiles is: [./www/asic.html ./www/asic2.html ./www/error.html ./www/footer.html ./www/header.html ./www/index.html ./www/internal.
html ./www/login.html]
Got config data, Greeting: Howdy, Username: admin, DeviceModel: SC1|DCR1
Handling index.html page.  indexPageData = &{Home Page Friday Howdy}
Handling asic.html page.  asicPageData = &{Asic Page Friday Howdy admin SC1|DCR1}
Handling asic2.html page.  asic2PageData = &{Asic2 Page Friday Howdy admin SC1|DCR1}
Found a page for index which was not index... 404 should be sent.
Here we are in renderErrorPage: error message is: Error 404: Page /badpage.html not found
Cookie is not set... user NOT logged in
internalPageHandler: We are not logged in yet.
Handling internal.html page.  internalPageData = &{Internal Page false }
Some debug -- method = GET
Cookie is not set... user NOT logged in
in handleLoginPage: We are not logged in yet... setting session and continuing.
Some debug -- method = POST
Got login request... checking variables in the form
LoginHandler here... creds entered are: name: (howie grapek), pass: (123)
Username/Pass entered does not match required credentials
Some debug -- method = POST
Got login request... checking variables in the form
LoginHandler here... creds entered are: name: (howie), pass: (123)
Username/Password entered matched required credentials... set session, etc. continuing
Cookie is set.... user IS logged in
internalPageHandler: We are logged in... coookie Username is: (howie) - go and display the internal page
Handling internal.html page.  internalPageData = &{Internal Page true howie}
Got logout request... clearing cookies and redirecting to index page
Cookie is set.... user IS logged in
Handling index.html page.  indexPageData = &{Home Page Friday Howdy}
````
