from abc import ABC, abstractmethod
from typing import Optional


class ITokenService(ABC):
    @abstractmethod
    def create_access_token(self, user_id: int) -> str: ...

    @abstractmethod
    def decode_access_token(self, token: str) -> Optional[int]: ...
