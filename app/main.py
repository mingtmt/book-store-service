from dotenv import load_dotenv
load_dotenv()

from fastapi import FastAPI
from app.api.v1.routes import api_router
from app.infrastructure.web.middlewares.error_handler import ErrorHandlerMiddleware

app = FastAPI()

# Register middlewares
app.add_middleware(ErrorHandlerMiddleware)

# Register API routes
app.include_router(api_router, prefix="/api/v1")
