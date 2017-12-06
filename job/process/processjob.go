package process

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"ejabberd-restbroker/lib/flight"

	"github.com/jrallison/go-workers"
)

func Job(message *workers.Msg) {
	c := flight.Context()
	args, _ := message.Args().Array()

	endpoint := args[0].(string) // ejabberd endpoint or command
	params := args[1]            // ejabberd API Arguments
	token := args[2].(string)    // This should be the token
	returnid := args[3].(string) // id for response job

	// Check to see whether the authentication token is valid
	if token != c.Config.Token {
		log.Printf("Token invalid")
		return
	}

	// Params must be a json formatted string
	// Use passed endpoint and add the arguments
	req, err := http.NewRequest("POST",
		c.Config.RestURL+"/"+endpoint,
		strings.NewReader(params.(string)))

	if err != nil {
		log.Printf("Error Occured. %+v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", c.Config.JabberRestHost)
	req.Header.Set("X-Admin", "true")
	req.SetBasicAuth(c.Config.RestUser, c.Config.RestPass)

	// use httpClient to send request
	response, err := c.HttpClient.Do(req)
	if err != nil && response == nil {
		log.Printf("Error sending request to API endpoint. %+v", err)
	} else {
		// Close the connection to reuse it
		defer response.Body.Close()

		// Let's check if the work actually is done
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Printf("Couldn't parse response body. %+v", err)
			return
		}

		// If the job has a returnid, send the http response
		// to the ResponseQueue along with the id
		if returnid != "" {
			response := []string{string(body), returnid}
			workers.Enqueue(c.Config.ResponseQueue, "Add", response)
			log.Println("Sent Response")
		}

		log.Println("Response Body:", string(body))
	}
	return
}
