from app.domain.entities.role import Role
from app.domain.repositories.role_repo import IRoleRepository


class CreateRoleUseCase:
    def __init__(self, repo: IRoleRepository):
        self.repo = repo

    def execute(self, name: str):
        role = Role(name=name)
        return self.repo.create(role)
