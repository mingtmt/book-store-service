import uuid
from abc import ABC, abstractmethod

from app.domain.entities.role import Role


class IRoleRepository(ABC):
    @abstractmethod
    def get_by_id(self, role_id: uuid.UUID) -> Role: ...

    @abstractmethod
    def get_all(self) -> list[Role]: ...

    @abstractmethod
    def create(self, role: Role) -> Role: ...

    @abstractmethod
    def save(self, role: Role) -> Role: ...

    @abstractmethod
    def delete(self, role_id: uuid.UUID) -> None: ...
