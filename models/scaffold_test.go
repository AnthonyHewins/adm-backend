package models

var DBConfigDefaultTest = &DB{
	Host:     "localhost",
	Port:     5432,
	Name:     "admtest",
	User:     "test",
	Password: "test",
}

func DBSetupTest(differentConfig *DB) {
	if differentConfig == nil {
		differentConfig = DBConfigDefaultTest
	}

	DBSetup(differentConfig)
}
