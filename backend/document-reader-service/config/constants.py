import json

from dotenv import load_dotenv
from config.app_setting import settings
from langchain.prompts import PromptTemplate
load_dotenv()

qa_prompt = {}
raw_extract_prompt = {}
structure_extract_prompt = {}

f = open(settings.qa_promt)
qa_prompt_json = json.load(f)
for language, promt_str in qa_prompt_json.items():
    qa_prompt[language] = PromptTemplate.from_template(template=promt_str)
f.close()

f = open(settings.extract_promt)
extract_prompts = json.load(f)
for language, promt_dict in extract_prompts.items():
    raw_extract_prompt[language] = PromptTemplate.from_template(template=promt_dict["raw_extract"])
    structure_extract_prompt[language] = PromptTemplate.from_template(template=promt_dict["structure_extract"])
f.close()