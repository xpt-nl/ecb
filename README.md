# ecb

    import "github.com/xpt-nl/ecb"

[![Go Reference](https://pkg.go.dev/badge/github.com/xpt-nl/ecb.svg)](https://pkg.go.dev/github.com/xpt-nl/ecb#section-documentation)

Package `ecb` provides a function to get the current EURO exchange rate directly from the ECB denominated in another currency.

## Euro foreign exchange reference rates

The reference rates are usually updated around 16:00 CET on every working day,
except on TARGET closing days. They are based on a regular daily concertation
procedure between central banks across Europe, which normally takes place
at 14:15 CET.

## Example

```go
package main

import "github.com/xpt-nl/ecb"

func main() {
	rate, err := ecb.EUR(ecb.USD)
	if err != nil {
		panic(err)
	}
	println(rate)
}
```
