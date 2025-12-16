# bytes-human

Convert between raw byte counts and human-readable sizes with precision control and multiple unit standards

## Features

- Convert raw byte counts to human-readable format (e.g., 1024 -> 1.0 KiB)
- Parse human-readable sizes back to bytes (e.g., '1.5 GB' -> 1500000000)
- Support SI decimal units (KB, MB, GB using 1000) and binary units (KiB, MiB, GiB using 1024)
- Configurable precision for decimal places (0-6)
- Auto-detect optimal unit or force specific unit output
- Batch processing from stdin or files (one value per line)
- Round up, down, or to nearest based on flags
- Handle very large numbers (up to uint64 max)
- Validate input format and provide clear error messages
- Output in plain text or JSON format for scripting

## Installation

```bash
# Clone the repository
git clone https://github.com/KurtWeston/bytes-human.git
cd bytes-human

# Install dependencies
go build
```

## Usage

```bash
./main
```

## Built With

- go

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
