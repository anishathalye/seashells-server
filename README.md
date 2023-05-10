# Seashells

The source code for the [seashells.io](https://seashells.io/) server. See the [anishathalye/seashells](https://github.com/anishathalye/seashells) repo for the client.

For more information, see [seashells.io](https://seashells.io) or the [launch blog post](https://www.anishathalye.com/2017/07/10/seashells/).

![Seashells.io website preview](https://raw.githubusercontent.com/anishathalye/assets/master/seashells/seashells-web.png)

## Build and run

### Development mode

1. Build the binary with `go build`
2. Run `./seashells-server`

### Production mode

1. Build the binary with `go build`
2. Configure some settings, by `cp env.sample env` and editing the resulting `env` file
3. Run with `./run.bash`

### Environment Variables

- `NETCAT_BINDING`: interface/port to listen for terminal connections.  
  Defaults to `:1337`. 📝 This should be directly accessible to the network/internet.
- `WEBAPP_BINDING`: interface/port to listen for web app connections.  
  Defaults to `:8080`. 📝 _Strongly_ consider using `127.0.0.1:port` and using a secure reverse proxy.
- `BASE_URL`: URL prefix to print out with randomized connection path string when terminal connections are established.  
  Defaults to `https://seashells.io/v/`. 📝 You will almost certainly want to change this.  
- `GIN_MODE`: [Gin web framework](https://gin-gonic.com/) mode (`debug` or `release`).  
  Defaults to `debug`.
- `GTAG`: Header tag to pass along with requests.  
  Defaults to `g-tag`. # 🚨 CHANGE THIS!
- `ADMIN_PASSWORD`: Password to use for all admin actions.  
  Defaults to `xxx`. # 🚨 CHANGE THIS!

### Paths

- `/v/:id` the default endpoint provided after establishing a terminal session (full xterm.js terminal emulation)
- `/p/:id` plaintext version of ^^
- `/inspect` admin login endpoint to see an overview of sessions username is `admin` and password is in `ADMIN_PASSWORD` env var

## License

Copyright (c) Anish Athalye. Released under AGPLv3. See [LICENSE.txt](LICENSE.txt) for details.
