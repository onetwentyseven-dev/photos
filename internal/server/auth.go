package server

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"photos"
	"photos/internal"

	"github.com/fatih/structs"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

func (s *Server) sesssionFromRequest(r *http.Request) (*sessions.Session, error) {
	return s.sessions.Get(r, "photos-session")
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {

	session, err := s.sesssionFromRequest(r)
	if err != nil {
		// Create an error page and redirect to that. Use session flashing to flash an internal error message of sorts
		s.logger.WithError(err).Error("failed to load session")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	state, err := generateRandomState()
	if err != nil {
		// Create an error page and redirect to that. Use session flashing to flash an internal error message of sorts
		s.logger.WithError(err).Error("failed to generate state for authentication request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	session.Values["state"] = state

	err = session.Save(r, w)
	if err != nil {
		s.logger.WithError(err).Error("failed to save state for authentication request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", s.authenticator.AuthCodeURL(state))

	w.WriteHeader(http.StatusTemporaryRedirect)

}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {

	session, err := s.sesssionFromRequest(r)
	if err != nil {
		// Create an error page and redirect to that. Use session flashing to flash an internal error message of sorts
		s.logger.WithError(err).Error("failed to load session")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		s.logger.WithError(err).Error("failed to exchange code for token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%sv2/logout?client_id=%s&returnTo=%s", s.authenticator.IssuerURL, s.authenticator.ClientID, url.QueryEscape(s.appURL)))
	w.WriteHeader(http.StatusTemporaryRedirect)

}

func (s *Server) handleOauthCallback(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	query := r.URL.Query()

	state := query.Get("state")
	code := query.Get("code")

	if state == "" || code == "" {
		s.writeRedirectRouteName(w, "homepage")
		return
	}

	session, err := s.sesssionFromRequest(r)
	if err != nil {
		// Create an error page and redirect to that. Use session flashing to flash an internal error message of sorts
		s.logger.WithError(err).Error("failed to load session")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessionState, ok := session.Values["state"]
	if !ok {
		s.logger.Error("invalid session, no state stored in session")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if sessionState.(string) != state {
		s.logger.Error("session state does equal query state, discarding")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := s.authenticator.Exchange(ctx, code)
	if err != nil {
		s.logger.WithError(err).Error("failed to exchange code for token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	idToken, err := s.authenticator.VerifyIDToken(ctx, token)
	if err != nil {
		s.logger.WithError(err).Error("failed to verify id token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var profile = make(map[string]any)
	err = idToken.Claims(&profile)
	if err != nil {
		s.logger.WithError(err).Error("failed to provision claims token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	emailInf, ok := profile["name"]
	if !ok {
		s.logger.WithError(err).Error("profile is missing informaiton necessary to identify user")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := s.userRepository.UserByEmail(ctx, emailInf.(string))
	if err != nil {
		s.logger.WithError(err).Error("failed to look up user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user == nil {
		user = &photos.User{
			ID:    uuid.New().String(),
			Name:  fmt.Sprintf("%s %s", profile["given_name"], profile["family_name"]),
			Email: emailInf.(string),
		}

		// Support reaching out to the profile api to retrive emaployee id and profile uri

		err := s.userRepository.CreateUser(ctx, user)
		if err != nil {
			s.logger.WithError(err).Error("failed to save user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	session.Values["userID"] = user.ID

	err = session.Save(r, w)
	if err != nil {
		s.logger.WithError(err).Error("failed to exchange code for token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.writeRedirectRouteName(w, "dashboard")

}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}

// API Routes
func (s *Server) handleValidateAuth(w http.ResponseWriter, r *http.Request) {

	user := internal.UserFromContext(r.Context())
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(structs.Map(user))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
