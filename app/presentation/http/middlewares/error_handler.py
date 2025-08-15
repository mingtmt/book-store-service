from typing import Callable, Awaitable, Any
import logging
from fastapi.responses import JSONResponse
from fastapi.exceptions import RequestValidationError
from pydantic import ValidationError
from starlette.exceptions import HTTPException as StarletteHTTPException

from app.domain.errors import (
    DomainError, UserNotFound, BookNotFound, EmailAlreadyExists,
    InvalidEmail, PasswordTooWeak, InvalidCredentials,
)
from app.usecases.errors import (
    UseCaseError, NotFound, Unauthorized, BadRequest, Conflict, Forbidden,
)

logger = logging.getLogger(__name__)

DOMAIN_HTTP_MAP: dict[type[Exception], tuple[int, str]] = {
    UserNotFound: (404, "USER_NOT_FOUND"),
    BookNotFound: (404, "BOOK_NOT_FOUND"),
    EmailAlreadyExists: (409, "EMAIL_ALREADY_EXISTS"),
    InvalidEmail: (400, "INVALID_EMAIL"),
    PasswordTooWeak: (400, "PASSWORD_TOO_WEAK"),
    InvalidCredentials: (401, "INVALID_CREDENTIALS"),
}
USECASE_HTTP_MAP: dict[type[Exception], tuple[int, str]] = {
    NotFound: (404, "NOT_FOUND"),
    Unauthorized: (401, "UNAUTHORIZED"),
    Forbidden: (403, "FORBIDDEN"),
    BadRequest: (400, "BAD_REQUEST"),
    Conflict: (409, "CONFLICT"),
}

def _match_map(exc: Exception) -> tuple[int, str] | None:
    for et, payload in DOMAIN_HTTP_MAP.items():
        if isinstance(exc, et):
            return payload
    for et, payload in USECASE_HTTP_MAP.items():
        if isinstance(exc, et):
            return payload
    return None


class ErrorHandlingMiddleware:
    """
    ASGI middleware để bắt mọi exception và trả JSON đồng nhất.
    Dùng ASGI thuần để tránh một số corner cases của BaseHTTPMiddleware.
    """
    def __init__(self, app):
        self.app = app

    async def __call__(self, scope: dict, receive: Callable[[], Awaitable[dict]], send: Callable[[dict], Awaitable[Any]]):
        # Chỉ xử lý HTTP; bỏ qua lifespan, websockets...
        if scope.get("type") != "http":
            await self.app(scope, receive, send)
            return

        try:
            await self.app(scope, receive, send)
        except Exception as exc:
            # 1) Map Domain/Usecase errors
            mapped = _match_map(exc)
            if mapped is not None:
                status, code = mapped
                logger.info("Handled business error %s -> %s", exc.__class__.__name__, code)
                response = JSONResponse(status_code=status, content={"code": code, "message": str(exc)})
                await response(scope, receive, send)
                return

            # 2) Validation (Pydantic/FastAPI)
            if isinstance(exc, (RequestValidationError, ValidationError)):
                response = JSONResponse(status_code=422, content={"code": "VALIDATION_ERROR", "message": str(exc)})
                await response(scope, receive, send)
                return

            # 3) Starlette/FastAPI HTTPException passthrough theo format chung
            if isinstance(exc, StarletteHTTPException):
                response = JSONResponse(status_code=exc.status_code, content={"code": "HTTP_EXCEPTION", "message": str(exc.detail)})
                await response(scope, receive, send)
                return

            # 4) Fallback 500
            logger.exception("Unhandled server error")
            response = JSONResponse(status_code=500, content={"code": "INTERNAL_ERROR", "message": "Internal server error"})
            await response(scope, receive, send)
