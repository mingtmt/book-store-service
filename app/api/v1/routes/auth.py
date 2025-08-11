from fastapi import APIRouter, Depends, HTTPException, status
from app.schemas.auth import LoginRequest, LoginResponse, RegisterRequest, RegisterResponse, UserOut
from app.use_cases.auth.login_user import LoginUserUseCase
from app.use_cases.auth.register_user import RegisterUserUseCase
from app.infrastructure.db.sqlalchemy.user_impl import SqlAlchemyUserRepository
from sqlalchemy.orm import Session
from app.api.dependencies.db import get_db


router = APIRouter()

@router.post("/register", response_model=RegisterResponse, status_code=status.HTTP_201_CREATED)
def register(data: RegisterRequest, db: Session = Depends(get_db)):
    repo = SqlAlchemyUserRepository(db)
    uc = RegisterUserUseCase(repo, issue_token=True)
    user, token = uc.execute(data.email, data.password)
    return RegisterResponse(user=UserOut(id=user.id, email=user.email), access_token=token)

@router.post("/login", response_model=LoginResponse)
def login(data: LoginRequest, db: Session = Depends(get_db)):
    repo = SqlAlchemyUserRepository(db)
    uc = LoginUserUseCase(repo)
    token = uc.execute(data.email, data.password)
    return LoginResponse(access_token=token)
