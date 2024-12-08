from typing import Any,Dict
from pydantic_settings import BaseSettings
from config.enum import Enum

class AppEnvTypes(Enum):
    prod: str = "prod"
    dev: str = "dev"
    test: str = "test"

class BaseAppSettings(BaseSettings):
     app_env: AppEnvTypes = AppEnvTypes.dev

     class Config:
        env_file = ".env"

class AppSettings(BaseAppSettings):
    debug: bool = False

    vectordb_qdrant_url: str

    qa_promt: str = "prompt_template/qa_prompt.json"
    extract_promt: str = "prompt_template/extract_prompt.json"

    model_path: str
    embeddings_model_name: str
    embeddings_cache_folder: str

    model_n_ctx: int
    model_n_batch: int

    target_source_chunks: int
    chunk_size: int
    temperature: float

    # app
    worker_nb: int = 3

    class Config:
        validate_assignment = True

    @property
    def fastapi_kwargs(self) -> Dict[str, Any]:
        return {
            "debug": self.debug,
            "docs_url": self.docs_url,
            "title": self.title,
            "version": self.version,
        }
    
settings = AppSettings()

