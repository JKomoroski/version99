package main

import (
	"bytes"
	"net/http"
	"text/template"
)

const (
	INFO_URL    = "http://day-to-day-stuff.blogspot.com/2007/10/announcement-version-99-does-not-exist.html"
	FAVICON_PNG = "iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAIAAACQkWg2AAAAOElEQVR4AWP4z0AaolwDGCCzKdeAiZAA0X4AA5p4mjCiix3/MQEOQTDCEfSYguTbQDgCKLeBJAQA5cd6lDERP0gAAAAASUVORK5CYII="
	EMPTY_JAR   = "UEsDBAoAAAAAAME+SDiyfwLuGQAAABkAAAAUAAQATUVUQS1JTkYvTUFOSUZFU1QuTUb+ygAATWFuaWZlc3QtVmVyc2lvbjogMS4wDQoNClBLAQIKAAoAAAAAAME+SDiyfwLuGQAAABkAAAAUAAQAAAAAAAAAAAAAAAAAAABNRVRBLUlORi9NQU5JRkVTVC5NRv7KAABQSwUGAAAAAAEAAQBGAAAATwAAAAAA"

	POM_XML = `<?xml version="1.0" encoding="iso-8859-1"?>
<project>
	<modelVersion>4.0.0</modelVersion>
	<groupId>{{.GroupId}}</groupId>
	<artifactId>{{.ArtifactId}}</artifactId>
	<name>{{.Name}}</name>
	<version>99.0-does-not-exist</version>
	<description>
		This is a generated pom.
		Version 99.0-does-not-exist is a dummy implementation which actually does nothing and has no dependencies. 
		VERSION 99.0-does-not-exist IS NOT IN ANY WAY AFFILIATED WITH THE ORIGINAL DEVELOPERS of {{.GroupId}}.
	</description>
	<url>{{.InfoURL}}</url>
</project>`

	INDEX_HTML = `<html>
<head>
<title>Version 99 Does Not Exist</title>
<style>
html { font-family: -apple-system, BlinkMacSystemFont, avenir next, avenir, segoe ui, helvetica neue, helvetica, Cantarell, Ubuntu, roboto, noto, arial, sans-serif; }
body { padding: 5px; max-width: 550px; }
h1 { font-size: 120%; }
img { width:1em; height:1em; position:relative; top:2px; }
</style>
<link rel="icon" href="/favicon.png" type="image/png">
</head>
<body>
<h1><img src="/favicon.png"> Version 99 Does Not Exist</h1>
<p>Please see <a href="http://day-to-day-stuff.blogspot.com/2007/10/announcement-version-99-does-not-exist.html">my blog</a> to read why I created Version 99 Does Not Exist and its predecessor no-commons-logging.</p>
<p>Version 99 Does Not Exist emulates a Maven 2 repository and serves empty jars for any valid package that has version number <i>99.0-does-not-exist</i>. It also generates poms, <span style="text-decoration: line-through">metadata files</span> (removed since 2.0) and of course the appropriate hashes.</p>
<p>For example the following links will give an <a href="/mvn2/commons-logging/commons-logging/99.0-does-not-exist/commons-logging-99.0-does-not-exist.jar">empty jar</a>, its <a href="/mvn2/commons-logging/commons-logging/99.0-does-not-exist/commons-logging-99.0-does-not-exist.pom">pom</a> and the <a href="/mvn2/commons-logging/commons-logging/maven-metadata.xml"><span style="text-decoration: line-through">maven metadata</span></a> for commons-logging.</p>
<p><a href="https://github.com/erikvanoosten/version99">Vesion 99 Does Not Exist source code on GitHub.</a></p>
</body>
</html>`

	NOT_FOUND_HTML = `<html>
<head>
<title>Version 99 Does Not Exist</title>
<style>
html { font-family: -apple-system, BlinkMacSystemFont, avenir next, avenir, segoe ui, helvetica neue, helvetica, Cantarell, Ubuntu, roboto, noto, arial, sans-serif; }
body { padding: 5px; max-width: 550px; }
h1 { font-size: 120%;  }
h1 > a { text-decoration: inherit; color: inherit; }
img { width:1em; height:1em; position:relative; top:2px; }
</style>
<link rel="icon" href="/favicon.png" type="image/png">
</head>
<body>
<h1><a href="/"><img src="/favicon.png"> Version 99 Does Not Exist - Error 404</a></h1>
<h2>404 Not Found: {{.URL}}</h2>
<p>
<a href="{{.InfoURL}}">Version 99 Does Not Exist</a> is a virtual Maven2 repository. 
It generates jars and poms for any artifact with version <tt>99.0-does-not-exist</tt>.
</p>
</body>
</html>`
)

// --------------------------------------------------------------------

type notFound struct {
	URL     string
	InfoURL string
}

// --------------------------------------------------------------------

type maven struct {
	GroupId    string
	ArtifactId string
	Name       string
	Ext        string
	Digest     string
	InfoURL    string
}

// --------------------------------------------------------------------

var (
	notFoundTemplate = template.Must(template.New("notFound").Parse(NOT_FOUND_HTML))
	pomTemplate      = template.Must(template.New("pom").Parse(POM_XML))
	emptyJar         = mustBase64Decode(EMPTY_JAR)
	favicon          = mustBase64Decode(FAVICON_PNG)
)

// --------------------------------------------------------------------

func handler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" && r.Method != "HEAD" {
		sendStatus(http.StatusMethodNotAllowed, w, r)
		return
	}

	// Handle HEAD request
	if r.Method == "HEAD" {
		sendStatus(http.StatusOK, w, r)
		return
	}

	if r.RequestURI == "/" || r.RequestURI == "/index.html" {
		sendOK([]byte(INDEX_HTML), "text/html", w, r)
		return
	}

	if r.RequestURI == "/favicon.png" {
		sendOK(favicon, "image/png", w, r)
		return
	}

	m := matchMavenURI(r.RequestURI)
	if m == nil {
		sendNotFound(w, r)
		return
	}

	switch m.Ext {
	case "jar":
		switch m.Digest {
		case "":
			sendOK(emptyJar, "application/jar", w, r)
			return
		case "sha1", "md5":
			sendDigest(emptyJar, m.Digest, w, r)
			return
		}
	case "pom":
		var buf bytes.Buffer
		err := pomTemplate.Execute(&buf, m)
		if err != nil {
			sendError(err, w, r)
			return
		}

		switch m.Digest {
		case "":
			sendOK(buf.Bytes(), "text/xml", w, r)
			return
		case "sha1", "md5":
			sendDigest(buf.Bytes(), m.Digest, w, r)
			return
		}
	}

	sendNotFound(w, r)
}
