package process

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

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
	if token == c.Config.Token {

		// First need to safely convert arguments to json
		argsJson, err := GetJson(params)
		if err != nil {
			log.Printf("Unable to decode params. %+v", err)
			return
		}

		// Use passed endpoint and add the arguments
		req, err := http.NewRequest("POST",
			c.Config.RestURL+"/"+endpoint,
			bytes.NewBuffer(argsJson))
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
			return
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
	}
}

// This handles some inconsistencies in how arguments are passed
func GetJson(args interface{}) ([]byte, error) {
	switch v := args.(type) {
	case map[string]interface{}:
		argsJson, err := json.Marshal(args.(map[string]interface{}))
		if err != nil {
			log.Printf("Error Occured. %+v", err)
			return nil , err
		}
		return argsJson, nil
	case []interface{}:
		argsJson, err := json.Marshal((args.(interface{})))
		if err != nil {
			log.Printf("Error Occured. %+v", err)
			return nil , err
		}
		return argsJson, nil
	default:
		log.Println("Unknown argument format: ", v.(string))
		return []byte(""), nil
	}
	return []byte(""), nil
}
