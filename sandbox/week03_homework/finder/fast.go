package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

// FastSearch should be 20% (at least) better (resource-wise, faster) than SlowSearch.
// Check: `go test -v $module; go test -bench . -benchmem $module`
func FastSearch(out io.Writer) {
	defer panic(`
	Task outline:
	For each line from file:
	- transform json => Record, where Record is (list-of-browsers, user-name, user-email)
	- loop over browsers: update global-set-of-unique-browsers, set MSIE-on-Android flag
	- if MSIE-on-Android: update result (print to out)
	print-to-out: len(global-set-of-unique-browsers)
	`)
}

func FastSearch_ExplainedSlowSearch(out io.Writer) {
	/*
		   file lines sample:
			{
			"browsers":[
				"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36",
				"LG-LX550 AU-MIC-LX550/2.0 MMP/2.0 Profile/MIDP-2.0 Configuration/CLDC-1.1",
				"Mozilla/5.0 (Android; Linux armv7l; rv:10.0.1) Gecko/20100101 Firefox/10.0.1 Fennec/10.0.1",
				"Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; MATBJS; rv:11.0) like Gecko"
			],
			"company":"Flashpoint",
			"country":"Dominican Republic",
			"email":"JonathanMorris@Muxo.edu",
			"job":"Programmer Analyst #{N}",
			"name":"Sharon Crawford",
			"phone":"176-88-49"
			}
		   {"browsers":["Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36","LG-LX550 AU-MIC-LX550/2.0 MMP/2.0 Profile/MIDP-2.0 Configuration/CLDC-1.1","Mozilla/5.0 (Android; Linux armv7l; rv:10.0.1) Gecko/20100101 Firefox/10.0.1 Fennec/10.0.1","Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; MATBJS; rv:11.0) like Gecko"],"company":"Flashpoint","country":"Dominican Republic","email":"JonathanMorris@Muxo.edu","job":"Programmer Analyst #{N}","name":"Sharon Crawford","phone":"176-88-49"}
		   {"browsers":["Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.10240","Mozilla/4.0 (compatible; GoogleToolbar 4.0.1019.5266-big; Windows XP 5.1; MSIE 6.0.2900.2180)","SAMSUNG-S8000/S8000XXIF3 SHP/VPP/R5 Jasmine/1.0 Nextreaming SMM-MMS/1.2.0 profile/MIDP-2.1 configuration/CLDC-1.1 FirePHP/0.3","Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.8.1) Gecko/20061024 Firefox/2.0 (Swiftfox)"],"company":"Youfeed","country":"Uruguay","email":"gLong@Zoombox.org","job":"Automation Specialist #{N}","name":"Michael Davis","phone":"946-49-01"}
		   {"browsers":["Mozilla/5.0 (Linux; Android 5.1.1; Nexus 7 Build/LMY47V) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.78 Safari/537.36 OPR/30.0.1856.93524","Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; Trident/6.0)","Mozilla/5.0 (Windows; U; Windows NT 5.2; en-US) AppleWebKit/532.9 (KHTML, like Gecko) Chrome/5.0.310.0 Safari/532.9","Mozilla/5.0 (compatible; Yahoo! Slurp China; http://misc.yahoo.com.cn/help.html)"],"company":"Feedbug","country":"Germany","email":"yKing@Leexo.net","job":"Cost Accountant","name":"Melissa Price","phone":"3-981-961-68-80"}
		   {"browsers":["SonyEricssonW580i/R6BC Browser/NetFront/3.3 Profile/MIDP-2.0 Configuration/CLDC-1.1","Mozilla/5.0 (Windows NT 6.2; ARM; Trident/7.0; Touch; rv:11.0; WPDesktop; NOKIA; Lumia 920) like Geckoo","Mozilla/5.0 (Linux; U; Android 2.0; en-us; Milestone Build/ SHOLS_U2_01.03.1) AppleWebKit/530.17 (KHTML, like Gecko) Version/4.0 Mobile Safari/530.17","Mozilla/5.0 (X11; Linux i686) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/51.0.2704.79 Chrome/51.0.2704.79 Safari/537.36"],"company":"Wordware","country":"Equatorial Guinea","email":"est_ratione_maxime@Thoughtsphere.gov","job":"Electrical Engineer","name":"Maria Watkins","phone":"5-525-472-31-88"}
		   {"browsers":["Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/537.78.1 (KHTML like Gecko) Version/7.0.6 Safari/537.78.1","Mozilla/1.22 (compatible; MSIE 5.01; PalmOS 3.0) EudoraWeb 2.1","Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2049.0 Safari/537.36","Mozilla/5.0 (SymbianOS/9.4; U; Series60/5.0 SonyEricssonP100/01; Profile/MIDP-2.1 Configuration/CLDC-1.1) AppleWebKit/525 (KHTML, like Gecko) Version/3.0 Safari/525"],"company":"Flashspan","country":"Turkmenistan","email":"est_adipisci@Brightbean.name","job":"Research Assistant #{N}","name":"Johnny Jones","phone":"3-553-621-66-19"}
		   {"browsers":["Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0; .NET CLR 2.0.50727; .NET CLR 3.0.04506.648; .NET CLR 3.5.21022; .NET CLR 3.0.4506.2152; .NET CLR 3.5.30729)","Mozilla/5.0 (iPod; U; CPU iPhone OS 2_2_1 like Mac OS X; en-us) AppleWebKit/525.18.1 (KHTML, like Gecko) Version/3.1.1 Mobile/5H11a Safari/525.20","Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9.1.2) Gecko/20090803 Ubuntu/9.04 (jaunty) Shiretoko/3.5.2","Mozilla/5.0 (X11; Linux x86_64; rv:15.0) Gecko/20120724 Debian Iceweasel/15.02"],"company":"Jazzy","country":"Bahamas","email":"sunt@Dynabox.name","job":"Senior Quality Engineer","name":"Anna Perkins","phone":"494-01-38"}
		   {"browsers":["Opera/9.60 (J2ME/MIDP; Opera Mini/4.2.14320/554; U; cs) Presto/2.2.0","Mozilla/5.0 (Linux; Android 4.4; Nexus 5 Build/BuildID) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/30.0.0.0 Mobile Safari/537.36","Mozilla/5.0 (X11; FreeBSD i386; rv:28.0) Gecko/20100101 Firefox/28.0 SeaMonkey/2.25","Mozilla/5.0 (iPad; CPU OS 10_0 like Mac OS X) AppleWebKit/601.1 (KHTML, like Gecko) CriOS/49.0.2623.109 Mobile/14A5335b Safari/601.1.46"],"company":"Jabberbean","country":"Vietnam","email":"wRamirez@Skajo.com","job":"Social Worker","name":"John Flores","phone":"9-257-322-49-96"}
		   {"browsers":["Opera/9.80 (X11; Linux i686) Presto/2.12.388 Version/12.16","Mozilla/5.0 (Linux; Android 5.1.1; Nexus 7 Build/LMY47V) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.78 Safari/537.36 OPR/30.0.1856.93524","Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_5_6; en-US) AppleWebKit/528.16 (KHTML, like Gecko, Safari/528.16) OmniWeb/v622.8.0","Opera/7.50 (Windows ME; U) [en]"],"company":"Livetube","country":"Palau","email":"qui@Meeveo.com","job":"Environmental Specialist","name":"Peter Palmer","phone":"7-732-179-16-86"}
	*/
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	fileContentsBytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(fileContentsBytes), "\n")

	atSymbolRegexp := regexp.MustCompile("@")
	seenBrowsers := []string{}
	uniqueBrowsers := 0
	foundUsers := ""

	users := make([]map[string]any, 0) // all lines, each parsed as map (string => any)

	for _, line := range lines {
		user := make(map[string]any)
		err := json.Unmarshal([]byte(line), &user)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}

	for i, user := range users {

		isAndroid := false
		isMSIE := false

		browsers, ok := user["browsers"].([]any)
		if !ok {
			log.Printf("Type assertion for `browsers` has failed, idx: %d\n", i)
			continue
		}

		for j, browserAny := range browsers { // loop on browsers #1
			browser, ok := browserAny.(string)
			if !ok {
				log.Printf("Type assertion for browser has failed, idxs: %d, %d\n", i, j)
				continue
			}

			// find `Android` in current browser, check uniqueness
			if ok, err := regexp.MatchString("Android", browser); ok && err == nil {
				isAndroid = true

				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			} // one browser processed.
		} // all user browsers processed (unique Android browsers).

		for j, browserAny := range browsers { // loop on browsers #2
			browser, ok := browserAny.(string)
			if !ok {
				log.Printf("Type assertion for browser has failed, idxs: %d, %d\n", i, j)
				continue
			}

			// find `MSIE` in current browser, check uniqueness
			if ok, err := regexp.MatchString("MSIE", browser); ok && err == nil {
				isMSIE = true

				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			} // one browser processed.
		} // all user browsers processed (unique MSIE browsers).

		// get result for current user (line)
		if isMSIE && isAndroid {
			// log.Println("MSIE on Android user:", user["name"], user["email"])
			email := atSymbolRegexp.ReplaceAllString(user["email"].(string), " [at] ") // replace `@` with ` [at] `
			// add line [42] Foo Bar <foo [at] some.domain>
			// should print it to `out`, no need for collecting
			foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user["name"], email)
		} else {
			continue
		}
	} // end of users (lines) loop

	// results (set of unique browsers; set of users where user(idx, name, email)-with-selected-browser):

	// line No, name, email-transformed (for each filtered user (have browser: MSIE on Android))
	// n.b. each user have a list of browsers
	fmt.Fprintln(out, "found users:\n"+foundUsers)

	// unique browsers count, total
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
