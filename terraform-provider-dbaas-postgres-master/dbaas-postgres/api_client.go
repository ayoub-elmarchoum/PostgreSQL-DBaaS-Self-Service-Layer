package dbaas

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	//	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

type apiClientOpt struct {
	uri      string
	cert     string
	key      string
	ca       string
	token    string
	insecure bool
	username string
	password string
	headers  map[string]string
	timeout  int
	debug    bool
}

type api_client struct {
	http_client *http.Client
	uri         string
	cert        string
	key         string
	ca          string
	insecure    bool
	token       string
	username    string
	password    string
	headers     map[string]string
	timeout     int
	debug       bool
}

// Constructeur
func NewAPIClient(opt *apiClientOpt) (*api_client, error) {
	if opt.uri == "" {
		return nil, errors.New("No URI defined, please set the dbaas uri.")
	}

	opt.uri = strings.TrimRight(opt.uri, "/")

	tlsConfig := &tls.Config{
		InsecureSkipVerify: opt.insecure,
	}

	if opt.cert != "" && opt.key != "" {
		var cert tls.Certificate
		var err error
		if strings.HasPrefix(opt.cert, "-----BEGIN") && strings.HasPrefix(opt.key, "-----BEGIN") {
			cert, err = tls.X509KeyPair([]byte(opt.cert), []byte(opt.key))
		} else {
			cert, err = tls.LoadX509KeyPair(opt.cert, opt.key)
		}
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	if opt.ca != "" {
		var caCert []byte
		var err error
		if strings.HasPrefix(opt.ca, "-----BEGIN") {
			caCert = []byte(opt.ca)
		} else {
			caCert, err = ioutil.ReadFile(opt.ca)
			if err != nil {
				return nil, err
			}
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &api_client{
		http_client: &http.Client{
			Timeout:   time.Second * time.Duration(opt.timeout),
			Transport: tr,
		},
		uri:      opt.uri,
		cert:     opt.cert,
		key:      opt.key,
		ca:       opt.ca,
		insecure: opt.insecure,
		token:    opt.token,
		username: opt.username,
		password: opt.password,
		headers:  make(map[string]string),
		timeout:  opt.timeout,
		debug:    opt.debug,
	}

	// Si tu veux, tu peux initialiser des headers par défaut ici
	for k, v := range opt.headers {
		client.headers[k] = v
	}

	return client, nil
}

// Envoyer la requête HTTP avec headers personnalisés
func (client *api_client) send_request(method string, path string, data string, headers map[string]string) (string, error) {
	full_uri := client.uri + path
	var req *http.Request
	var err error

	buffer := bytes.NewBuffer([]byte(data))
	if data == "" {
		req, err = http.NewRequest(method, full_uri, nil)
	} else {
		req, err = http.NewRequest(method, full_uri, buffer)
	}
	if err != nil {
		return "", err
	}

	// Toujours ajouter le header Authorization. Si token est vide, ça ne fait rien.
	if client.token != "" {
		req.Header.Set("Authorization", "Bearer "+client.token)
	}

	// Ajouter headers du client (headers globaux déclarés ou initialisés)
	for n, v := range client.headers {
		req.Header.Set(n, v)
	}

	// Ajouter headers passés explicitement à la fonction
	if len(headers) > 0 {
		for n, v := range headers {
			req.Header.Set(n, v)
		}
	}

	// Basic auth si configuré
	if client.username != "" && client.password != "" {
		req.SetBasicAuth(client.username, client.password)
	}

	// Debug
	if client.debug {
		reqDump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			return "", err
		}
		fmt.Printf("REQUÊTE:\n%s\n", string(reqDump))
	}

	resp, err := client.http_client.Do(req)
	if err != nil {
		return "", err
	}

	if client.debug {
		respDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Printf("Erreur dump réponse: %s\n", err)
		} else {
			fmt.Printf("RÉPONSE:\n%s\n", string(respDump))
		}
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}

	body := string(bodyBytes)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return body, errors.New(fmt.Sprintf("Code HTTP inattendu : %d - %s", resp.StatusCode, body))
	}

	return body, nil
}
