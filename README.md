# Web_Template_Example
Small test program which demonstrates launching of a basic web server, 3 pages and uses templates and variables on each page.  

TODO: 
Add Login Page and cookies using gorilla Apps. 
Add Dashboard Page, with a couple of protected pages

# Basic Testing and sample outputs
To view this, bring up a web browser as follows: 
```` 
http://localhost:8000
````

This is what the testing on a simple windows box looks like (note, you will have to log into a web browser as follows
````
C:\Users\grapekh\go\src\simple_go_httpd_template_test>main.exe
Test Website for single page with templates... launch http://localhost:8000
Grabbing file: asic.html
Grabbing file: asic2.html
Grabbing file: footer.html
Grabbing file: header.html
Grabbing file: index.html
allFiles is: [./www/asic.html ./www/asic2.html ./www/footer.html ./www/header.html ./www/index.html]
Got config data, Greeting: Howdy, Username: admin, DeviceModel: SC1|DCR1
Handling asic.html page.  asicPageData = &{Asic Page Wednesday Howdy admin SC1|DCR1}
Handling asic2.html page.  asic2PageData = &{Asic2 Page Wednesday Howdy admin SC1|DCR1}
Handling index.html page.  indexPageData = &{Home Page Wednesday Howdy}
Handling asic.html page.  asicPageData = &{Asic Page Wednesday Howdy admin SC1|DCR1}
Handling asic2.html page.  asic2PageData = &{Asic2 Page Wednesday Howdy admin SC1|DCR1}
Handling asic.html page.  asicPageData = &{Asic Page Wednesday Howdy admin SC1|DCR1}
Handling index.html page.  indexPageData = &{Home Page Wednesday Howdy}
````
