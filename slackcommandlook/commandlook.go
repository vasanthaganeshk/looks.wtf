package slackcommandlook

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	_ "google.golang.org/appengine"
)

var slackVerificationToken = os.Getenv("SLACK_VERIFICATION_TOKEN")

var tags = loadTags("tags.yml")
var looksByTags = loadLooks("looks.yml")

type slackCommandResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

func commandLookHandler(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	token := r.FormValue("token")
	if token != slackVerificationToken {
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return nil
	}

	teamDomain := r.FormValue("team_domain")
	channelName := r.FormValue("channel_name")
	userName := r.FormValue("user_name")
	command := r.FormValue("command")
	tag := r.FormValue("text")

	log.Printf("Request: TeamDomain: %s ChannelName: %s UserName: %s Command: %s Text: %s", teamDomain, channelName, userName, command, tag)

	looks := looksByTags[tag]
	if len(looks) == 0 {
		fmt.Fprintf(w, "Try using the /look command with one of these words: "+strings.Join(tags, ", "))
		return nil
	}

	l := looks[rand.Intn(len(looks))]

	w.Header().Add("Content-Type", "application/json")
	response := slackCommandResponse{
		ResponseType: "in_channel",
		Text:         l.Plain,
	}
	enc := json.NewEncoder(w)
	return enc.Encode(response)
}
