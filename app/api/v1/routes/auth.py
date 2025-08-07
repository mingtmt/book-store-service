from fastapi import APIRouter, Depends, HTTPException
from app.schemas.auth import LoginRequest, LoginResponse
from app.use_cases.auth.login_user import LoginUserUseCase
from app.infrastructure.db.sqlalchemy.user_impl import SqlAlchemyUserRepository
from sqlalchemy.orm import Session
from app.api.dependencies.db import get_db


router = APIRouter()

@router.post("/login", response_model=LoginResponse)
def login(data: LoginRequest, db: Session = Depends(get_db)):
    user_repo = SqlAlchemyUserRepository(db)
    use_case = LoginUserUseCase(user_repo)
    try:
        token = use_case.execute(data.email, data.password)
        return LoginResponse(access_token=token)
    except ValueError:
        raise HTTPException(status_code=401, detail="Invalid credentials")
