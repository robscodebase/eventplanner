package main

import (
 "net/http"
 )



 func main() {
    logger(nil, "main", "main.go", "main", "Listening on port 8080")
        http.ListenAndServe(":8080", router())
        }
