from dataclasses import dataclass

from app.domain.entities.user import User
from app.domain.errors import EmailAlreadyExists
from app.domain.repositories.user_repo import IUserRepository
from app.ports.password_service import IPasswordService
from app.ports.token_service import ITokenService
from app.utils.helper import normalize_email


@dataclass(frozen=True)
class RegisterUserCommand:
    email: str
    name: str
    age: int
    password: str


class RegisterUserUseCase:
    def __init__(
        self,
        user_repo: IUserRepository,
        token_service: ITokenService,
        password_service: IPasswordService,
        issue_token: bool = True,
    ):
        self.user_repo = user_repo
        self.token_service = token_service
        self.password_service = password_service
        self.issue_token = issue_token

    def execute(self, cmd: RegisterUserCommand) -> tuple[User, str | None]:
        normalized_email = normalize_email(cmd.email)
        if self.user_repo.get_by_email(normalized_email):
            raise EmailAlreadyExists(context={"email": normalized_email})

        hashed = self.password_service.hash_password(cmd.password)

        new_user = User(
            id=None,
            email=normalized_email,
            name=cmd.name,
            age=cmd.age,
            hashed_password=hashed,
        )
        created = self.user_repo.create(new_user)

        token = (
            self.token_service.create_access_token(str(created.id))
            if self.issue_token
            else None
        )
        return created, token
