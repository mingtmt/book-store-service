from pydantic import BaseModel
from typing import Generic, TypeVar, Optional

T = TypeVar("T")

class Envelope(BaseModel, Generic[T]):
    data: T
    meta: Optional[dict] = None

class PageMeta(BaseModel):
    page: int
    size: int
    total: int
