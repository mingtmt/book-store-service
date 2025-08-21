from fastapi import Request
from starlette.middleware.base import BaseHTTPMiddleware

from app.bootstrap.logging_config import logger


class LoggingMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next):
        client_ip = request.client.host
        method = request.method
        url = request.url.path

        logger.info(f"Request: {method} {url} from {client_ip}")

        try:
            response = await call_next(request)
        except Exception as exc:
            logger.error(
                f"Error handling request: {method} {url} from {client_ip} - {exc}"
            )
            raise

        status_code = response.status_code
        if 400 <= status_code < 500:
            logger.warning(
                f"Response: {method} {url} returned {status_code} to {client_ip}"
            )
        elif status_code >= 500:
            logger.error(
                f"Response: {method} {url} returned {status_code} to {client_ip}"
            )
        else:
            logger.info(
                f"Response: {method} {url} returned {status_code} to {client_ip}"
            )

        return response
