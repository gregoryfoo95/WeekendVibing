# 🦸‍♂️ FitHero - Gamified Fitness Application

A fun and engaging gamified fitness web application that transforms your daily exercise routine into an epic superhero adventure! Complete daily fitness challenges, earn points, level up your character, and unlock achievements.

## 🌟 Features

### 🎯 Daily Challenges
- **Random Fitness Tasks**: Get 3 personalized fitness tasks every day
- **Difficulty Levels**: Easy, Medium, and Hard tasks with varying point rewards
- **Categories**: Cardio, Strength, Flexibility, and Wellness exercises
- **Point System**: Earn 5-60 points per completed task

### 🦸‍♂️ Character Progression
- **Level System**: Progress from Level 1 to Level 5 based on points earned
- **Character Evolution**: Start as "Rookie Hero" and evolve to "Ultimate Hero"
- **Job Advancement**: Unlock professional titles from "Fitness Novice" to "Health Guru"

### 🏆 Achievement Store
- **Character Upgrades**: Unlock new superhero personas
- **Job Titles**: Advance your fitness career
- **Special Badges**: Earn unique badges for specific accomplishments
- **Point Economy**: Spend earned points to unlock achievements

### 📊 Social Features
- **Leaderboard**: Compete with other fitness heroes
- **Progress Tracking**: Visual progress bars and statistics
- **Community Rankings**: See how you stack up against others

### 🎨 Modern UI/UX
- **Material-UI Design**: Beautiful, responsive interface
- **Smooth Animations**: Framer Motion powered transitions
- **Cute & Relatable**: Designed for all ages, especially sedentary beginners
- **Mobile-First**: Responsive design that works on all devices

## 🏗️ Architecture

### Frontend (React + TypeScript)
- **Framework**: React 18 with TypeScript
- **UI Library**: Material-UI (MUI) v5
- **Animations**: Framer Motion
- **State Management**: React Query for server state
- **Routing**: React Router v6
- **HTTP Client**: Axios

### Backend (Go + Gin)
- **Framework**: Gin web framework
- **Language**: Go 1.21
- **Database**: PostgreSQL 15
- **CORS**: Enabled for frontend communication
- **API**: RESTful JSON API

### Database (PostgreSQL)
- **Users**: Store user profiles and progress
- **Tasks**: Predefined fitness challenges
- **Daily Tasks**: User-specific daily assignments
- **Achievements**: Unlockable rewards and upgrades
- **User Achievements**: Track unlocked achievements

### Infrastructure (Docker)
- **Containerization**: Each service runs in its own Docker container
- **Orchestration**: Docker Compose for multi-service setup
- **Networking**: Internal Docker network for service communication
- **Volumes**: Persistent PostgreSQL data storage

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose installed
- Git for cloning the repository

### Installation & Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd WeekendVibing
   ```

2. **Start the application**
   ```bash
   docker-compose up --build
   ```

3. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Database: localhost:5432

### First Time Setup

1. **Create your hero account** on the home page
2. **Generate daily tasks** on the dashboard
3. **Complete tasks** to earn points
4. **Visit the Achievement Store** to spend points
5. **Check the Leaderboard** to see your ranking

## 📁 Project Structure

```
WeekendVibing/
├── docker-compose.yml          # Multi-service orchestration
├── README.md                   # This file
├── backend/                    # Go backend service
│   ├── Dockerfile             # Backend container config
│   ├── go.mod                 # Go dependencies
│   ├── go.sum                 # Dependency checksums
│   ├── main.go                # Main application file
│   └── migrations/            # Database migrations
│       └── 001_init.sql       # Initial schema and data
└── frontend/                   # React frontend service
    ├── Dockerfile             # Frontend container config
    ├── package.json           # Node.js dependencies
    ├── tsconfig.json          # TypeScript configuration
    ├── public/                # Static assets
    │   └── index.html         # Main HTML template
    └── src/                   # Source code
        ├── index.tsx          # Application entry point
        ├── App.tsx            # Main app component
        ├── types/             # TypeScript type definitions
        ├── api/               # API client and services
        └── components/        # React components
            ├── Navigation.tsx  # Main navigation
            ├── Home.tsx       # Landing page
            ├── Dashboard.tsx  # Main dashboard
            ├── Achievements.tsx # Achievement store
            ├── Leaderboard.tsx # Community rankings
            └── Profile.tsx    # User profile
```

## 🎮 How to Play

### Getting Started
1. **Register**: Create your superhero account with a username and email
2. **Dashboard**: Your mission control center showing daily tasks and progress
3. **Complete Tasks**: Click "Complete Task" when you finish a fitness challenge
4. **Earn Points**: Gain 5-60 points per completed task based on difficulty
5. **Level Up**: Progress through 5 character levels as you earn points

### Task Categories
- **💪 Strength**: Push-ups, squats, planks, lunges
- **❤️ Cardio**: Walking, jogging, jumping jacks, dancing
- **🤸 Flexibility**: Stretching, yoga, mobility exercises
- **🧘 Wellness**: Hydration, breathing exercises, mindfulness

### Achievement System
- **🦸‍♂️ Characters**: Rookie Hero → Fitness Apprentice → Health Guardian → Wellness Warrior → Fitness Champion → Ultimate Hero
- **💼 Job Titles**: Fitness Novice → Personal Trainer → Fitness Coach → Wellness Expert → Fitness Director → Health Guru
- **🏅 Badges**: Special achievements for consistency, category mastery, and milestones

## 🛠️ Development

### Backend Development
```bash
cd backend
go mod download
go run main.go
```

### Frontend Development
```bash
cd frontend
npm install
npm start
```

### Database Access
```bash
docker exec -it fithero-db psql -U fithero_user -d fithero
```

## 🔧 Configuration

### Environment Variables
- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database username (default: fithero_user)
- `DB_PASSWORD`: Database password (default: fithero_password)
- `DB_NAME`: Database name (default: fithero)
- `PORT`: Backend server port (default: 8080)
- `REACT_APP_API_URL`: Frontend API URL (default: http://localhost:8080)

## 📊 API Endpoints

### Users
- `POST /api/users` - Create new user
- `GET /api/users/:id` - Get user by ID
- `PUT /api/users/:id` - Update user

### Tasks
- `GET /api/tasks` - Get all available tasks
- `GET /api/tasks/daily/:user_id` - Get user's daily tasks
- `POST /api/tasks/daily` - Generate new daily tasks
- `PUT /api/tasks/daily/:id/complete` - Complete a task

### Achievements
- `GET /api/achievements` - Get all achievements
- `GET /api/achievements/user/:user_id` - Get user's achievements
- `POST /api/achievements/unlock` - Unlock an achievement

### Leaderboard
- `GET /api/leaderboard` - Get top users

## 🎯 Target Audience

**Primary**: Sedentary individuals looking to start their fitness journey
**Secondary**: Fitness enthusiasts who enjoy gamification
**Age Range**: All ages (family-friendly design)
**Experience Level**: Beginners to intermediate fitness levels

## 🌈 Design Philosophy

- **Cute & Approachable**: Friendly superhero theme removes fitness intimidation
- **Progressive**: Start with easy 5-minute tasks, gradually increase difficulty
- **Rewarding**: Immediate positive feedback through points and achievements
- **Social**: Community aspect encourages continued participation
- **Flexible**: Tasks accommodate different fitness levels and preferences

## 🚀 Future Enhancements

- **Mobile App**: Native iOS/Android applications
- **Wearable Integration**: Sync with fitness trackers
- **Social Features**: Friend connections and challenges
- **Custom Tasks**: User-created fitness challenges
- **Nutrition Tracking**: Meal logging and nutrition goals
- **Workout Plans**: Structured multi-day fitness programs
- **AI Recommendations**: Personalized task suggestions
- **Achievements Expansion**: More badges and rewards

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by fitness gamification research and successful apps like Strava
- Built with modern web technologies for optimal performance
- Designed with accessibility and inclusivity in mind
- Community-driven development approach

---

**Ready to become a fitness superhero? Start your journey today! 🦸‍♂️💪**

## Environment Configuration

### Security Best Practices

This application uses environment variables to manage configuration and keep sensitive information secure. **Never commit the `.env` file to version control.**

### Setup Instructions

1. **Copy the example environment file:**
   ```bash
   cp env.example .env
   ```

2. **Edit the `.env` file with your configuration:**
   ```bash
   # Database Configuration
   POSTGRES_DB=fithero
   POSTGRES_USER=fithero_user
   POSTGRES_PASSWORD=your_secure_password_here
   POSTGRES_PORT=5432

   # Backend Configuration
   DB_HOST=db
   DB_PORT=5432
   DB_USER=fithero_user
   DB_PASSWORD=your_secure_password_here
   DB_NAME=fithero
   BACKEND_PORT=8080

   # Frontend Configuration
   VITE_API_URL=http://localhost:8080
   FRONTEND_PORT=3000

   # Project Configuration
   COMPOSE_PROJECT_NAME=fithero
   ```

3. **Important Security Notes:**
   - The `.env` file is already included in `.gitignore`
   - Use strong, unique passwords for `POSTGRES_PASSWORD` and `DB_PASSWORD`
   - For production, consider using Docker secrets or a secure secret management service
   - Never hardcode credentials in source code or Docker Compose files

## Troubleshooting

### Common Issues

1. **Port conflicts**: Make sure ports 3000, 8080, and 5432 are available
2. **Environment variables**: Ensure `.env` file exists and contains all required variables
3. **Docker issues**: Try `docker-compose down` and `docker-compose up --build`
4. **Database connection**: Verify database credentials in `.env` file

### Logs

View service logs:
```bash
docker-compose logs frontend
docker-compose logs backend
docker-compose logs db
```

For more detailed information about specific components, refer to the individual service documentation.