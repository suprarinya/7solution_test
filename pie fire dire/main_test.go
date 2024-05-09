package main

import (
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestMeatSummaryHandler(t *testing.T) {
    // create server for mock response
	testData := "Fatback t-bone t-bone, pastrami t-bone. pork, meatloaf jowl enim. Bresaola t-bone."
    testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// return test data as API response 
        w.Write([]byte(testData)) 
    }))
    defer testServer.Close()

    // change to url of mok server
    originalURL := url
	url = testServer.URL
	// change back to original url after the test
    defer func() { url = originalURL }() 

    // create new request
    recorder := httptest.NewRecorder()
    request, _ := http.NewRequest("GET", "/beef/summary", nil)

	meatSummaryHandler(recorder, request)

	// check status code
    assert.Equal(t, http.StatusOK, recorder.Code, "handler returned wrong status code")

	// check body
    expected := `{"beef":{"bresaola":1,"enim":1,"fatback":1,"jowl":1,"meatloaf":1,"pastrami":1,"pork":1,"t-bone":4}}`
    response, _ := ioutil.ReadAll(recorder.Body)
    assert.JSONEq(t, expected, string(response), "handler returned unexpected body")
}

