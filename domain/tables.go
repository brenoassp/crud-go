package domain

import "github.com/vingarcia/ksql"

var ClientsTable = ksql.NewTable("clients", "id")
