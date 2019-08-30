package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/deadmanssnitch/snshttp"
)

type EventHandler struct {
	// DefaultHandler provides auto confirmation of subscriptions and ignores
	// unsubscribe events.
	snshttp.DefaultHandler
}

// Notification is called for messages published to the SNS Topic. When using
// DefaultHandler as above this is the only event you need to implement.
/*
  Sample output:
  id="5c1b2a20-748f-5f6e-88c5-7fc25a30a8e4"
  subject="UPDATE: AWS CodeCommit us-west-2 push: mandrill-puppet"
  message= {
    "Records":
    [
      {
        "awsRegion":"us-west-2",
        "codecommit":{
          "references":
          [
            {
              "commit":"8212f68b310a885ce430dc5b702e280c02cab2e1",
              "ref":"refs/heads/production"
            }
          ]
        },
        "eventId":"099cfebf-5503-4ed6-8200-417ff2cf1637",
        "eventName":"ReferenceChanges",
        "eventPartNumber":1,
        "eventSource":"aws:codecommit",
        "eventSourceARN":"arn:aws:codecommit:us-west-2:175749097276:mandrill-puppet",
        "eventTime":"2019-08-29T19:50:23.648+0000",
        "eventTotalParts":1,
        "eventTriggerConfigId":"1a3b824a-68ad-4b02-adbe-9b5757cb2043",
        "eventTriggerName":"r10k",
        "eventVersion":"1.0",
        "userIdentityARN":"arn:aws:iam::175749097276:user/puppet"
      }
    ]
  }
  timestamp="2019-08-29 19:50:23.678 +0000 UTC"
*/
func (h *EventHandler) Notification(ctx context.Context, event *snshttp.Notification) error {
	// fmt.Printf("id=%q subject=%q message=%q timestamp=%q\n",
	//	event.MessageID,
	//	event.Subject,
	//	event.Message,
	//	event.Timestamp,
	// )
	r10k := "/opt/puppetlabs/puppet/bin/r10k"
	output, err := exec.Command(r10k, "deploy", "environment").CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	fmt.Println(string(output))
	fmt.Println("Received a Notification and executed r10k with exit code 0")

	return nil
}

func main() {
	// snshttp.New returns an http.Handler that will parse the payload and pass the
	// event to the provided EventHandler.
	snsHandler := snshttp.New(&EventHandler{}, snshttp.WithAuthentication("", ""))
	http.Handle("/hooks/sns", snsHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.ListenAndServe(":8080", nil)
}
