package proc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"keycloakslackbot/api"
	"keycloakslackbot/logs"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	kc       *api.KeyCloak
	slackurl string
	interval int
	last     int64
}

func NewServer(
	slackurl string,
	keycloakurl string,
	user string,
	password string,
	interval int,
	realm string,
) Server {
	kc := api.NewKeyCloak(keycloakurl, realm, user,
		password)
	return Server{
		&kc,
		slackurl,
		interval,
		time.Now().UnixNano(),
	}
}

func (s *Server) Run() {
	logs.Logger.Info("Starting main loop")
	// Kick off subroutine
	go func() {
		for {
			newUsers := s.checkKC()
			s.sendToSlack(newUsers)
			time.Sleep(time.Duration(s.interval) * time.Second)
		}
	}()
	// Start server for K8s health
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("OK")) })
	http.ListenAndServe(":5000", nil)
}

// Checks for new users with creation after the last timestamp
func (s *Server) checkKC() []string {
	rs, err := s.kc.ListUsers()
	if err != nil {
		logs.Logger.Error("Cannot parse users because: " + err.Error())
	}

	var newUsers []string
	for _, r := range rs {
		now := s.last / 1e6
		if r.CreatedTimeStamp > now { // KC uses milliseconds
			newUsers = append(newUsers, r.Email)
		}
	}

	s.last = time.Now().UnixNano()

	return newUsers
}

// Decides whether to post to slack
func (s *Server) sendToSlack(newUsers []string) {
	if len(newUsers) == 0 {
		logs.Logger.Info("No new users to post")
		return
	}
	logs.Logger.Info(fmt.Sprintf("Sending %d new users to slack", len(newUsers)))
	values := map[string]string{
		"text": "New users: " + strings.Join(newUsers, ", "),
	}
	jsonValue, _ := json.Marshal(values)
	_, err := http.Post(s.slackurl, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		logs.Logger.Error("Error sending to slack: " + err.Error())
	}
}
