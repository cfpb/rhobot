## simple makefile to log workflow
# .PHONY: all test clean build install

GOFLAGS ?= -v -a $(GOFLAGS:)

PGHOST ?=localhost
PGPORT ?=5432
PGDATABASE ?=my_database
PGUSER ?=my_user
PGPASSWORD ?=my_password

all: build test

.PHONY: build
build:
	@mkdir -p build
	@go build $(GOFLAGS) -o build github.com/cfpb/rhobot/cmd/rhobot

test:
	@go test $(GOFLAGS) ./...

ccdb.csv:
	curl -o ccdb.csv https://data.consumerfinance.gov/api/views/s6ew-h6mp/rows.csv?accessType=DOWNLOAD

fixtures: ccdb.csv
	psql -c ' BEGIN; CREATE SCHEMA ccdb; CREATE TABLE ccdb.record ( "Date received" text, "Product" text, "Sub-product" text, "Issue" text, "Sub-issue" text, "Consumer complaint narrative" text, "Company public response" text, "Company" text, "State" text, "ZIP code" text, "Tags" text, "Consumer consent provideded" text, "Submitted via" text, "Date sent to company" text, "Company response to consumer" text, "Timely response?" text, "Consumer disputed?" text, "Complaint ID" text); COMMIT;'
	psql -c "\copy ccdb.record FROM 'ccdb.csv' DELIMITER ',' CSV HEADER;"

clean:
	psql -c ' DROP SCHEMA ccdb CASCADE; '
