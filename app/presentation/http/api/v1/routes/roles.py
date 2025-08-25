from fastapi import APIRouter, Depends, Response, status
from sqlalchemy.orm import Session

from app.infrastructure.db.sqlalchemy.repositories.role_impl import (
    SqlAlchemyRoleRepository,
)
from app.presentation.http.dependencies.db import get_db
from app.presentation.http.schemas.base import Envelope
from app.presentation.http.schemas.roles import CreateRoleIn, RoleOut
from app.usecases.roles.create_role import CreateRoleUseCase

router = APIRouter()


@router.post(
    "/",
    response_model=Envelope[RoleOut],
    status_code=status.HTTP_201_CREATED,
)
def create(payload: CreateRoleIn, response: Response, db: Session = Depends(get_db)):
    repo = SqlAlchemyRoleRepository(db)
    uc = CreateRoleUseCase(repo)
    created = uc.execute(payload.name)
    response.headers["Location"] = f"/api/v1/roles/{created.id}"
    return Envelope(data=RoleOut.model_validate(created))
