package db

type determineFactoryCase struct {
	databaseType string
	validate     validateConfigFunc
	init         initDbFunc
	isError      bool
}

var determineFactoryCases = []determineFactoryCase{
	{"sqlite3", sqliteValidateConfig, sqliteInitDb, false},
	{"postgres", postgresValidateConfig, postgresInitDb, false},
	{"non-existant", nil, nil, true},
	{"", nil, nil, true},
}
