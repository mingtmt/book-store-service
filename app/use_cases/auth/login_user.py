from app.domain.repos.user_repo import UserRepository
from passlib.context import CryptContext
from app.infrastructure.services.jwt_service import create_access_token
from app.core.exceptions import UnauthorizedException

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")

class LoginUserUseCase:
    def __init__(self, user_repo: UserRepository):
        self.user_repo = user_repo

    def execute(self, email: str, password: str) -> str:
        user = self.user_repo.get_by_email(email)
        if not user or not pwd_context.verify(password, user.hashed_password):
            raise UnauthorizedException("Invalid credentials")
        return create_access_token({"user_id": str(user.id)})
