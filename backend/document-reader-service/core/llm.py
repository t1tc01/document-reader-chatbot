import os
from loguru import logger
from langchain_google_genai import ChatGoogleGenerativeAI
# from langchain.embeddings import HuggingFaceEmbeddings
from langchain_google_genai import GoogleGenerativeAIEmbeddings
from config.app_setting import settings

os.environ["GOOGLE_API_KEY"] = settings.google_api_key

def create_llm_component():
    logger.info("Using local LLM")
    model = ChatGoogleGenerativeAI(
        model="gemini-1.5-flash",
        temperature=0,
        max_tokens=None,
        timeout=None,
        max_retries=2,
    )
    embeddings = GoogleGenerativeAIEmbeddings(
        model="models/embedding-001",
    )
    logger.info("Created LLM Successfully")
    return model, embeddings

llm_model, llm_embeddings = create_llm_component()