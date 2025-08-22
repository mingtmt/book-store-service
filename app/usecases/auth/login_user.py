from app.domain.repositories.user_repo import IUserRepository
from app.ports.password_service import IPasswordService
from app.ports.token_service import ITokenService
from app.usecases.errors import Unauthorized
from app.utils.helper import normalize_email


class LoginUserUseCase:
    def __init__(
        self,
        user_repo: IUserRepository,
        token_service: ITokenService,
        password_service: IPasswordService,
    ):
        self.user_repo = user_repo
        self.token_service = token_service
        self.password_service = password_service

    def execute(self, email: str, password: str) -> str:
        normalized_email = normalize_email(email)
        user = self.user_repo.get_by_email(normalized_email)
        if not user or not self.password_service.verify_password(
            password, user.hashed_password
        ):
            raise Unauthorized("Invalid credentials")
        return self.token_service.create_access_token(str(user.id))
