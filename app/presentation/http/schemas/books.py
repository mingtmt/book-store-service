import uuid
from decimal import Decimal
from pydantic import BaseModel, Field, ConfigDict
from typing import Optional

class CreateBookRequest(BaseModel):
    title: str
    author: str
    price: Decimal
    description: str | None = None
    category: str

class UpdateBookRequest(BaseModel):
    title: Optional[str] = Field(default=None, min_length=1)
    author: Optional[str] = Field(default=None, min_length=1)
    price: Optional[Decimal] = Field(default=None, gt=0)
    description: Optional[str] = None
    category: Optional[str] = Field(default=None, min_length=1)

class BookOut(BaseModel):
    model_config = ConfigDict(from_attributes=True)
    id: uuid.UUID
    title: str
    author: str
    price: Decimal
    description: Optional[str]
    category: str
