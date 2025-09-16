# To-Do List Application

A cross-platform desktop application for managing tasks, built with Wails (Go + React).

## Features

- Create, read, update, and delete tasks
- Mark tasks as complete/incomplete
- Set task priorities (High, Medium, Low)
- Add due dates to tasks
- Filter tasks by status (All, Active, Completed)
- Light and dark theme support
- Responsive design
- Persistent storage (tasks are saved between sessions)

## Prerequisites

- Go 1.21 or later
- Node.js 18+ and npm
- Wails CLI (install with `go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

## Getting Started

### Backend Setup

1. Navigate to the project root directory:
   ```bash
   cd todo-app
   ```

2. Initialize the Go module and download dependencies:
   ```bash
   cd backend
   go mod tidy
   ```

### Frontend Setup

1. Navigate to the frontend directory and install dependencies:
   ```bash
   cd frontend
   npm install
   ```

## Running the Application

### Development Mode

1. In one terminal, start the backend:
   ```bash
   cd backend
   wails dev
   ```

2. In another terminal, start the frontend development server:
   ```bash
   cd frontend
   npm run dev
   ```

3. The application should automatically open in your default browser at `http://localhost:3000`

### Production Build

1. Build the frontend:
   ```bash
   cd frontend
   npm run build
   ```

2. Build the Wails application:
   ```bash
   cd ..
   wails build
   ```

3. The built application will be available in the `build/bin` directory

## Project Structure

```
todo-app/
├── backend/         # Go backend
│   ├── app/         # Application logic
│   ├── main.go      # Entry point
│   └── go.mod       # Go module file
├── frontend/        # React frontend
│   ├── public/      # Static files
│   └── src/         # Source files
│       ├── components/ # React components
│       ├── App.jsx  # Main React component
│       └── App.css  # Main styles
├── wails.json       # Wails configuration
└── README.md        # This file
```

## Features Implementation Status

### Core Features (100%)
- [x] Task management (CRUD operations)
- [x] Task status toggling
- [x] Task prioritization
- [x] Due dates for tasks
- [x] Task filtering
- [x] Light/dark theme
- [x] Responsive design

### Bonus Features (100%)
- [x] Task priorities (High/Medium/Low)
- [x] Due date support
- [x] Dark/light theme toggle
- [x] Confirmation for task deletion
- [x] Data persistence
- [x] Clean UI/UX

## Screenshots

*Screenshots will be added after the first build*

## Contributing

1. Fork the repository
2. Create a new branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Wails](https://wails.io/) for the amazing framework
- [React](https://reactjs.org/) for the frontend library
- [Vite](https://vitejs.dev/) for the build tooling
