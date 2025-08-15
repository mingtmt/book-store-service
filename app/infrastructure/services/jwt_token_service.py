from jose import jwt
from datetime import datetime, timedelta
from typing import Optional
from app.ports.token_service import ITokenService

from app.settings.config import settings

class JWTTokenService(ITokenService):
    def create_access_token(self, user_id: int) -> str:
        expire = datetime.utcnow() + timedelta(minutes=settings.access_token_expire_minutes)
        payload = {
            "sub": str(user_id),
            "exp": expire
        }
        return jwt.encode(payload, settings.jwt_secret_key, algorithm=settings.jwt_algorithm)

    def decode_access_token(self, token: str) -> Optional[int]:
        try:
            payload = jwt.decode(token, settings.jwt_secret_key, algorithms=[settings.jwt_algorithm])
            user_id = payload.get("sub")
            if user_id is None:
                return None
            return int(user_id)
        except (jwt.ExpiredSignatureError, jwt.InvalidTokenError):
            return None
