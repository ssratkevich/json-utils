# json-utils

Go json-utils package for removing comments and trailing commas from JSON files. It also checks for `BOM` and removes it if necessary.


Fast usage sample:

``` go
import (
    "encoding/json"
    "os"
    json_utils "github.com/ssratkevich/json-utils"
)

func getData(name string) (any, error) {
    src, err := os.ReadFile(name)
    if err != nil {
        return nil, err
    }
    src = json_utils.FixJson(src)
    // parsing and handling JSON
    var data any
    err = json.Unmarshal(src, &data)
    return data, err
}
```
