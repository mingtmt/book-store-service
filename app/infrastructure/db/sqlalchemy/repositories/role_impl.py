import uuid

from sqlalchemy import select
from sqlalchemy.orm import Session

from app.domain.entities.role import Role
from app.domain.errors import RoleNotFound
from app.domain.repositories.role_repo import IRoleRepository
from app.infrastructure.db.sqlalchemy.mappers.orm_mapper import (
    apply_domain_to_orm,
    domain_to_orm,
    orm_to_domain,
)
from app.infrastructure.db.sqlalchemy.models.role_model import RoleModel


class SqlAlchemyRoleRepository(IRoleRepository):
    def __init__(self, db: Session):
        self.db = db

    def get_by_id(self, id: uuid.UUID) -> Role:
        m = (
            self.db.execute(select(RoleModel).where(RoleModel.id == id))
            .scalars()
            .first()
        )
        if not m:
            raise RoleNotFound
        return orm_to_domain(m, Role)

    def get_all(self) -> list[Role]:
        m = self.db.execute(select(RoleModel)).scalars().all()
        return [orm_to_domain(r, Role) for r in m]

    def create(self, role: Role) -> Role:
        m = domain_to_orm(role, RoleModel)
        self.db.add(m)
        try:
            self.db.commit()
            self.db.refresh(m)
        except Exception as e:
            self.db.rollback()
            raise e
        return orm_to_domain(m, Role)

    def save(self, role: Role) -> Role:
        return super().save(role)

    def delete(self, id: uuid.UUID) -> Role:
        return super().delete(id)
