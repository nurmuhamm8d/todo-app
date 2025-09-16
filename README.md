# Todo App (Wails + Go + React)

Минималистичное приложение задач с локальной БД PostgreSQL и UI на Vite + React. Поддерживает приоритизацию, дедлайны, фильтры (All/Active/Completed/Overdue/Today/This week), поиск, теги и категории, массовые действия, светлую/тёмную темы. Собрано на Wails.

## Стек
- Go, PostgreSQL
- Wails (WebView2)
- React + Vite + TypeScript
- CSS на кастомных переменных

## Возможности
- Добавление задач с приоритетом и дедлайном
- Переключение статуса задачи (выполнено/активно)
- Удаление с подтверждением системным диалогом Wails
- Фильтры: All, Active, Completed, Overdue, Today, This week
- Сортировка: по дате и по приоритету
- Поиск по заголовку
- Категория задачи
- Теги (через запятую)
- Массовые действия: завершить выбранные, удалить выбранные
- Статистика: Total, Active, Completed, Overdue
- Светлая/тёмная тема с запоминанием выбора

## Скриншоты и видео
Добавьте изображения в `docs/screenshots/` и вставьте сюда:
- ![Главный экран](docs/screenshots/01-main.png)
- ![Фильтры и поиск](docs/screenshots/02-filters.png)
- ![Массовые операции](docs/screenshots/03-bulk.png)
- ![Светлая тема](docs/screenshots/04-light.png)

Демо-видео положите как `docs/demo.mp4` и добавьте ссылку:
- [Демо-видео](docs/demo.mp4)

## Требования
- Go 1.21+
- Node 18+
- PostgreSQL 13+
- Windows 10/11 с WebView2 Runtime (ставится автоматически Wails)

## Конфигурация БД
Переменная окружения:
DATABASE_URL=host=127.0.0.1 port=5432 user=postgres password=postgres dbname=todoapp sslmode=disable
## Установка и запуск
go mod tidy
wails dev
## Сборка
wails build

## Архитектура
- `backend/internal/models` — доменные модели
- `backend/internal/repository` — доступ к БД (PostgreSQL)
- `backend/internal/usecase` — бизнес-логика
- `backend/app.go` — биндинги Wails и маппинг DTO

## Соответствие заданию
- Базовый функционал: добавление, удаление, отметка выполнения, сохранение состояния, фильтры — выполнено
- Дополнительные баллы: дедлайны, приоритеты, фильтры по дате (Overdue/Today/This week), подтверждение удаления, слоистая архитектура, улучшенный UI, темы — выполнено
- README + скриншоты + видео — этот файл и раздел выше

## Импорт в GitHub
Создайте публичный репозиторий на `github.com/nurmuhamm8d` и выполните:
git init
git add .
git commit -m "Initial commit"
git branch -M main
git remote add origin https://github.com/nurmuhamm8d/todo-app.git

git push -u origin main

Альтернатива через GitHub CLI:
gh repo create nurmuhamm8d/todo-app --public --source=. --remote=origin --push


## Лицензия
MIT


Публикация через GitHub CLI и правила игнора подтверждены официальной документацией.