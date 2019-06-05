package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/coreos/go-oidc"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spark-cluster/pkg/apis"
	"github.com/spark-cluster/pkg/util/k8sutil"
	"golang.org/x/oauth2"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type APIHandler struct {
	frontDir string

	kubeConfig *rest.Config
	client     client.Client
	kubeClient kubernetes.Interface

	oauthConfig *oauth2.Config
	verifier    *oidc.IDTokenVerifier
	provider    *oidc.Provider
}

func NewAPIHandler(frontDir string) (*APIHandler, error) {
	kubeConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	// setup client set
	clientset, err := setupClient(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to setup kubernetes client: %v", err)
	}

	// setup kubernetes rest client
	kubeClient, err := k8sutil.NewKubeClient()
	if err != nil {
		return nil, fmt.Errorf("Failed to setup kubernetes client: %v", err)
	}

	apihandler := &APIHandler{
		frontDir: frontDir,

		kubeConfig: kubeConfig,
		client:     clientset,
		kubeClient: kubeClient,
	}
	// Setup authentication plugin
	authProvider := os.Getenv("OAUTH_PROVIDER")
	logrus.Infof("Setting up authentication provider: %v", authProvider)
	if authProvider == AuthProviderOAuth {
		issuerURL := os.Getenv("OAUTH_URL")
		apihandler.provider, err = oidc.NewProvider(context.Background(), issuerURL)
		if err != nil {
			return nil, errors.Wrap(err, "failed to connect oauth provider")
		}

		apihandler.oauthConfig = &oauth2.Config{
			ClientID:     os.Getenv("OAUTH_CLIENT"),
			Endpoint:     apihandler.provider.Endpoint(),
			ClientSecret: os.Getenv("OAUTH_SECRET"),
			RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
			Scopes:       []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess, oauthScopeProfile, oauthScopeEmail},
		}

		apihandler.verifier = apihandler.provider.Verifier(&oidc.Config{ClientID: apihandler.oauthConfig.ClientID})
	}

	return apihandler, nil
}

func setupClient(config *rest.Config) (client.Client, error) {
	scheme := runtime.NewScheme()
	for _, addToSchemeFunc := range []func(s *runtime.Scheme) error{
		apis.AddToScheme,
		v1.AddToScheme,
	} {
		if err := addToSchemeFunc(scheme); err != nil {
			return nil, err
		}
	}

	clientset, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

type Message struct {
	Message string `json:"message"`
}

func responseJSON(body interface{}, w http.ResponseWriter, statusCode int) {
	jsonResponse, err := json.Marshal(body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}
