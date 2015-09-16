package server

import "net/http"

func authorizeSigning(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("X-API-KEY") != config["sign-token"] {
		res.WriteHeader(http.StatusUnauthorized)
	}
}
