import logging
from logging.handlers import RotatingFileHandler
from app.settings.config import settings

def setup_logging():
    handler = RotatingFileHandler(
        settings.log_file,
        maxBytes=settings.log_max_bytes,
        backupCount=settings.log_backup_count
    )
    formatter = logging.Formatter("%(asctime)s - %(levelname)s - %(message)s")
    handler.setFormatter(formatter)

    root = logging.getLogger()
    root.setLevel(settings.log_level)
    root.addHandler(handler)

setup_logging()
logger = logging.getLogger(__name__)