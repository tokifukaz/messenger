package handler

import (
	"database/sql"
    "log"
    "net/http"
    "os"
    "fmt"
    "time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	_ "github.com/lib/pq"

	"messenger/internal/db"
    "messenger/internal/dispatcher"
)

type Handler struct {
    dispatcher *dispatcher.Dispatcher
    repo       *db.MessageRepository
    redis      *redis.Client
}

func NewHandler() *Handler {
    dbConn := InitDB()
    repo := &db.MessageRepository{DB: dbConn}
    redisClient := InitRedis()
    
    if redisClient == nil {
        log.Fatal("Failed to connect to Redis")
    }

    return &Handler{
        dispatcher: dispatcher.New(repo, redisClient),
        repo:       repo,
        redis:      redisClient,
    }
}

// StartDispatcher godoc
// @Summary Start auto dispatcher
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/start [post]
func (h *Handler) StartDispatcher(c *gin.Context) {
    h.dispatcher.Start()
    c.JSON(http.StatusOK, gin.H{"status": "started"})
}

// StopDispatcher godoc
// @Summary Stop auto dispatcher
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/stop [post]
func (h *Handler) StopDispatcher(c *gin.Context) {
    h.dispatcher.Stop()
    c.JSON(http.StatusOK, gin.H{"status": "stopped"})
}

// GetSentMessages godoc
// @Summary Get sent messages
// @Produce json
// @Success 200 {array} model.Message
// @Router /api/sent-messages [get]
func (h *Handler) GetSentMessages(c *gin.Context) {
    messages, err := h.repo.GetSentMessages()
    if err != nil {
        log.Println("Error fetching sent messages:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch sent messages"})
        return
    }
    c.JSON(http.StatusOK, messages)
}


var RedisClient *redis.Client

func InitDB() *sql.DB {
    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
    )

    var db *sql.DB
    var err error

    for attempts := 1; attempts <= 10; attempts++ {
        db, err = sql.Open("postgres", connStr)
        if err == nil {
            if err = db.Ping(); err == nil {
                // Check if table exists
                if checkMessagesTable(db) {
                    log.Println("DB connection established and messages table exists.")
                    return db
                }
            }
        }
        log.Printf("Waiting for DB (%d/10)...", attempts)
        time.Sleep(3 * time.Second)
    }

    log.Fatal("Could not connect to DB or table does not exist:", err)
    return nil
}

func checkMessagesTable(db *sql.DB) bool {
    var exists bool
    query := `SELECT EXISTS (
        SELECT FROM information_schema.tables 
        WHERE table_name = 'messages'
    )`
    err := db.QueryRow(query).Scan(&exists)
    if err != nil {
        log.Println("Error checking messages table:", err)
        return false
    }
    return exists
}

func InitRedis() *redis.Client {
    RedisClient = redis.NewClient(&redis.Options{
        Addr: os.Getenv("REDIS_ADDR"),
    })
    return RedisClient
}