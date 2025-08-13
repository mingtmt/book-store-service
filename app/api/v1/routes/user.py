import uuid
from fastapi import APIRouter, Depends, status
from app.use_cases.user.delete_user import DeleteUserUseCase
from app.infrastructure.db.sqlalchemy.user_impl import SqlAlchemyUserRepository
from sqlalchemy.orm import Session
from app.api.dependencies.db import get_db

router = APIRouter()

@router.delete("/{user_id}", status_code=status.HTTP_204_NO_CONTENT)
def delete(user_id: uuid.UUID, db: Session = Depends(get_db)):
    repo = SqlAlchemyUserRepository(db)
    uc = DeleteUserUseCase(repo)
    uc.execute(user_id)
    return
