# metadata-engine-crawler

A CLI tool that crawls a file system and sends file/directory metadata to the metadata-engine API.

Each crawled directory appears as a separate volume in the web UI. The crawler stores paths relative to the scanned directory, so no host filesystem structure is exposed.

---

## Installation

### macOS

1. Install Go using Homebrew:
   ```bash
   brew install go
   ```

2. Verify the installation:
   ```bash
   go version
   ```

3. Clone this repository:
   ```bash
   git clone https://github.com/seanluce/metadata-engine-crawler.git
   cd metadata-engine-crawler
   ```

4. Build the binary:
   ```bash
   go build -o crawler .
   ```

5. Run the crawler:
   ```bash
   ./crawler --root /path/to/scan --api https://your-api-url.com
   ```

### Linux

1. Download and install Go (replace the version number with the latest from https://go.dev/dl/):
   ```bash
   wget https://go.dev/dl/go1.22.5.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz
   ```

2. Add Go to your PATH by adding this line to `~/.bashrc` or `~/.zshrc`:
   ```bash
   export PATH=$PATH:/usr/local/go/bin
   ```

3. Reload your shell:
   ```bash
   source ~/.bashrc
   ```

4. Verify the installation:
   ```bash
   go version
   ```

5. Clone this repository:
   ```bash
   git clone https://github.com/seanluce/metadata-engine-crawler.git
   cd metadata-engine-crawler
   ```

6. Build the binary:
   ```bash
   go build -o crawler .
   ```

7. Run the crawler:
   ```bash
   ./crawler --root /path/to/scan --api https://your-api-url.com
   ```

### Windows

1. Download the Go installer from https://go.dev/dl/ (choose the `.msi` file for Windows).

2. Run the installer and follow the prompts. The default settings are fine.

3. Open a new Command Prompt or PowerShell window and verify the installation:
   ```cmd
   go version
   ```

4. Clone this repository:
   ```cmd
   git clone https://github.com/seanluce/metadata-engine-crawler.git
   cd metadata-engine-crawler
   ```

5. Build the executable:
   ```cmd
   go build -o crawler.exe .
   ```

6. Run the crawler:
   ```cmd
   crawler.exe --root C:\path\to\scan --api https://your-api-url.com
   ```

---

## Options

| Flag | Env Var | Required | Default | Description |
|------|---------|----------|---------|-------------|
| `--root` | | Yes | | Root path to crawl |
| `--api` | `API_URL` | Yes | | API base URL |
| `--name` | | No | Root directory name | Volume name displayed in the web UI |
| `--version` | | No | | Version label for this crawl run (e.g. `snap001`) |
| `--workers` | | No | 8 | Number of concurrent workers |

The `--api` flag takes precedence over the `API_URL` environment variable. If neither is provided, the crawler will exit with an error.

### Volume naming

By default, the volume name is the name of the directory being scanned. For example, scanning `/mnt/data` creates a volume called `data`.

Use `--name` to override this. This is useful when the directory name isn't meaningful (e.g. a Windows drive letter for a mounted network share):

**macOS / Linux:**
```bash
./crawler --root /mnt/share --name my-volume --api https://your-api-url.com
```

**Windows:**
```cmd
crawler.exe --root Z:\ --name my-volume --api https://your-api-url.com
```

The volume will appear in the web UI at `https://your-web-url.com/my-volume`.

### Version labeling

Use `--version` to label a crawl run. Each volume maintains its own version history, and the label appears in the version selector dropdown in the web UI.

```bash
./crawler --root /mnt/share --name my-volume --version "snap001" --api https://your-api-url.com
```

Running multiple crawls with different version labels lets you compare the state of a volume at different points in time.

### Using environment variables instead of flags

**macOS / Linux:**
```bash
export API_URL=https://your-api-url.com
./crawler --root /path/to/scan
```

**Windows (Command Prompt):**
```cmd
set API_URL=https://your-api-url.com
crawler.exe --root C:\path\to\scan
```

**Windows (PowerShell):**
```powershell
$env:API_URL = "https://your-api-url.com"
.\crawler.exe --root C:\path\to\scan
```

---

## Docker

If you prefer not to install Go, you can run the crawler using Docker:

```bash
docker build -t metadata-engine-crawler .

docker run --rm \
  -e API_URL=https://your-api-url.com \
  -v /path/to/scan:/data:ro \
  metadata-engine-crawler --root /data --name my-volume
```
