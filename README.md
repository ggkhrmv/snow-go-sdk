# snow-go-sdk

A small, focused Go SDK for ServiceNow REST APIs.

Right now this repo provides:

- a configurable base client (`snow`)
- a generic Table API client (`snow/table`)
- an encoded query builder for `sysparm_query`

## Installation

```bash
go get github.com/ggkhrmv/snow-go-sdk
```

## Quick start

```go
package main

import (
	"context"
	"log"

	"github.com/ggkhrmv/snow-go-sdk/snow"
	"github.com/ggkhrmv/snow-go-sdk/snow/table"
)

func main() {
	client, err := snow.NewClient(
		snow.WithInstanceURL("https://dev12345.service-now.com"),
		snow.WithBasicAuth("admin", "password"),
	)
	if err != nil {
		log.Fatal(err)
	}

	incidents := table.NewMap(client, "incident")
	resp, err := incidents.List(context.Background(), &table.ListOptions{
		Limit: table.Int(10),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("records: %d", len(resp.Result))
}
```

## Typed records (generics)

You can use your own struct type instead of `map[string]any`.

```go
type Incident struct {
	SysID            string `json:"sys_id"`
	Number           string `json:"number"`
	ShortDescription string `json:"short_description"`
}

incidents := table.New[Incident](client, "incident")
resp, err := incidents.Get(ctx, "46d44f2c0a0a0b5700c77f9bf387afe3", nil)
item := resp.Result
```

## Encoded query builder

`ListOptions.Query` accepts a raw encoded query string.  
If you don't want to hand-write that string, use `table.QueryBuilder`:

```go
query, err := table.NewQueryBuilder().
	Eq("active", true).
	Or().
	Eq("priority", 1).
	Build()
if err != nil {
	log.Fatal(err)
}

resp, err := incidents.List(ctx, &table.ListOptions{
	Query: query,
})
```

It supports common operators (`Eq`, `NotEq`, `GT`, `GTE`, `LT`, `LTE`, `In`, `NotIn`, `Contains`, `StartsWith`, `EndsWith`, `IsEmpty`, `IsNotEmpty`) and logical chaining (`And`, `Or`, `NewQuery`).

## Error handling

Non-2xx responses return `*snow.APIError` when possible, including ServiceNow `error.message` and `error.detail` fields when present.

## Development

Run the full test suite:

```bash
go test ./...
```

Run one test:

```bash
go test ./snow/table -run TestQueryBuilderBuild
```

Build:

```bash
go build ./...
```
