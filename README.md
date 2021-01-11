# qe
Go query encoding library inspired by encoding/json's Marshal/Unmarshal functions

### Installation

```bash
go get github.com/insprac/qe
```

### Basic Usage

An example of a generic search query.

```go
import "github.com/insprac/qe"

type SearchParams struct {
  Term      string `q:"term" required:"true"`
  Order     string `q:"order"`
  PageLimit uint8  `q:"page_limit"`
  Page      uint8  `q:"page"`
}

params = SearchParams{"some search term", "name", 10, 1}

queryString, err = qe.Marshal(params)
// term=some+search+term&order=name&page_limit=10&page=1
```

Another example using a list of IDs.

```go
import "github.com/insprac/qe"

type ItemQueryParams struct {
  IDs   []uint32 `q:"ids" required:"true"`
  Order string   `q:"order"`
}

params = ItemQueryParams{[]uint32{3, 20, 43}, "id"}

queryString, err = qe.Marshal(params)
// ids=3%2C20%2C43&order=id
```
