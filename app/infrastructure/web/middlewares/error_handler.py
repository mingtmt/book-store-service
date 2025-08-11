import logging
from fastapi.responses import JSONResponse
from starlette.middleware.base import BaseHTTPMiddleware
from app.core.exceptions import AppException

logger = logging.getLogger(__name__)

class ErrorHandlerMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request, call_next):
        try:
            return await call_next(request)
        except AppException as ex:
            logger.warning(f"[AppException] {ex.code}: {ex.message}")
            return JSONResponse(
                status_code=ex.status_code,
                content={"error": {"code": ex.code, "message": ex.message}}
            )
        except Exception as ex:
            logger.exception("[Unhandled] %s", ex)
            return JSONResponse(
                status_code=500,
                content={"error": {"code": "INTERNAL_ERROR", "message": "Internal server error"}}
            )
