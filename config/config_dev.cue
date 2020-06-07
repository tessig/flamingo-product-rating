core: zap: {
	json:     false
	loglevel: "Debug"
	colored:  true
}

flamingo: debug: mode: true

mysql: db: {
	host:         ""
	port:         "33306"
	databaseName: "ratings"
	user:         "ratings"
	password:     "ratings"
}

productservice: {
	baseurl:            "http://localhost:8080/"
	"endpoints.list":   "products"
	"endpoints.detail": "products/id/:pid"
}
