import uuid
from pydantic import BaseModel, EmailStr, Field

class RegisterRequest(BaseModel):
    email: EmailStr
    password: str = Field(min_length=8, max_length=128)

class UserOut(BaseModel):
    id: uuid.UUID
    email: EmailStr

class RegisterResponse(BaseModel):
    user: UserOut
    access_token: str | None = None
    token_type: str | None = "bearer"

class LoginRequest(BaseModel):
    email: str
    password: str
    
class LoginResponse(BaseModel):
    access_token: str
    token_type: str = "bearer"
