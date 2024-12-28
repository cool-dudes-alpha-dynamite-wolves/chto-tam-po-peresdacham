from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.chrome.options import Options
from webdriver_manager.chrome import ChromeDriverManager
import os
import time

# Настройка папки для загрузки файлов
download_folder = os.path.abspath("downloaded_files")
os.makedirs(download_folder, exist_ok=True)

# Настройка опций для Chrome
chrome_options = Options()
chrome_options.add_argument('--headless')  # Фоновый режим (можно убрать для отладки)
chrome_options.add_argument('--disable-gpu')
chrome_options.add_argument('--no-sandbox')
chrome_options.add_argument('--disable-dev-shm-usage')
chrome_options.add_experimental_option("prefs", {
    "download.default_directory": download_folder,  # Папка для загрузки
    "download.prompt_for_download": False,          # Отключение запросов
    "download.directory_upgrade": True,             # Автоматическое обновление папки
    "safebrowsing.enabled": True                    # Включение безопасной загрузки
})

# Инициализация драйвера
driver = webdriver.Chrome(service=Service(ChromeDriverManager().install()), options=chrome_options)

# Открытие целевой страницы
base_url = "https://misis.ru"
page_url = f"{base_url}/students/likvidacia/#tab-1-3"
driver.get(page_url)
time.sleep(5)  # Ожидание загрузки страницы

# Поиск всех ссылок на файлы
links = driver.find_elements(By.TAG_NAME, 'a')
for link in links:
    href = link.get_attribute('href')  # Получение атрибута href
    if href and ('.xls' in href or '.xlsx' in href):
        if href.startswith("/"):  # Если ссылка относительная, добавляем домен
            href = base_url + href
        try:
            print(f"Кликаем по ссылке: {href}")
            link.click()  # Симуляция клика
            time.sleep(2)  # Пауза, чтобы файл успел начать загружаться
        except Exception as e:
            print(f"Ошибка при клике по ссылке {href}: {e}")

# Закрытие браузера
driver.quit()

print(f"Все файлы сохранены в папке: {download_folder}")
