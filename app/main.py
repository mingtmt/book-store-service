from dotenv import load_dotenv
load_dotenv()

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
from fastapi import status
from app.api.v1.routes import api_router
from app.infrastructure.web.middlewares.error_handler import ErrorHandlerMiddleware

app = FastAPI(
	title="Book Store Service",
	version="1.0.0",
	description="Backend service for a book store."
)

# Register middlewares
app.add_middleware(ErrorHandlerMiddleware)
app.add_middleware(
	CORSMiddleware,
	allow_origins=["*"],
	allow_credentials=True,
	allow_methods=["*"],
	allow_headers=["*"],
)

# Health check endpoint
@app.get("/ping", tags=["Health"])
async def ping():
	return JSONResponse(content={"message": "pong"}, status_code=status.HTTP_200_OK)

# Register API routes
app.include_router(api_router, prefix="/api/v1")
