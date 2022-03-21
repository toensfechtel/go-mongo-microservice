package utils

import (
    "crypto/rand"
    "fmt"
)

/**
*Generates uuid => for ProduceCode value
**/
func UUID() (*string, error) {

    b := make([]byte, 8)
    _, err := rand.Read(b)
    if err != nil {
        return nil, err
    }
    idgened := fmt.Sprintf("%X-%X-%X-%X", b[0:2], b[2:4], b[4:6], b[6:8])

    return &idgened, nil
}

