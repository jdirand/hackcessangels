package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func (s *Server) handleAgentPosition(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Latitude  *float64 `json:"latitude,omitempty"`
		Longitude *float64 `json:"longitude,omitempty"`
		Precision *float64 `json:"precision,omitempty"`
	}

	w.Header().Add("Content-Type", "application/json")
	err := getJSONRequest(r, &data)
	if err != nil {
		log.Print(err)
		returnError(400, "Invalid request", w)
		return
	}

	session, err := s.store.Get(r, "user")
	if err != nil {
		returnError(401, "Please log in", w)
		return
	}

	loggedEmail, ok := session.Values["email"].(string)
	if !ok {
		returnError(401, "Please log in", w)
		return
	}

	user, err := s.model.GetUserByEmail(loggedEmail)
	if err != nil {
		returnError(400, "Error while getting user data", w)
		return
	}

	if !user.IsAgent {
		returnError(403, "You are not an agent", w)
		return
	}

	switch r.Method {
	case "POST":
		if data.Latitude == nil || data.Longitude == nil {
			returnError(400, "latitude and longitude must be present", w)
			return
		}
		var precision float64
		if data.Precision == nil {
			precision = 0.0
		} else {
			precision = *data.Precision
		}
		station, err := s.model.FindStationByLocation(
			*data.Longitude, *data.Latitude, precision)
		if err != nil {
			log.Print("FindStationByLocation error: ", err)
			returnError(500, "Error while processing request", w)
			return
		}
		if station == nil {
			user.CurrentStation = nil
		} else {
			user.CurrentStation = &station.Id
		}
		json.NewEncoder(w).Encode(station)
		return
	}
}

func (s *Server) handleAgentRequests(w http.ResponseWriter, r *http.Request) {
	var data struct {
		RequestId     string
		TakeRequest   bool
		FinishRequest bool
	}

	w.Header().Add("Content-Type", "application/json")
	err := getJSONRequest(r, &data)
	if err != nil {
		log.Print(err)
		returnError(400, "Invalid request", w)
		return
	}

	session, err := s.store.Get(r, "user")
	if err != nil {
		returnError(401, "Please log in", w)
		return
	}

	loggedEmail, ok := session.Values["email"].(string)
	if !ok {
		returnError(401, "Please log in", w)
		return
	}

	user, err := s.model.GetUserByEmail(loggedEmail)
	if err != nil {
		returnError(400, "Error while getting user data", w)
		return
	}

	if !user.IsAgent {
		returnError(403, "You are not an agent", w)
		return
	}

	switch r.Method {
	case "GET":
		// Returns all active requests that may be answered by this agent, or a precise request if RequestId is provided.
		station, err := user.GetStation()
		if err != nil {
			log.Print(err)
			returnError(500, "Unable to get station", w)
			return
		}
		if station == nil {
			returnError(404, "Agent not in station", w)
			return
		}
		requests, err := s.model.GetActiveRequestsByStation(station)
		if err != nil {
			log.Print(err)
			returnError(500, "Unable to get active requests", w)
			return
		}

		apiRequests := make([]*APIRequest, 0)
		for _, helpRequest := range requests {
			apiRequests = append(apiRequests, NewAPIRequestFromHelpRequest(helpRequest))
		}
		json.NewEncoder(w).Encode(apiRequests)
		return
	case "POST":
		// Sets the response to the help request
		return
	}

}
