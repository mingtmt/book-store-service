from abc import ABC, abstractmethod


class IPasswordService(ABC):
    @abstractmethod
    def hash_password(self, raw: str) -> str: ...

    @abstractmethod
    def verify_password(self, raw: str, hashed: str) -> bool: ...
