from fastapi import APIRouter, Depends, Response, status
from sqlalchemy.orm import Session

from app.infrastructure.db.sqlalchemy.repositories.user_impl import (
    SqlAlchemyUserRepository,
)
from app.infrastructure.services.jwt_token_service import JWTTokenService
from app.infrastructure.services.password_service import PasswordService
from app.presentation.http.dependencies.db import get_db
from app.presentation.http.schemas.auth import (
    LoginRequest,
    LoginResponse,
    RegisterIn,
    RegisterResponse,
    UserOut,
)
from app.presentation.http.schemas.base import Envelope
from app.usecases.auth.login_user import LoginUserUseCase
from app.usecases.auth.register_user import RegisterUserCommand, RegisterUserUseCase

router = APIRouter()


@router.post(
    "/register",
    response_model=Envelope[RegisterResponse],
    status_code=status.HTTP_201_CREATED,
)
def register(payload: RegisterIn, response: Response, db: Session = Depends(get_db)):
    token_service = JWTTokenService()
    password_service = PasswordService()
    repo = SqlAlchemyUserRepository(db)
    cmd = RegisterUserCommand(
        email=payload.email,
        name=payload.name,
        age=payload.age,
        password=payload.password,
    )
    uc = RegisterUserUseCase(repo, token_service, password_service, issue_token=True)
    created_user, token = uc.execute(cmd)
    response.headers["Location"] = f"/api/v1/users/{created_user.id}"
    return Envelope(
        data=RegisterResponse(
            user=UserOut.model_validate(created_user), access_token=token
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
