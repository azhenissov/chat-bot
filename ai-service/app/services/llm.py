import os
import google.generativeai as genai 

genai.configure(api_key=os.environ.get("GEMINI_API_KEY"))
# system prompt
student_persona = """
Ты — обычный студент. Общайся просто, по-человечески, на "ты", без излишней вежливости и роботоподобных фраз (никаких "Чем я могу вам помочь сегодня?"). 
Используй умеренный студенческий и айтишный сленг. Ты увлекаешься бэкендом, пишешь на Go, разбираешься в микросервисах и докере.
Периодически ты немного жалуешься на недосып, горящие дедлайны по лабам и сложную подготовку к экзаменам по вычислительной математике (всякие формулы численного интегрирования и интерполяции).
Если ты чего-то не знаешь, так и скажи: "Блин, без понятия" или "Я такое еще не проходил", не пытайся выдумывать.
Твои ответы должны быть короткими и по делу, как в мессенджере.
"""

model = genai.GenerativeModel(
    'gemini-2.5-flash',
    system_instruction=student_persona
    )


async def generate_response(new_message:str, context: list[dict]) -> str:
    """
    context is in format [{"role": "user", "content": "Hi"}, {"role": "model", "content": "Hello"}] 
    """

    gemini_history = []
    for msg in context:
        role = "user" if msg["role"] == "user" else "model"
        gemini_history.append({
            "role" : role, 
            "parts": [msg["content"]]
        })
    
    chat = model.start_chat(history=gemini_history)

    response = await chat.send_message_async(new_message)

    return response.text
