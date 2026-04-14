# Yalsh-Go (Yet Another Link Shortener)

A production-ready, high-performance URL shortener built with **Go** and **PostgreSQL**. This backend service handles user authentication, custom URL aliasing, and link click tracking, all orchestrated within a containerized infrastructure.

## 🚀 Key Features

- **Custom URL Aliasing**: Users can create short links with custom keys.
- **User Authentication**: Secure JWT-based authentication for user account management.
- **Automated SSL/TLS**: Every deployment is secured with automated HTTPS via **Caddy**.
- **Real-time Analytics**: Tracks clicks on all shortened links.
- **Containerized Deployment**: Fully orchestrated with **Docker** and **Docker Compose**.
- **GitOps CI/CD**: Fully automated Build, Push, and Deploy pipeline via **GitHub Actions**.

## 🛠️ Technology Stack

- **Backend**: [Go (Golang)](https://go.dev/) with `go-chi` router.
- **Database**: [PostgreSQL](https://www.postgresql.org/) (Hosted on [Neon DB](https://neon.tech/)).
- **Reverse Proxy**: [Caddy](https://caddyserver.com/) (Handles SSL/TLS).
- **DevOps**: [Docker](https://www.docker.com/), [GitHub Actions](https://github.com/features/actions).
- **Infrastructure**: [Vultr](https://www.vultr.com/) (Rocky Linux 9), [Cloudflare](https://www.cloudflare.com/) (DNS).
- **IaC**: [Terraform](https://www.terraform.io/).

## 🚦 Getting Started (Local)

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/Kenwanjohi/yalsh-go.git
    cd yalsh-go
    ```
2.  **Set up your environment**:
    Copy `.env.example` to `.env` and fill in your local PostgreSQL connection string and JWT secret.
3.  **Run with Docker Compose**:
    ```bash
    docker-compose up -d --build
    ```
    The API will be available at `http://localhost:3000`.

## 🚢 Deployment (CI/CD)

Deployments are fully automated. Every push to the `main` branch triggers a GitHub Actions workflow that:
1.  **Builds** the Go binary and Docker image.
2.  **Pushes** the image to Docker Hub.
3.  **Copies** `docker-compose.yml` and `Caddyfile` to the Vultr server.
4.  **Deploys** the new image to the production environment on Rocky Linux.

### Required GitHub Secrets:
- `DOCKER_USERNAME`: Your Docker Hub username.
- `DOCKER_TOKEN`: Your Docker Hub Access Token.
- `SERVER_IP`: Your Vultr server's IP address.
- `SERVER_SSH_KEY`: Your private SSH key for server access.

## 📡 API Endpoints (v1)

- `POST /v1/accounts`: Register a new user account.
- `POST /v1/session`: Login and receive an authentication cookie.
- `POST /v1/links`: Create a new shortened link (Authenticated).
- `GET /v1/links`: List all links for the authenticated user.
- `GET /{key}`: Redirect from a short key to the original URL.
- `GET /v1/healthcheck`: Verify the API status and environment.

## 📄 License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
