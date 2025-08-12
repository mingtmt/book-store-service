# Book Store Service

A backend service for a book store, built with FastAPI, SQLAlchemy, and Alembic.

## Features
- User registration and authentication (JWT)
- Modular architecture (domain, infrastructure, use_cases, schemas)
- Database migrations with Alembic
- Dockerized for easy setup

## Requirements
- Python 3.12+
- Docker & Docker Compose (for containerized setup)

## Setup (Local)

1. **Clone the repository:**
   ```bash
   git clone <your-repo-url>
   cd book-store-service
   ```

2. **Create a virtual environment and activate it:**
   ```bash
   python3 -m venv venv
   source venv/bin/activate
   ```

3. **Install dependencies:**
   ```bash
   pip install -r requirements.txt
   ```

4. **Set up environment variables:**
   - Copy `.env.example` to `.env` and update values as needed.

5. **Run database migrations:**
   ```bash
   alembic upgrade head
   ```

6. **Start the application:**
   ```bash
   uvicorn app.main:app --reload
   ```

## Setup (Docker)

1. **Build and start services:**
   ```bash
   docker-compose up --build
   ```

2. The API will be available at `http://localhost:8000`

## API Documentation
- Swagger UI: [http://localhost:8000/docs](http://localhost:8000/docs)
- Redoc: [http://localhost:8000/redoc](http://localhost:8000/redoc)

## Project Structure
```
app/
  main.py
  api/
  core/
  domain/
  infrastructure/
  schemas/
  use_cases/
alembic/
requirements.txt
docker-compose.yml
```

## License
MIT
