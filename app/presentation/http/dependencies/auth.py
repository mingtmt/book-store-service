import uuid

from fastapi import Depends, HTTPException, Request, Security, status
from fastapi.security import HTTPAuthorizationCredentials, HTTPBearer
from sqlalchemy.orm import Session

from app.domain.errors import AccessTokenExpired, AccessTokenInvalid
from app.infrastructure.db.sqlalchemy.repositories.user_impl import (
    SqlAlchemyUserRepository,
)
from app.infrastructure.services.jwt_token_service import JWTTokenService
from app.presentation.http.dependencies.db import get_db

token_service = JWTTokenService()
bearer_scheme = HTTPBearer(auto_error=True)


async def get_current_user(
    creds: HTTPAuthorizationCredentials = Security(bearer_scheme),
    db: Session = Depends(get_db),
):
    if creds is None or (creds.scheme or "").lower() != "bearer":
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED, detail="Not authenticated"
        )

    token = creds.credentials
    try:
        claims = JWTTokenService().decode_access_token(token)
        user_id = uuid.UUID(claims["sub"])
    except Exception:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED, detail="Invalid token"
        )

    user = SqlAlchemyUserRepository(db).get_by_id(user_id)
    if not user:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED, detail="User not found"
        )

    return user


def require_auth(
    credentials: HTTPAuthorizationCredentials = Security(bearer_scheme),
    request: Request = None,
):
    # Swagger now passes the parsed credentials here
    token = credentials.credentials.strip()
    try:
        payload = token_service.decode_access_token(token)
    except AccessTokenExpired:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED, detail="Token has expired"
        )
    except AccessTokenInvalid:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED, detail="Invalid token"
        )

    request.state.user = payload  # optional: attach context
    return payload
