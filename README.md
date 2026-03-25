# FreightShip Backend

A Go (Gin) backend API for the FreightShip platform — handling MBL/HBL documents, bookings, shipments, and PDF generation.

## Prerequisites

- **Go 1.24+**
- **MongoDB** (Atlas or local instance)

## Configuration

All configuration is defined in `config.yaml`:

```yaml
mongo:
  uri: "mongodb+srv://..."
  database: "fs-backend-db"
server:
  port: ":5000"
pdf_service:
  base_url: "http://localhost:3000"
extraction_service:
  base_url: "http://localhost:10000/extract"
```

## Running Locally

```bash
go run main.go
```

The server will start at `http://localhost:5000`.

---

## 🐳 Run with Docker

The easiest way to run the backend on any machine — no need to install Go or configure a build toolchain.

> **Note:** The `config.yaml` file is **NOT** baked into the Docker image for security. You must provide it at runtime using a volume mount.

### Prerequisites

- [Docker](https://www.docker.com/products/docker-desktop) installed on your system.

### Step 1 — Pull the Image

```bash
docker pull jainras/freightship-backend
```

### Step 2 — Create a `config.yaml` File

Create a file named `config.yaml` in your working directory with the following content:

```yaml
mongo:
  uri: "mongodb+srv://your_mongo_connection_string"
  database: "fs-backend-db"
server:
  port: ":5000"
pdf_service:
  base_url: "http://localhost:3000"
extraction_service:
  base_url: "http://localhost:10000/extract"
```

### Step 3 — Run the Container

```bash
docker run -d -p 5000:5000 -v ./config.yaml:/app/config.yaml --name fs-backend jainras/freightship-backend
```

| Flag | Description |
|---|---|
| `-d` | Run in background (detached mode) |
| `-p 5000:5000` | Map host port 5000 → container port 5000 |
| `-v ./config.yaml:/app/config.yaml` | Mount your local config file into the container |
| `--name fs-backend` | Assign a name to the container |

### Step 4 — Verify It's Running

Open your browser and visit:

```
http://localhost:5000
```

Or check from the terminal:

```bash
docker ps
```

### Useful Docker Commands

```bash
# View logs
docker logs fs-backend

# Follow logs in real-time
docker logs -f fs-backend

# Stop the container
docker stop fs-backend

# Restart the container
docker start fs-backend

# Remove the container (must be stopped first)
docker rm fs-backend
```

### Updating the Image After Code Changes

After making changes to the code, rebuild and push the updated image:

```bash
# 1. Rebuild the image
docker build -t jainras/freightship-backend:latest .

# 2. Push the updated image to Docker Hub
docker push jainras/freightship-backend:latest
```

On the target machine, pull and restart with the latest version:

```bash
# 3. Pull the latest image
docker pull jainras/freightship-backend:latest

# 4. Stop and remove the old container
docker stop fs-backend
docker rm fs-backend

# 5. Run the new version
docker run -d -p 5000:5000 -v ./config.yaml:/app/config.yaml --name fs-backend jainras/freightship-backend:latest
```
