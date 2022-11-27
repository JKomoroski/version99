package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"hash"
	"log"
	"net/http"
	"regexp"
)

// --------------------------------------------------------------------

func logRequest(status int, r *http.Request) {
	log.Printf("%d %s", status, r.RequestURI)
}

// --------------------------------------------------------------------

func matchMavenURI(uri string) *maven {
	re := regexp.MustCompile("/mvn2/(.+)/([^/]+)/99\\.0-does-not-exist/(.+)-99\\.0-does-not-exist\\.(jar|pom)(\\.(sha1|md5))?")
	m := re.FindStringSubmatch(uri)

	// no match -> return
	if m == nil {
		return nil
	}

	// artifactId != package name -> return
	if m[2] != m[3] {
		return nil
	}

	return &maven{GroupId: m[1], ArtifactId: m[2], Name: m[3], Ext: m[4], Digest: m[6], InfoURL: INFO_URL}
}

// --------------------------------------------------------------------

func makeDigest(data []byte, digest string) (string, error) {
	var h hash.Hash

	switch digest {
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	}

	if h == nil {
		return "", fmt.Errorf("Invalid hash extension")
	}

	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// --------------------------------------------------------------------

func mustBase64Decode(s string) []byte {
	buf, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return buf
}

// --------------------------------------------------------------------

func send(data []byte, contentType string, status int, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", contentType)
	w.Write(data)
	logRequest(status, r)
}

// --------------------------------------------------------------------

func sendOK(data []byte, contentType string, w http.ResponseWriter, r *http.Request) {
	send(data, contentType, http.StatusOK, w, r)
}

// --------------------------------------------------------------------

func sendDigest(data []byte, digest string, w http.ResponseWriter, r *http.Request) {
	d, err := makeDigest(data, digest)
	if err != nil {
		sendError(err, w, r)
		return
	}

	sendOK([]byte(d), "text/plain", w, r)
}

// --------------------------------------------------------------------

func sendError(err error, w http.ResponseWriter, r *http.Request) {
	log.Print(err)
	sendStatus(http.StatusInternalServerError, w, r)
}

// --------------------------------------------------------------------

func sendNotFound(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	err := notFoundTemplate.Execute(&buf, &notFound{URL: r.RequestURI, InfoURL: INFO_URL})
	if err != nil {
		sendError(err, w, r)
		return
	}

	send(buf.Bytes(), "text/html", http.StatusNotFound, w, r)
}

// --------------------------------------------------------------------

func sendStatus(status int, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(status)
	logRequest(status, r)
}
