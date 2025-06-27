package initializers

import (
	"encoding/json"
	"net/http"

	"example/hostverifier"
	"example/serverclassfactory"
)

type FeatureCountHandler struct {
	hostVerifier      hostverifier.HostVerifier
	serverClassFactory serverclassfactory.ServerClassFactory
}

func NewFeatureCountHandler(hostVerifier hostverifier.HostVerifier, serverClassFactory serverclassfactory.ServerClassFactory) *FeatureCountHandler {
	return &FeatureCountHandler{
		hostVerifier:      hostVerifier,
		serverClassFactory: serverClassFactory,
	}
}

func (f *FeatureCountHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if !f.hostVerifier.VerifyHost(r) {
		http.Error(w, "Host verification failed", http.StatusUnauthorized)
		return
	}

	serverClass, err := f.serverClassFactory.InstantiateServerClass()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	featureCount, err := serverClass.GetFeatureCount()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(featureCount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}