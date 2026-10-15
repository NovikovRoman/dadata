# Client for DaData.ru

Forked from <https://github.com/ekomobile/dadata>.

DaData API v2

## Installation

`go get github.com/NovikovRoman/dadata/v2`

## Usage

```go
import (
 "context"
 "fmt"

 "github.com/NovikovRoman/dadata/v2"
 "github.com/NovikovRoman/dadata/v2/api/suggest"
)

func DaDataExample()  {
 api := dadata.NewSuggestApi()

 params := suggest.RequestParams{
  Query: "ул Свободы",
 }

 suggestions, err := api.Address(context.Background(), &params)
 if err != nil {
  return
 }

 for _, s := range suggestions {
  fmt.Printf("%s", s.Value)
 }
}
```

## Configuration

### Credentials

`DADATA_API_KEY` and `DADATA_SECRET_KEY` environment variables are used by default to authenticate client.

Custom credential provider may be used by implementing `client.CredentialProvider` interface.

Also, there is a "simple" credential provider `client.Credentials` you may utilize.

```go
creds := client.Credentials{
    ApiKeyValue:    "<YOUR_API_KEY>",
    SecretKeyValue: "<YOUR_SECRET_KEY>",
}

api := NewSuggestApi(client.WithCredentialProvider(&creds))
```

### HTTP client

HTTP client may be overridden with custom one:

```go
httpClient := &http.Client{}

api := NewSuggestApi(WithHttpClient(httpClient))
```

### Low level client usage

Pass and consume raw bytes/strings:

```go
    endpointUrl, err := url.Parse("https://suggestions.dadata.ru/suggestions/api/4_1/rs/")
    if err != nil {
        return
    }

    cli := client.NewClient(endpointUrl,
        client.WithEncoderFactory(encoder.RawEncoderFactory()),
        client.WithDecoderFactory(encoder.RawDecoderFactory()),     
    )

    request := bytes.NewBufferString("{ \"query\": \"москва хабар\" }")
    response := &bytes.Buffer{}

    err = cli.Post(context.Background(), "suggest/address", request, response)
    if err != nil {
        return
    }

    fmt.Print(response.String()) // Output: `{"suggestions":[{"value":"г Москва, ул Хабаровская",...`
```

## Licence

MIT see [LICENSE](LICENSE)
