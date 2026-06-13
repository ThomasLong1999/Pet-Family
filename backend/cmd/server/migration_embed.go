package main

import _ "embed"

//go:embed migration.sql
var migrationSQL string
