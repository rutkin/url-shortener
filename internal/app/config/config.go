package config

import "net/url"

var Address = url.URL{Scheme: "http", Host: "localhost:8080"}
