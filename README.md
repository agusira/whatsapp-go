# WhatsGO - Multi-Module Message Handler Application

## **Project Overview**

**WhatsGO** is a comprehensive Go-based messaging application that integrates with WhatsApp Web API and provides advanced event handling capabilities. The system is designed to handle real-time message processing, session management, call handling, and plugin extensibility.

### Key Features
- **Real-time Messaging**: Handle incoming/outgoing messages via WhatsApp Web API
- **Session Management**: Persistent device sessions with SQLite storage
- **Call Handling**: Automated call answering with audio playback
- **Plugin Architecture**: Extensible plugin system for custom functionality
- **Event-Driven**: Reactive architecture for handling WhatsApp events

---

## **Project Structure**

```
whatsapp-go/
├── main.go                    # Entry point and core application logic
├── go.mod                     # Go module dependencies
├── go.sum                     # Dependency checksums
├── .env                       # Environment configuration
├── .gitignore                 # Git ignore rules
├── call.mp3                   # Call audio playback
├── tes.db                     # SQLite database for sessions
├── configs/                   # Configuration files
│   ├── env.go                 # Environment variable loading
│   └── json.go                # JSON configuration handling
├── handler/                   # Event handlers and routing logic
│   └── handler.go             # Main event handler
├── lib/                       # Core library modules
│   ├── client.go              # Client connection management
│   ├── plugins.go             # Plugin loading and management
│   ├── serialize.go           # Data serialization utilities
│   └── tiktok.go               # TikTok-specific functionality
├── plugins/                   # Plugin ecosystem
│   ├── group/                 # Group management plugins
│   ├── other/                 # Miscellaneous plugins
│   ├── owner/                 # Owner-specific plugins
│   └── plugins.go             # Plugin entry point
└── docs/                      # Documentation directory (README placeholder)
```

---

## **Installation & Setup**

### Prerequisites
- Go 1.25.0 or higher
- Git
- WhatsApp Web API access

### Installation Steps

1. **Clone the repository**:
   ```bash
   git clone https://github.com/agusira/whatsapp-go
   cd whatsapp-go
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**:
   ```bash
   # Copy .env.example to .env if it exists
   cp .env.example .env  # (if available)
   # Edit .env with your configuration
   ```

4. **Build the application**:
   ```bash
   go build -o whatsapp-go ./...
   ```

---

## **Usage**

### Basic Operations

Run the application:
```bash
./whatsapp-go
```

The application will start with these steps:

1. **Session Selection**: Choose between creating a new session or loading an existing one
2. **Authentication**: If no session exists, you'll be prompted to enter your WhatsApp number
3. **Pairing**: QR code will be generated for WhatsApp pairing
4. **Operation**: The application will connect and await events

### Available Commands

While running, the application accepts:
- **New Session**: Creates a fresh device session
- **Load Session**: Restores a previously saved session
- **Graceful Shutdown**: Ctrl+C to disconnect and exit

---

## **Configuration**

### Environment Variables
| Variable | Description | Default |
|----------|-------------|---------|
| `LOGGER` | Log level for WhatsApp and database logs | unset |
| `DB_ADDRESS` | Path to SQLite database file | `./tes.db` |
| `DEVICE_NAME` | Name of the connected device | `Agus` |

### Database
- Uses SQLite with foreign key constraints enabled
- Stores device sessions and connection state
- File location configurable via `DB_ADDRESS`

---

## **Features & Capabilities**

### Message Handling
- Receive and process incoming messages
- Forward messages to event handlers
- Serialize message data for external processing

### Call Handling
- Answer incoming calls automatically
- Play audio notification (`call.mp3`)
- Hang up calls after audio playback

### Plugin System
- Dynamic plugin loading
- Support for group, owner, and custom plugins
- Plugin state management

### Serialization
- Client serialization for external use
- Message serialization for processing
- Integration with handler modules

---

## **Development**

### Building from Source
```bash
# Build for current platform
go build -o whatsapp-go

# Build for specific platform (example: linux)
go build -o whatsapp-go-linux -tags=linux
```

### Running Tests
```bash
# Run unit tests
go test ./...

# Run specific package tests
go test ./lib
```

### Plugin Development
Plugins should follow the structure defined in `plugins/plugins.go`. New plugins can be added to the `plugins/` subdirectory according to their type (group, owner, other).

---

## **Contributing**

### Code Quality Standards
- Follow Go formatting standards (`gofmt`)
- Use meaningful variable and function names
- Add comments for complex logic
- Maintain consistent error handling patterns

### Pull Request Process
1. **Fork** the repository
2. **Create** a feature branch
3. **Implement** changes with tests
4. **Document** new functionality
5. **Submit** pull request

### Testing Standards
- Write unit tests for new functionality
- Ensure existing tests pass
- Test edge cases and error conditions

---

## **Dependencies**

### Direct Dependencies
- `github.com/mattn/go-sqlite3` - SQLite driver
- `github.com/purpshell/meowcaller` - Call handling
- `go.mau.fi/whatsmeow` - WhatsApp Web API client
- `github.com/joho/godotenv` - Environment variable management
- `github.com/HugoSmits86/nativewebp` - WebP image handling
- `github.com/PuerkitoBio/goquery` - HTML parsing
- `github.com/robertkrimen/otto` - JavaScript runtime

### Development Tools
- Go 1.25.0 or higher
- Git

---

## **Troubleshooting**

### Common Issues

**Q: What if the application fails to start?**
- Check `.env` file configuration
- Verify database file permissions
- Ensure all Go dependencies are installed

**Q: How to debug issues?**
- Set `LOGGER` environment variable to `debug`
- Use `go test -v` for verbose test output
- Check application logs for detailed error information

**Q: Why does the application exit unexpectedly?**
- Incomplete session configuration
- Network connectivity issues
- Authentication failures

### Get Help
For additional help:
- Check the `docs/` directory for detailed guides
- Review the source code comments
- Join the community forums or Discord server

---

## **License**

This project is licensed under the MIT License.

```
MIT License

Copyright (c) [2026] [Agus Irawan]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR
IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```

---

## **Contact & Support**

### Communication
- GitHub Issues: For bug reports and feature requests
- Discord: Community discussions and support
- Email: agusirawan2834@gmail.com

### Social Media
- TikTok: @ikiagusss
- IG: @guzagussss
---

*Version: 1.0.0*
