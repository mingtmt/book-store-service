from passlib.context import CryptContext

from app.ports.password_service import IPasswordService

_ctx = CryptContext(schemes=["bcrypt"], deprecated="auto")


class PasswordService(IPasswordService):
    def hash_password(self, raw: str) -> str:
        return _ctx.hash(raw)

    def verify_password(self, raw: str, hashed: str) -> bool:
        return _ctx.verify(raw, hashed)
