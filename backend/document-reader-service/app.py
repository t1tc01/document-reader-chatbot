from loguru import logger
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from internal.document.router.v1 import document_router as document_v1 

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
    expose_headers=["*"],
)

# Example router
@app.get("/")
async def root():
    return {"message": "Hello, the APIs are now ready for your embeds and queries!"}

app.include_router(document_v1.router, prefix="/document/v1", tags=["document"])