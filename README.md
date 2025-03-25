# Shorten

A simple URL shortener built with **Svelte** for the frontend and **Go Fiber + SQLite** for the backend.

## Features
- 🔗 Shorten long URLs instantly
- 📋 Copy shortened URLs easily
- 📊 Track the number of visits (coming soon)
- ⚡️ Fast and lightweight

## Tech Stack
- **Frontend:** SvelteKit, Tailwind CSS
- **Backend:** Go Fiber, SQLite, Redis (for caching)
- **Deployment:** Docker, Nginx (optional)

## Installation & Setup

### 1️⃣ Backend Setup (Go Fiber + SQLite)
```sh
# Clone the repository
git clone https://github.com/yourusername/shorten.git
cd shorten/backend

# Install dependencies
go mod tidy

# Run the server
go run main.go
```

### 2️⃣ Frontend Setup (Svelte)
```sh
cd frontend
npm install
npm run dev
```

### 3️⃣ Deploy with Docker
```sh
docker-compose up -d
```

## API Endpoints
| Method | Endpoint       | Description            |
|--------|--------------|------------------------|
| POST   | `/shorten`   | Shorten a URL         |
| GET    | `/:code`     | Redirect to full URL  |

## License
MIT License © 2025 [Your Name]

