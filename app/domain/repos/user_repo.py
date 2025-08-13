import uuid
from abc import ABC, abstractmethod
from app.domain.models.user import User

class UserRepository(ABC):
    @abstractmethod
    def get_by_email(self, email: str) -> User | None:
        pass

    @abstractmethod
    def create(self, user: User) -> User:
        pass

    @abstractmethod
    def update(self, user: User) -> User:
        pass

    @abstractmethod
    def delete(self, user_id: uuid.UUID) -> bool:
        pass
    