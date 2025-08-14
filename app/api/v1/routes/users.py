import uuid
from fastapi import APIRouter, Depends, status
from app.schemas.base import Envelope
from app.schemas.users import UpdateRequest, UpdateResponse
from app.use_cases.users.delete_user import DeleteUserUseCase
from app.use_cases.users.update_user import UpdateUserUseCase
from app.infrastructure.db.sqlalchemy.repos.user_impl import SqlAlchemyUserRepository
from sqlalchemy.orm import Session
from app.infrastructure.web.dependencies.db import get_db

router = APIRouter()

@router.put("/", response_model=Envelope[UpdateResponse], status_code=status.HTTP_200_OK)
def update(data: UpdateRequest, db: Session = Depends(get_db)):
    repo = SqlAlchemyUserRepository(db)
    uc = UpdateUserUseCase(repo)
    user = uc.execute(data.id, data.email, data.password)
    return Envelope(data=UpdateResponse(id=user.id, email=user.email))


@router.delete("/{user_id}", status_code=status.HTTP_204_NO_CONTENT)
def delete(user_id: uuid.UUID, db: Session = Depends(get_db)):
    repo = SqlAlchemyUserRepository(db)
    uc = DeleteUserUseCase(repo)
    uc.execute(user_id)
    return
