from urllib.request import localhost

from fastapi import FastAPI,HTTPException
from .models import ChatRequest,ChatResponse
from .ai_log import ai_pipe
from .. import config

app = FastAPI()

@app.post("/chat",response_model = ChatResponse)

async def chat_endpoint(request: ChatRequest):
    try:
        response = ai_pipe.get_response(request.message,request.conversation_id)
        return response
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/health")
async def health_check():
    return {"status": "OK"}

import sys
import os
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host=localhost(), port=config.setting.port)