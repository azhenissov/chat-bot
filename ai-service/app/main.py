from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from dotenv import load_dotenv

load_dotenv() 

from app.services.llm import generate_response

app = FastAPI(title="Gemini AI Service")

class Message(BaseModel):
    role:str
    content:str

class ChatRequest(BaseModel):
    chat_id:str | None = None 
    new_message: str
    context:list[Message] = []

@app.post("/generate")
async def generate_chat_response(request:ChatRequest):
    try:
        answer = await generate_response(
            new_message=request.new_message,
            context=[msg.model_dump() for msg in request.context]
        )
        return {"response": answer}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e)) 



@app.get("/health")
async def health_check():
    return {"status": "ok"}

