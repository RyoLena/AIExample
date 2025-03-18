from pydantic import BaseModel
from typing import List,Optional


class ImageData(BaseModel):
    mime_type : str # 例如 "image/jpeg"
    data:str    # Base64 编码的图片数据

class ChatRequest(BaseModel):
    message:str
    image:Optional[ImageData] = None #图片列表
    conversation_id:Optional[str]=None

class ChatResponse(BaseModel):
    reply:str
    image_url: Optional[str] = None  # 新增字段，图片 URL
    conversation_id:Optional[str]=None