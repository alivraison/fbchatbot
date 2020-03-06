package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fnproject/fn/api/server"
	"github.com/fnproject/fn/fnext"
)

const (
	fbAPI = "https://graph.facebook.com/v6/me/messages?access_token=%s"
	image = "http://37.media.tumblr.com/e705e901302b5925ffb2bcf3cacb5bcd/tumblr_n6vxziSQD11slv6upo3_500.gif"
)

// Callback from
type Callback struct {
	Object string `json:"object,omitempty"`
	Entry  []struct {
		ID        string      `json:"id,omitempty"`
		Time      int         `json:"time,omitempty"`
		Messaging []Messaging `json:"messaging,omitempty"`
	} `json:"entry,omitempty"`
}

// Messaging from
type Messaging struct {
	Sender    User    `json:"sender,omitempty"`
	Recipient User    `json:"recipient,omitempty"`
	Timestamp int     `json:"timestamp,omitempty"`
	Message   Message `json:"message,omitempty"`
}

// User from
type User struct {
	ID string `json:"id,omitempty"`
}

// Message from
type Message struct {
	MID        string `json:"mid,omitempty"`
	Text       string `json:"text,omitempty"`
	QuickReply *struct {
		Payload string `json:"payload,omitempty"`
	} `json:"quick_reply,omitempty"`
	Attachments *[]Attachment `json:"attachments,omitempty"`
	Attachment  *Attachment   `json:"attachment,omitempty"`
}

// Attachment from
type Attachment struct {
	Type    string  `json:"type,omitempty"`
	Payload Payload `json:"payload,omitempty"`
}

// Response from
type Response struct {
	Recipient User    `json:"recipient,omitempty"`
	Message   Message `json:"message,omitempty"`
}

// Payload from
type Payload struct {
	URL string `json:"url,omitempty"`
}

func init() {
	server.RegisterExtension(&fbmsgExt{})
}

type fbmsgExt struct {
}

func (e *fbmsgExt) Name() string {
	return "github.com/alivraison/fbchatbot/message"
}

func (e *fbmsgExt) Setup(s fnext.ExtServer) error {
	s.AddEndpoint("POST", "/webhook", &messageEndpoint{})
	return nil
}

type messageEndpoint struct {
}

func (v *messageEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var callback Callback
	json.NewDecoder(r.Body).Decode(&callback)
	if callback.Object == "page" {
		for _, entry := range callback.Entry {
			for _, event := range entry.Messaging {
				//ProcessMessage(event)
				fmt.
			}
		}
		w.WriteHeader(200)
		w.Write([]byte("Got your message"))
	} else {
		w.WriteHeader(404)
		w.Write([]byte("Message not supported"))
	}
}

// ProcessMessage for
func ProcessMessage(event Messaging) {
	client := &http.Client{}
	response := Response{
		Recipient: User{
			ID: event.Sender.ID,
		},
		Message: Message{
			Attachment: &Attachment{
				Type: "image",
				Payload: Payload{
					URL: image,
				},
			},
		},
	}
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(&response)
	url := fmt.Sprintf(fbAPI, os.Getenv("PAGE_ACCESS_TOKEN"))
	req, err := http.NewRequest("POST", url, body)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}
