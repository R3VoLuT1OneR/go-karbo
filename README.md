# Karbo Go
Golang's implementation of [Karbo](https://github.com/Karbovanets/karbo) cryptocurrency.

### Running

```shell
go run krbd.go
```

## Development Notes

### Development Issues
List of tasks to be done.

#### Cryptonote
  * Implement blockchain "inmemory" store for usage in unit tests and define interface for blockchain storage

#### Crypto
  * Implement crypto functions that we can find in test data file found in C++ implementation
    [crypto/fixtures/tests.txt]()

#### P2P
  * P2P: Handle incoming connections
  * P2P: Create proper logging in p2p node
  * Transaction serialize/deserialize signatures