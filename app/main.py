from dotenv import load_dotenv
from fastapi import FastAPI, status
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse

from app.bootstrap.logging_config import setup_logging
from app.presentation.http.api.v1.routes import api_router
from app.presentation.http.middlewares.error_handler import ErrorHandlingMiddleware
from app.presentation.http.middlewares.logger import LoggingMiddleware


def create_app() -> FastAPI:
    load_dotenv()
    setup_logging()

    app = FastAPI(
        title="Book Store Service",
        version="1.0.0",
        description="Backend service for a book store.",
    )

    # Middlewares
    app.add_middleware(LoggingMiddleware)
    app.add_middleware(ErrorHandlingMiddleware)
    app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"],
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )

    # Routes
    app.include_router(api_router, prefix="/api/v1")

    @app.get("/ping", tags=["Health"])
    async def ping():
        return JSONResponse(content={"message": "pong"}, status_code=status.HTTP_200_OK)

    return app


app = create_app()
