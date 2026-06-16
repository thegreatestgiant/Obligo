# Deploying Charity-Tracker

Charity-Tracker is designed to be completely self-contained and simple to deploy using Docker Compose.

## Prerequisites

* Docker and Docker Compose installed on your host machine.

## Quickstart

**1. Clone the repository:**
\`\`\`bash
git clone <https://github.com/yourusername/Charity-Tracker.git>
cd Charity-Tracker
\`\`\`

**2. Configure your Environment Variables:**
Copy the example environment file and fill in your details.
\`\`\`bash
cp .env.example .env
\`\`\`

Inside your `.env` file, ensure the following variables are set.
*Note: Your `JWT_SECRET` must be a long, highly secure string. You can generate a standard secure string by running `openssl rand -base64 64` in your terminal.*

\`\`\`env

# .env

APP_URL=<http://localhost:8080>

DB_USER=charity_user
DB_PASSWORD=super_secure_db_password
DB_NAME=charity_db

JWT_SECRET=your_super_long_random_string_here
\`\`\`

**3. Launch the Application:**
\`\`\`bash
docker-compose up -d --build
\`\`\`

**4. Access the App:**
Navigate to your configured `APP_URL` (default: `http://localhost:8080`).

*Note: Database schemas are automatically generated and migrated on startup.*
