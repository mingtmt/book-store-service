import uuid
from decimal import Decimal
from pydantic import BaseModel

class CreateBookRequest(BaseModel):
    title: str
    author: str
    price: Decimal
    description: str | None = None
    category: str

class CreateBookResponse(BaseModel):
    id: uuid.UUID
    title: str
    author: str
    price: Decimal
    description: str | None = None
    category: str

class GetBookResponse(BaseModel):
    id: uuid.UUID
    title: str
    author: str
    price: Decimal
    description: str | None = None
    category: str