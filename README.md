# Daily 150

Daily 150 is a journaling application that encourages users to write at least 150 words daily. Leveraging the power of Gemini 2.0 Flash, it automatically summarizes journal entries from the past week every Monday, providing users with a concise overview of their thoughts and experiences.

Beyond a standard CRUD application, Daily 150 was built with a strong focus on scalability, security, and robust architecture.

## Table of Contents

- [Daily 150](#daily-150)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Architecture Overview](#architecture-overview)
  - [Tech Stack](#tech-stack)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Clone the Repository](#clone-the-repository)
    - [Backend Setup (Go)](#backend-setup-go)
    - [Frontend Setup (React)](#frontend-setup-react)
    - [Running the Applications](#running-the-applications)
  - [Environment Variables](#environment-variables)
  - [Ethical Considerations](#ethical-considerations)
  - [License](#license)

## Features

*   **Daily Word Count Goal:** Encourages consistent journaling with a 150-word daily target.
*   **AI-Powered Summaries:** Utilizes Gemini 2.0 Flash to automatically summarize weekly journal entries.
*   **End-to-End Encryption:** Journal entries and AI-generated summaries are encrypted to ensure privacy.
*   **Decoupled Summary Service:** Summary generation logic is isolated into a separate service for enhanced scalability and maintainability.
*   **Redis Queue for Batch Processing:** Efficiently processes summary requests in batches via a Redis queue, managed by a Go server.
*   **Rate-Limited API Calls:** Implements token-based rate limiting on the summarization service to respect external API constraints.
*   **Robust Error Handling:** Incorporates retries and chunking mechanisms for graceful failure handling during API interactions.
*   **Chrome Extension Integration:** A companion Chrome extension blocks social media access until the user completes their daily journal entry, with status cached in Redis for performance.

## Architecture Overview

The Daily 150 system is designed for scalability and resilience, as illustrated in the diagram below:

![image](https://github.com/user-attachments/assets/c9c826c7-ecff-4f10-bb2c-77a1e7ad0ef0)

1.  **Weekly Trigger:** A cron job initiates the Go server every Monday.
2.  **Entry Collection:** The Go server fetches all journal entries from the previous week from the **Database**.
3.  **Task Enqueuing:** Instead of direct API calls, the Go server places summary generation tasks into a **Redis Queue**.
4.  **Background Processing:** A dedicated Go routine (worker process) continuously pulls tasks from the Redis queue.
5.  **Batched Processing:** The worker processes user journal entries in batches to optimize API calls to the summarization service.
6.  **Rate-Limited API Calls:** The **Express Server** (summarization service) manages rate limits when interacting with the **Gemini 2.0 Flash** API, ensuring compliance and stable operation.
7.  **Data Storage:** Once generated, the summaries are securely stored back in your **Database**.
8.  **Chrome Extension Interaction:** The Chrome extension interacts with the Go server to check journaling status, leveraging Redis for cached responses to minimize database load.

## Tech Stack

**Backend (Go Server)**
*   **Language:** Go
*   **Web Framework:** Fiber
*   **Database:** PostgreSQL
*   **ORM/Database Toolkit:** GORM
*   **Queue:** Redis
*   **Authentication:** JWT
*   **Encryption:** Custom End-to-End Encryption

**Frontend (React App)**
*   **Framework:** React
*   **Build Tool:** Vite
*   **Language:** TypeScript
*   **Styling:** Tailwind CSS

**Summarization Service (Express Server)**
*   **Language:** Node.js
*   **Web Framework:** Express.js
*   **AI Model Integration:** Google Gemini API
*   **Rate Limiting:** Custom token bucket implementation

**Other**
*   **Version Control:** Git
*   **Deployment:** Railway (as indicated by `0.0.0.0` in `app.Listen`)

## Getting Started

Follow these instructions to set up and run the Daily 150 project on your local machine.

### Prerequisites

*   **Go:** [Install Go](https://go.dev/doc/install) (version 1.20 or higher recommended)
*   **Node.js & npm:** [Install Node.js and npm](https://nodejs.org/en/download/) (LTS version recommended)
*   **PostgreSQL:** [Install PostgreSQL](https://www.postgresql.org/download/)
*   **Redis:** [Install Redis](https://redis.io/download/)
*   **Git:** [Install Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

### Clone the Repository

```bash
git clone https://github.com/actuallyakshat/daily-150.git
cd daily-150
```

### Backend Setup (Go)

1.  **Create `.env` file:**
    Create a file named `.env` in the root directory of the project and populate it with the following content. Replace placeholder values (especially for keys) with your own securely generated ones.

    ```dotenv
    PORT=8080
    GEMINI_KEY=YOUR_GEMINI_API_KEY # Obtain from Google AI Studio
    API_SECRET="your_secure_api_secret_key"
    SUMMARY_SERVICE_URL=http://localhost:3001/api/summary # Or the deployed URL if running separately
    JOURNAL_ENCRYPTION_KEY="your_secure_journal_encryption_key" # Must be 32 bytes for AES-256
    CRON_ACTIVATION_KEY="your_secure_cron_activation_key"
    SUMMARISER_KEY="your_secure_summariser_key"
    JWT_SECRET="your_long_and_secure_jwt_secret" # At least 256 bits (32 bytes) for HS256
    DATABASE_URL=postgresql://postgres:admin@localhost:5432/daily_150 # Update if your PG setup is different
    COOKIE_ENCRYPTION_KEY="your_secure_cookie_encryption_key" # A 32-byte key for AES-256
    ```

    **Important Security Note:** For `JOURNAL_ENCRYPTION_KEY`, `COOKIE_ENCRYPTION_KEY`, `API_SECRET`, `CRON_ACTIVATION_KEY`, `SUMMARISER_KEY`, and `JWT_SECRET`, it is crucial to generate strong, random keys. For AES-256, keys should be 32 bytes (256 bits). You can generate them using tools like OpenSSL or a programming language's cryptographically secure random number generator.

    Example for Go:
    ```go
    import (
        "crypto/rand"
        "encoding/base64"
        "fmt"
        "log"
    )

    func generateRandomBytes(n int) ([]byte, error) {
        b := make([]byte, n)
        _, err := rand.Read(b)
        if err != nil {
            return nil, err
        }
        return b, nil
    }

    func main() {
        key, err := generateRandomBytes(32) // 32 bytes for AES-256
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println("Base64 encoded key:", base64.StdEncoding.EncodeToString(key))
    }
    ```

2.  **Set up PostgreSQL Database:**
    Ensure your PostgreSQL server is running. Create a new database named `daily_150` (or whatever you've configured in `DATABASE_URL`).
    ```sql
    CREATE DATABASE daily_150;
    ```

3.  **Install Go Dependencies:**
    ```bash
    go mod tidy
    ```

4.  **Run Migrations:**
    The `init()` function in `main.go` will automatically run database migrations when the server starts.

### Frontend Setup (React)

1.  **Navigate to the client directory:**
    ```bash
    cd client
    ```

2.  **Install Node.js Dependencies:**
    ```bash
    npm install
    ```
    or
    ```bash
    yarn install
    ```

### Running the Applications

You need to run both the Go backend and the React frontend concurrently.

1.  **Start the Go Backend Server:**
    Open a new terminal, navigate to the project root (`daily-150/`), and run:
    ```bash
    go run main.go
    ```
    The server will start on the port specified in your `.env` file (default `8080`).

2.  **Start the React Frontend Development Server:**
    Open another terminal, navigate to the `client` directory (`daily-150/client`), and run:
    ```bash
    npm run dev
    ```
    or
    ```bash
    yarn dev
    ```
    The frontend will typically start on `http://localhost:5173`.

3.  **Access the Application:**
    Once both servers are running, open your web browser and navigate to `http://localhost:5173` (or the address where your frontend is served).

**Note:** If you are running the summarization service locally as a separate Node.js Express app, ensure it's also running, ideally on `http://localhost:3001` to match the default `SUMMARY_SERVICE_URL`.

## Ethical Considerations

While Daily 150 implements end-to-end encryption for journal entries, it's important to note that these entries are temporarily decrypted and processed by an external Large Language Model (LLM), Gemini 2.0 Flash, for summarization. Users who prioritize absolute privacy and wish to avoid any external processing of their sensitive data might have concerns.

A more privacy-conscious solution would involve running the LLM locally within the project. However, the compute requirements for hosting such a powerful model are currently beyond the scope and budget of this project. This remains a significant future goal, as we explore ways to make the system more self-sufficient and enhance user data privacy further.

---
