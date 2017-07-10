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

## License

Copyright (c) Anish Athalye. Released under AGPLv3. See [LICENSE.txt](LICENSE.txt) for details.
