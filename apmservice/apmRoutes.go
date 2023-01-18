package apmservice

import "net/http"

func Init() {
	http.HandleFunc("/test", test)
	http.HandleFunc("/compareResults", compareResults)
	http.HandleFunc("/htmlReport/", htmlReport)
}
