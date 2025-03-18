from abc import ABC, abstractmethod


class AIModel(ABC):
    @abstractmethod
    def generate_text(self,message:str,conversation_id:str | None=None) -> str:
        pass
    @abstractmethod
    def generate_image(self, prompt: str) -> str:  # 返回图片 URL
        pass

import google.generativeai as genai
from .. import  config

class GeminiModel(AIModel):
    def __init__(self):
        genai.configure(api_key=config.setting.gemin_api_key)
        self.model = genai.GenerativeModel(config.setting.model_name)

    def generate_text(self, message: str, conversation_id: str | None = None) -> str:
        response = self.model.generate_content(message)
        return response.text

    def generate_image(self, prompt: str) -> str:
        # TODO: 实现图片生成逻辑，调用 Gemini 的图片生成 API
        # 目前返回空字符串
        return ""

from .models import ChatResponse

class AIPipeline:
    def __init__(self, ai_model: AIModel):
        self.ai_model = ai_model

    def get_response(self, message: str, conversation_id: str | None = None) -> ChatResponse:
        reply = self.ai_model.generate_text(message, conversation_id)
        image_url = ""  # 默认为空
        # TODO: 判断是否需要生成图片，如果需要，则调用 self.ai_model.generate_image()
        return ChatResponse(reply=reply, image_url=image_url, conversation_id=conversation_id)

gemini_model = GeminiModel()
ai_pipe = AIPipeline(gemini_model)

