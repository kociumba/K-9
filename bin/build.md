you can build K-9 with `go build -C . -o ./bin -ldflags "-s -w"`

if you have task installed just use `task build-r`

building with just `go build` works but the binary is 2mb bigger 🤷

if you are on windows .syso files are provided so the app will compile with an icon and metadata
