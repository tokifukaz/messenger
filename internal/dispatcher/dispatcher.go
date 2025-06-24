package dispatcher

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "strings"
    "sync"
    "time"
	"fmt"

	"github.com/redis/go-redis/v9"

    "messenger/internal/db"
)

type Dispatcher struct {
    repo       *db.MessageRepository
    redis      *redis.Client
    ctx        context.Context
    cancel     context.CancelFunc
    running    bool
    mu         sync.Mutex
}

func New(repo *db.MessageRepository, redisClient *redis.Client) *Dispatcher {
    return &Dispatcher{
        repo:  repo,
        redis: redisClient,
    }
}

func (d *Dispatcher) Start() {
    d.mu.Lock()
    defer d.mu.Unlock()

    if d.running {
        log.Println("Dispatcher already running")
        return
    }

    d.ctx, d.cancel = context.WithCancel(context.Background())
    d.running = true

    go func() {
        ticker := time.NewTicker(2 * time.Minute)
        defer ticker.Stop()

        for {
            select {
            case <-d.ctx.Done():
                log.Println("Dispatcher stopped")
                return
            default:
                d.processMessages()
                <-ticker.C
            }
        }
    }()
}

func (d *Dispatcher) Stop() {
    d.mu.Lock()
    defer d.mu.Unlock()

    if d.cancel != nil {
        d.cancel()
        d.running = false
    }
}

func (d *Dispatcher) IsRunning() bool {
    d.mu.Lock()
    defer d.mu.Unlock()
    return d.running
}

func (d *Dispatcher) processMessages() {
    messages, err := d.repo.GetUnsentMessages(2)
    if err != nil {
        log.Println("Error fetching messages:", err)
        return
    }

    for _, msg := range messages {
        success, messageId := d.sendMessage(msg.PhoneNumber, msg.Content)
        if success {
            _ = d.repo.MarkMessageAsSent(msg.ID)

            // Cache in Redis
            d.redis.HSet(context.Background(), fmt.Sprintf("message:%d", msg.ID), map[string]interface{}{
                "messageId": messageId,
                "sentAt":    time.Now().Format(time.RFC3339),
            })
        }
    }
}

func (d *Dispatcher) sendMessage(phone, content string) (bool, string) {
    payload := map[string]string{
        "to":      phone,
        "message": content,
    }

    body, _ := json.Marshal(payload)
    resp, err := http.Post("https://webhook.site/2d55a82e-d299-4260-a41f-e859b27bc583", "application/json", strings.NewReader(string(body)))
    if err != nil {
        return false, ""
    }
    defer resp.Body.Close()

    var res map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&res)
	log.Printf("Response from webhook: %v", res)

    return resp.StatusCode == 200, res["messageId"].(string)
}
