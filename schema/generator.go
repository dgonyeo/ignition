package schema

//go:generate schematyper --package=types ignition.json -o ../config/types/schema.go --root-type=Config

/*
 * This file exists solely to provide the go generate directive.
 *
 * The above directive parses the json schema contained in this directory and
 * writes out go structs to schema.go that represent this schema. This happens
 * during the build script, where go generate is called.
 *
 * schematyper must be present on the system building ignition.
 */
