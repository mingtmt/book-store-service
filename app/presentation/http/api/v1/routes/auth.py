from fastapi import APIRouter, Depends, status
from sqlalchemy.orm import Session

from app.infrastructure.db.sqlalchemy.repositories.user_impl import (
    SqlAlchemyUserRepository,
)
from app.infrastructure.services.jwt_token_service import JWTTokenService
from app.presentation.http.dependencies.db import get_db
from app.presentation.http.schemas.auth import (
    LoginRequest,
    LoginResponse,
    RegisterRequest,
    RegisterResponse,
    UserOut,
)
from app.presentation.http.schemas.base import Envelope
from app.usecases.auth.login_user import LoginUserUseCase
from app.usecases.auth.register_user import RegisterUserUseCase

router = APIRouter()


@router.post(
    "/register",
    response_model=Envelope[RegisterResponse],
    status_code=status.HTTP_201_CREATED,
)
def register(data: RegisterRequest, db: Session = Depends(get_db)):
    token_service = JWTTokenService()
    repo = SqlAlchemyUserRepository(db)
    uc = RegisterUserUseCase(repo, token_service=token_service, issue_token=True)
    user, token = uc.execute(data.email, data.password)
    return Envelope(
        data=RegisterResponse(
            user=UserOut(id=user.id, email=user.email), access_token=token
        )
    )


@router.post(
    "/login", response_model=Envelope[LoginResponse], status_code=status.HTTP_200_OK
)
def login(data: LoginRequest, db: Session = Depends(get_db)):
    token_service = JWTTokenService()
    repo = SqlAlchemyUserRepository(db)
    uc = LoginUserUseCase(repo, token_service=token_service)
    token = uc.execute(data.email, data.password)
    return Envelope(data=LoginResponse(access_token=token))
