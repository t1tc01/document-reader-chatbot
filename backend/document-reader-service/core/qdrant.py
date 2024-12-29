from qdrant_client import QdrantClient
from config.app_setting import settings

from config.app_setting import settings
q_client = QdrantClient(settings.vectordb_qdrant_url) 