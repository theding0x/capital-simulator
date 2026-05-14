# Go Project Structure

> Load when: Project layout, packages, testing conventions.

## Standard Layout

```
myapp/
├── cmd/
│   └── server/
│       └── main.go          # Entry point
├── internal/                 # Private packages
│   ├── handler/              # HTTP handlers
│   ├── service/              # Business logic
│   ├── repository/           # Data access
│   └── model/                # Domain types
├── pkg/                      # Public packages (reusable)
├── go.mod
├── go.sum
├── Makefile
└── Dockerfile
```

## Testing

```go
func TestGetUser(t *testing.T) {
    tests := []struct {
        name    string
        id      string
        want    *User
        wantErr error
    }{
        {name: "found", id: "1", want: &User{ID: "1", Name: "Alice"}, wantErr: nil},
        {name: "not found", id: "999", want: nil, wantErr: ErrNotFound},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := repo.GetUser(tt.id)
            if !errors.Is(err, tt.wantErr) {
                t.Errorf("error = %v, want %v", err, tt.wantErr)
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("got = %v, want %v", got, tt.want)
            }
        })
    }
}
```
