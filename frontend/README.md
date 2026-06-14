# PriceTracker: Frontend SPA

Клиентская часть (SPA) для проекта **PriceTracker**, предназначенная для визуализации отслеживаемых товаров и изменения их цен, демонстрирует современный подход к созданию пользовательских интерфейсов с минималистичным дизайном.

## 🚀 О проекте

Приложение отвечает за взаимодействие с пользователем и отображение данных, полученных от API Gateway.
Основные возможности:
- Добавление ссылок на товары для парсинга.
- Просмотр списка всех отслеживаемых товаров с их текущими статусами.
- Визуализация истории изменения цен на графике.
- Плавные анимации и обработка состояний загрузки/ошибок с уведомлениями (fallback механизм на моковые данные при недоступности API).

## 🛠️ Стек технологий
<div> 
<img src="https://img.shields.io/badge/React-20232A?style=flat&logo=react&logoColor=61DAFB" alt="React"/> 
<img src="https://img.shields.io/badge/TypeScript-007ACC?style=flat&logo=typescript&logoColor=white" alt="TypeScript"/>
<img src="https://img.shields.io/badge/Vite-B73BFE?style=flat&logo=vite&logoColor=FFD62E" alt="Vite"/> 
<img src="https://img.shields.io/badge/Tailwind_CSS-38B2AC?style=flat&logo=tailwind-css&logoColor=white" alt="Tailwind CSS"/> 
<img src="https://img.shields.io/badge/Recharts-22B573?style=flat&logo=react&logoColor=white" alt="Recharts"/>
</div>

## ⚙️ Как запустить локально

Чтобы развернуть frontend-часть, выполните следующие шаги:

### Установка и запуск

1.  **Установите зависимости:**
    Убедитесь, что у вас установлен Node.js, затем выполните:
    ```sh
    npm install
    ```

2.  **Запустите сервер для разработки:**
    ```sh
    npm run dev
    ```
    *Приложение будет доступно по адресу `http://localhost:5173`.*

## 🌐 Сборка и деплой

Для сборки оптимизированной production-версии проекта выполните команду:
```sh
npm run build
```
Готовые статические файлы появятся в директории `dist` и будут готовы к деплою на любой хостинг (например, Vercel, Netlify или Nginx).

---

**Примечание:** Это учебный проект, в нём могут быть ошибки и упрощения.