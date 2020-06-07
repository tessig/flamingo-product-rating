flamingo: {
	"cmd.name":                   "rating"
	"systemendpoint.serviceAddr": ":13210"
}

core: {
	"locale.date.location": "Europe/Berlin"
	gotemplate: engine: {
		"templates.basepath": "templates"
		"layout.dir":         "layouts"
	}
}

mysql: {
	// env vars are always string values, so we have to write the bool value manually
    if flamingo.os.env.AUTOMIGRATE=="true" {
		migration: automigrate: true
    }

	db: {
		host:         flamingo.os.env.DBHOST
		port:         flamingo.os.env.DBPORT
		databaseName: flamingo.os.env.DBNAME
		user:         flamingo.os.env.DBUSER
		password:     flamingo.os.env.DBPASSWORD
	}
}

productservice: baseurl: "flamingo.os.env.PRODUCTSERVICEURL"
