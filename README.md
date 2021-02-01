
# Run the application:

```
go run main.go server.go model.go
```

# Send a text:

```
curl --header "Content-Type: application/json" \
     --request POST \
     --data '{"description":"mój fragment kodu","content":"fmt.Println(\"Hello, world!\")"}' \
     http://localhost:8080/
```

# Sample response from the server:

```
{"id":"019fdf09-cf59-4091-9a61-3fbae0e7276c"}
```

# Read previously entered text:

```
curl --header "Content-Type: application/json" \
     --request GET \
     http://localhost:8080/paste/019fdf09-cf59-4091-9a61-3fbae0e7276c
```

# Sample response from the server:

```
{"id":"019fdf09-cf59-4091-9a61-3fbae0e7276c","description":"mój fragment kodu","content":"fmt.Println(\"Hello, world!\")"}
```
