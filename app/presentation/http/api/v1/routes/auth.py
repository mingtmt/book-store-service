from fastapi import APIRouter, Depends, Response, status
from sqlalchemy.orm import Session

from app.infrastructure.db.sqlalchemy.repositories.user_impl import (
    SqlAlchemyUserRepository,
)
from app.infrastructure.services.jwt_token_service import JWTTokenService
from app.infrastructure.services.password_service import PasswordService
from app.presentation.http.dependencies.auth import get_current_user
from app.presentation.http.dependencies.db import get_db
from app.presentation.http.schemas.auth import (
    LoginIn,
    LoginOut,
    RegisterIn,
    RegisterOut,
    UpdateMeIn,
    UserOut,
)
from app.presentation.http.schemas.base import Envelope
from app.usecases.auth.login_user import LoginUserUseCase
from app.usecases.auth.register_user import RegisterUserCommand, RegisterUserUseCase
from app.usecases.auth.update_profile import UpdateProfileCmd, UpdateProfileUseCase

router = APIRouter()


@router.post(
    "/register",
    response_model=Envelope[RegisterOut],
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
        data=RegisterOut(user=UserOut.model_validate(created_user), access_token=token)
    )


@router.post(
    "/login",
    response_model=Envelope[LoginOut],
    status_code=status.HTTP_200_OK,
)
def login(payload: LoginIn, db: Session = Depends(get_db)):
    token_service = JWTTokenService()
    password_service = PasswordService()
    repo = SqlAlchemyUserRepository(db)
    uc = LoginUserUseCase(repo, token_service, password_service)
    token = uc.execute(payload.email, payload.password)
    return Envelope(data=LoginOut(access_token=token))


@router.get(
    "/me",
    response_model=Envelope[UserOut],
    status_code=status.HTTP_200_OK,
)
def me(current_user=Depends(get_current_user)):
    return Envelope(data=UserOut.from_domain(current_user))


@router.patch("/me", response_model=Envelope[UserOut], status_code=status.HTTP_200_OK)
def update_me(
    payload: UpdateMeIn,
    db: Session = Depends(get_db),
    current_user=Depends(get_current_user),
):
    repo = SqlAlchemyUserRepository(db)
    uc = UpdateProfileUseCase(repo)
    cmd = UpdateProfileCmd(
        **payload.model_dump(exclude_unset=True), user_id=current_user.id
    )
    updated = uc.execute(cmd)
    db.commit()
    return Envelope(data=UserOut.from_domain(updated))
