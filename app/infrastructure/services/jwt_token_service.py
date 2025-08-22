# app/infrastructure/services/jwt_token_service.py
from datetime import datetime, timedelta, timezone
from typing import Any, Dict, Optional

from jose import ExpiredSignatureError, JWTError, jwt

from app.ports.token_service import ITokenService
from app.settings.config import settings


class JWTTokenService(ITokenService):
    def __init__(
        self,
        secret: Optional[str] = None,
        algorithm: Optional[str] = None,
        access_ttl_minutes: Optional[int] = None,
        issuer: Optional[str] = None,
    ) -> None:
        self.secret = secret or settings.jwt_secret_key
        self.alg = algorithm or settings.jwt_algorithm  # e.g. "HS256"
        self.ttl_minutes = access_ttl_minutes or settings.access_token_expire_minutes
        self.issuer = issuer or getattr(settings, "jwt_issuer", None)  # optional

    def create_access_token(self, user_id: str) -> str:
        """
        user_id should be a string (UUID string is fine). Store in 'sub'.
        """
        now = datetime.now(timezone.utc)
        exp = now + timedelta(minutes=self.ttl_minutes)
        payload: Dict[str, Any] = {
            "sub": str(user_id),
            "iat": int(now.timestamp()),
            "exp": int(exp.timestamp()),
        }
        if self.issuer:
            payload["iss"] = self.issuer
        return jwt.encode(payload, self.secret, algorithm=self.alg)

    def decode_access_token(self, token: str) -> Dict[str, Any]:
        """
        Return the decoded claims dict. Raise ValueError on invalid/expired tokens.
        """
        try:
            return jwt.decode(
                token,
                self.secret,
                algorithms=[self.alg],
                options={"verify_aud": False},
                issuer=self.issuer if self.issuer else None,
            )
        except ExpiredSignatureError as e:
            raise ValueError("token_expired") from e
        except JWTError as e:
            raise ValueError("invalid_token") from e
