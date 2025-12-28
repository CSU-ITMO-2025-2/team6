# import asyncio
# import json
# import os
# import tempfile
# from nats.aio.client import Client as NATS
# from nats.errors import TimeoutError
# import base64
# from ultralytics import YOLO
# import numpy as np
# import logging
#
# # Настройка логирования
# logging.basicConfig(level=logging.INFO)
# logger = logging.getLogger(__name__)
#
# class ImageAnalysisService:
#     def __init__(self, nats_url, model_path):
#         self.nats_url = nats_url
#         self.nc = None
#         self.model = YOLO(model_path)
#
#         # Названия очередей
#         self.request_subject = "core-task-ml.request"
#         self.response_subject = "core-task-ml.response"
#
#     async def connect(self):
#         """Подключение к NATS"""
#         self.nc = NATS()
#         await self.nc.connect(self.nats_url)
#         logger.info(f"Connected to NATS at {self.nats_url}")
#
#     async def process_image_message(self, msg):
#         """Обработка входящего сообщения с изображением"""
#         try:
#             # Парсим входящее сообщение
#             data = json.loads(msg.data.decode())
#             study_id = data.get("study_id")
#             image_data = data.get("image")  # base64 encoded image
#             reply_to = msg.reply
#
#             logger.info(f"Received image analysis request for study_id: {study_id}")
#
#             if not image_data:
#                 await self._send_error(reply_to, "No image data provided")
#                 return
#
#             # Декодируем изображение из base64
#             try:
#                 image_bytes = base64.b64decode(image_data)
#             except Exception as e:
#                 await self._send_error(reply_to, f"Invalid study data: {str(e)}")
#                 return
#
#             # Анализируем изображение
#             result = await self.analyze_image(image_bytes, study_id)
#
#             # Отправляем результат обратно
#             response = {
#                 "status": "success",
#                 "study_id": study_id,
#                 "predict": result
#             }
#
#             if reply_to:
#                 await self.nc.publish(reply_to, json.dumps(response).encode())
#                 logger.info(f"Sent analysis result for study_id: {study_id}")
#
#         except Exception as e:
#             logger.error(f"Error processing message: {str(e)}")
#             if msg.reply:
#                 await self._send_error(msg.reply, f"Processing error: {str(e)}")
#
#     async def analyze_image(self, image_bytes, image_guid):
#         """Анализ изображения с помощью YOLO модели"""
#         try:
#             # Сохраняем изображение во временный файл
#             with tempfile.NamedTemporaryFile(delete=False, suffix='.jpg') as temp_file:
#                 temp_file.write(image_bytes)
#                 temp_path = temp_file.name
#
#             # Анализируем изображение
#             results = self.model(source=temp_path)
#
#             # Обрабатываем результаты
#             if hasattr(results[0], 'probs') and results[0].probs is not None:
#                 # Классификация
#                 pred_scores = results[0].probs.data.numpy()
#                 max_index = np.argmax(pred_scores)
#                 max_score = np.max(pred_scores)
#                 predicted_class = results[0].names[max_index]
#
#                 result = {
#                     'type': predicted_class,
#                     'score': str(max_score)
#                 }
#             else:
#                 # Детекция объектов
#                 boxes = results[0].boxes
#                 if len(boxes) > 0:
#                     detections = []
#                     for box in boxes:
#                         cls = int(box.cls[0])
#                         confidence = float(box.conf[0])
#                         detections.append({
#                             'class': results[0].names[cls],
#                             'confidence': confidence,
#                             'bbox': box.xyxy[0].tolist()
#                         })
#                     result = {
#                         'detections': detections,
#                         'count': len(detections)
#                     }
#                 else:
#                     result = {
#                         'type': 'no_detections',
#                         'score': '0.0'
#                     }
#
#             # Удаляем временный файл
#             os.unlink(temp_path)
#
#             return result
#
#         except Exception as e:
#             logger.error(f"Error analyzing image {image_guid}: {str(e)}")
#             return {
#                 'type': 'error',
#                 'score': '0.0',
#                 'error': str(e)
#             }
#
#     async def _send_error(self, reply_to, error_message):
#         """Отправка сообщения об ошибке"""
#         error_response = {
#             "status": "error",
#             "error": error_message
#         }
#         await self.nc.publish(reply_to, json.dumps(error_response).encode())
#
#     async def subscribe(self):
#         """Подписка на очередь запросов"""
#         if not self.nc:
#             await self.connect()
#
#         await self.nc.subscribe(self.request_subject, cb=self.process_image_message)
#         logger.info(f"Subscribed to {self.request_subject}")
#
#     async def run(self):
#         """Запуск сервиса"""
#         await self.connect()
#         await self.subscribe()
#
#         logger.info("Image Analysis Service started. Waiting for messages...")
#
#         try:
#             # Бесконечный цикл для поддержания работы сервиса
#             while True:
#                 await asyncio.sleep(1)
#         except KeyboardInterrupt:
#             logger.info("Shutting down...")
#         finally:
#             await self.nc.close()
#
#
# async def main():
#     service = ImageAnalysisService(
#         nats_url="nats:4222",
#         model_path="/app/model/best.pt"
#     )
#     await service.run()
#
# if __name__ == "__main__":
#     asyncio.run(main())

import asyncio
import json
import os
import tempfile
from nats.aio.client import Client as NATS
import base64
from ultralytics import YOLO
import numpy as np
import logging
from dotenv import load_dotenv

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

def load_configuration():
    """Загружает конфигурацию из .env файла или переменных окружения"""
    load_dotenv()

    config = {
        'nats_url': os.getenv('NATS_URL', 'nats://nats-srv:4222'),
        'model_path': os.getenv('MODEL_PATH', '/app/model/best.pt')
    }

    print(f"Loaded configuration: {config}")
    return config

class ImageAnalysisService:
    def __init__(self, nats_url=None, model_path=None):
        config = load_configuration()

        # Используем параметры или конфигурацию
        self.nats_url = nats_url or config['nats_url']
        self.model_path = model_path or config['model_path']

        self.nc = None
        self.model = YOLO(self.model_path)

        # Стримы
        self.request_subject = "core-task-ml.request"
        self.response_subject = "core-task-ml.response"

        logger.info(f"Initialized with NATS: {self.nats_url}, Model: {self.model_path}")

    async def connect(self):
        """Подключение к NATS"""
        self.nc = NATS()
        await self.nc.connect(self.nats_url)
        logger.info(f"Connected to NATS at {self.nats_url}")

    async def subscribe(self):
        """Подписка на запросы"""
        if self.nc is None:
            await self.connect()

        await self.nc.subscribe(self.request_subject, cb=self.process_image_message)
        logger.info(f"Subscribed to {self.request_subject}")

    async def run(self):
        """Запуск сервиса"""
        await self.connect()
        await self.subscribe()
        logger.info("Image Analysis Service started. Waiting for messages...")

        try:
            while True:
                await asyncio.sleep(1)
        except KeyboardInterrupt:
            logger.info("Shutting down...")
        finally:
            if self.nc:
                await self.nc.close()

    # ============================================================
    #                    PROCESS MESSAGE
    # ============================================================
    async def process_image_message(self, msg):
        try:
            data = json.loads(msg.data.decode())
            study_id = data.get("study_id")
            image_data = data.get("image")

            logger.info(f"Received image analysis request for study_id: {study_id}")

            if not image_data:
                await self._send_error(study_id, "no image data provided")
                return

            try:
                image_bytes = base64.b64decode(image_data)
            except Exception as e:
                await self._send_error(study_id, f"invalid base64 image: {str(e)}")
                return

            result = await self.analyze_image(image_bytes)

            # Формируем Go-совместимый ответ
            response = {
                "study_id": study_id,
                "type": result["type"],
                "score": result["score"]
            }

            if result.get("error"):
                response["error"] = result["error"]

            await self.nc.publish(self.response_subject, json.dumps(response).encode())
            logger.info(f"Sent ML analysis result for study_id={study_id}")

        except Exception as e:
            logger.error(f"Error processing message: {str(e)}")
            await self._send_error(None, f"processing error: {str(e)}")

    # ============================================================
    #                    ANALYZE IMAGE
    # ============================================================
    async def analyze_image(self, image_bytes):
        """Анализ изображения YOLO (классификация / детекция)"""
        try:
            with tempfile.NamedTemporaryFile(delete=False, suffix=".jpg") as tmp:
                tmp.write(image_bytes)
                tmp_path = tmp.name

            results = self.model(source=tmp_path)

            # ----- Классификация -----
            if hasattr(results[0], "probs") and results[0].probs is not None:
                scores = results[0].probs.data.numpy()
                max_idx = int(np.argmax(scores))
                max_score = float(np.max(scores))

                predicted_class = results[0].names[max_idx]

                os.unlink(tmp_path)

                return {
                    "type": predicted_class,
                    "score": int(max_score * 100)   # int для Go
                }

            # ----- Детекция -----
            boxes = results[0].boxes
            if len(boxes) > 0:
                cls = int(boxes[0].cls[0])
                conf = float(boxes[0].conf[0])
                predicted_class = results[0].names[cls]

                os.unlink(tmp_path)

                return {
                    "type": predicted_class,
                    "score": int(conf * 100)
                }

            os.unlink(tmp_path)

            return {
                "type": "no_detections",
                "score": 0
            }

        except Exception as e:
            logger.error(f"Error analyzing image: {str(e)}")
            return {
                "type": "error",
                "score": 0,
                "error": str(e)
            }

    # ============================================================
    #                        ERROR RESPONSE
    # ============================================================
    async def _send_error(self, study_id, error_message):
        """Отправка ошибки"""
        response = {
            "study_id": study_id,
            "type": "error",
            "score": 0,
            "error": error_message
        }
        if self.nc:
            await self.nc.publish(self.response_subject, json.dumps(response).encode())
        logger.info(f"Sent error response: {error_message}")


# ============================================================
#                       ENTRYPOINT
# ============================================================
async def main():
    # Используем конфигурацию из окружения
    service = ImageAnalysisService()  # Без параметров - загрузит из конфигурации
    await service.run()

if __name__ == "__main__":
    asyncio.run(main())