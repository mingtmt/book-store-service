from app.domain.repos.user_repo import IUserRepository
from passlib.context import CryptContext
from app.domain.services.token_service import ITokenService
from app.core.exceptions import UnauthorizedException

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")

class LoginUserUseCase:
    def __init__(self, user_repo: IUserRepository, token_service: ITokenService):
        self.user_repo = user_repo
        self.token_service = token_service

    def execute(self, email: str, password: str) -> str:
        user = self.user_repo.get_by_email(email)
        if not user or not pwd_context.verify(password, user.hashed_password):
            raise UnauthorizedException("Invalid credentials")
        return self.token_service.create_access_token(str(user.id))
