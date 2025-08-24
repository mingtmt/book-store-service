import uuid
from typing import Annotated, Optional

from pydantic import (
    BaseModel,
    ConfigDict,
    EmailStr,
    Field,
    field_validator,
    model_validator,
)

Name = Annotated[str, Field(min_length=1, max_length=100)]
Age = Annotated[int, Field(ge=0, le=150)]
Password = Annotated[str, Field(min_length=8, max_length=128)]


class UserOut(BaseModel):
    model_config = ConfigDict(from_attributes=True, extra="ignore")

    id: uuid.UUID
    email: EmailStr
    name: str
    age: int

    @classmethod
    def from_domain(cls, user) -> "UserOut":
        return cls.model_validate(user)


class RegisterIn(BaseModel):
    model_config = ConfigDict(
        extra="forbid",
        json_schema_extra={
            "example": {
                "email": "ada.lovelace@example.com",
                "name": "Ada Lovelace",
                "age": 28,
                "password": "S3cureP@ssw0rd",
            }
        },
    )

    email: EmailStr
    name: Name
    age: Age
    password: Password

    @field_validator("email", mode="before")
    @classmethod
    def _normalize_email(cls, v):
        if isinstance(v, str):
            v = v.strip()
            v = " ".join(v.split())
            v = v.lower()
        return v

    @field_validator("name", mode="before")
    @classmethod
    def _normalize_name(cls, v):
        if isinstance(v, str):
            v = v.strip()
            v = " ".join(v.split())
        return v


class RegisterOut(BaseModel):
    user: UserOut
    access_token: str | None = None
    token_type: str | None = "bearer"


class LoginIn(BaseModel):
    model_config = ConfigDict(
        extra="forbid",
        json_schema_extra={
            "example": {
                "email": "ada.lovelace@example.com",
                "password": "S3cureP@ssw0rd",
            }
        },
    )

    email: EmailStr
    password: Password

    @field_validator("email", mode="before")
    @classmethod
    def _normalize_email(cls, v):
        if isinstance(v, str):
            v = v.strip()
            v = " ".join(v.split())
            v = v.lower()
        return v


class LoginOut(BaseModel):
    access_token: str
    token_type: str = "bearer"


class UpdateMeIn(BaseModel):
    model_config = ConfigDict(
        extra="forbid",
        json_schema_extra={"example": {"name": "Ada L.", "age": 29}},
    )
    name: Optional[Name] = None
    age: Optional[Age] = None

    @field_validator("name", mode="before")
    @classmethod
    def _normalize_name(cls, v):
        if isinstance(v, str):
            v = " ".join(v.strip().split())
        return v


class ChangePasswordIn(BaseModel):
    model_config = ConfigDict(
        extra="forbid",
        json_schema_extra={
            "example": {
                "old_password": "OldP@ssw0rd",
                "new_password": "N3wS3cureP@ssw0rd",
            }
        },
    )

    old_password: Password
    new_password: Password

    @model_validator(mode="after")
    def check_different_passwords(self):
        if self.old_password == self.new_password:
            raise ValueError("New password must be different from old password")
        return self
