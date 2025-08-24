import uuid
from dataclasses import dataclass
from typing import Optional

from app.domain.entities.user import User
from app.domain.errors import UserNotFound
from app.domain.repositories.user_repo import IUserRepository


@dataclass
class UpdateProfileCmd:
    user_id: uuid.UUID
    name: Optional[str] = None
    age: Optional[int] = None


class UpdateProfileUseCase:
    def __init__(self, repo: IUserRepository) -> None:
        self.repo = repo

    def execute(self, cmd: UpdateProfileCmd) -> User:
        user = self.repo.get_by_id(cmd.user_id)
        if not user:
            raise UserNotFound(context={"user_id": str(cmd.user_id)})

        new_name = cmd.name if cmd.name is not None else user.name
        new_age = cmd.age if cmd.age is not None else user.age

        updated = User(
            id=user.id,
            email=user.email,
            name=new_name,
            age=new_age,
            hashed_password=user.hashed_password,
        )
        return self.repo.save(updated)
