
<!-- PROJECT LOGO -->
<h1 align="center">TBL</h1>
<p align="center">
  package for unmarshal Diablo 2 TextTables from byte-encoded ".tbl" files
  <br />
  <br />
  <a href="https://github.com/gravestench/tbl_text/issues">Report Bug</a>
  ·
  <a href="https://github.com/gravestench/tbl_text/issues">Request Feature</a>
</p>

<!-- ABOUT THE PROJECT -->
## About

The Diablo 2 TextTable Unmarshaller package provides a Go implementation to
unmarshal Diablo 2 TextTables from byte-encoded ".tbl" files. These files are
used for storing locale-specific translations of game text in Diablo 2, among 
other things.

The main data structure used to hold the key-value pairs is the `TextTable`.
Each key in the `TextTable` corresponds to a text identifier, while the value
corresponds to the translated text for that identifier.

## Usage

### Prerequisites
To use this TextTable unmarshaller package, ensure you have Go 1.16 or a later
version installed, and your Go environment is set up correctly.

### Installation
To install the package, you can use Go's standard `go get` command:

```shell
go get -u github.com/gravestench/tbl_text
```

### Unmarshal TextTable
To unmarshal a Diablo 2 TextTable from a byte slice, use the `Unmarshal`
function:

```golang
fileData := // Load your ".tbl" file data here as a byte slice
textTable, err := Unmarshal(fileData)
if err != nil {
    // Handle error
}
// Use the textTable object to access the key-value pairs
```

### Accessing Key-Value Pairs
After unmarshalling the TextTable, you can access the individual key-value pairs:

```golang
// Access the key-value pairs within the TextTable
for key, value := range textTable {
    // Use the key and value for your translation needs
    fmt.Printf("Key: %s, Value: %s\n", key, value)
}
```

### Features
The Diablo 2 TextTable Unmarshaller package offers the following features:
- Efficiently unmarshal Diablo 2 TextTables from ".tbl" files.
- Access individual key-value pairs in the TextTable.
- Provides a convenient `TextTable` data structure for working with translations.

<!-- CONTRIBUTING -->
## Contributing

Contributions to the Diablo 2 TextTable Unmarshaller package are welcome and
encouraged. If you find any issues or have improvements to suggest, feel free to
open an issue or submit a pull request.

To contribute to the project, follow these steps:

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<!-- MARKDOWN LINKS & IMAGES -->
[tbl_text]: https://github.com/gravestench/tbl_text
