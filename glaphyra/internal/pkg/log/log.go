package log

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/rs/zerolog"
)

var clickhouseConn clickhouse.Conn
var Logger zerolog.Logger

func init() {
	// Настройка формата времени для zerolog
	zerolog.TimeFieldFormat = time.RFC3339

	// Создаем глобальный логгер
	Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Добавляем хук для записи логов в ClickHouse
	Logger = Logger.Hook(&ClickHouseHook{})

	// Устанавливаем уровень логирования
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Инициализация ClickHouse
	if err := initClickHouse(); err != nil {
		panic(fmt.Errorf("failed to initialize ClickHouse: %w", err))
	}
}

func initClickHouse() error {
	var err error
	clickhouseConn, err = clickhouse.Open(&clickhouse.Options{
		Addr: []string{"127.0.0.1:9000"},
		Auth: clickhouse.Auth{
			Database: "test",
			Username: "test",
			Password: "test",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	// Проверка подключения
	if err := clickhouseConn.Ping(context.Background()); err != nil {
		return fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	return nil
}

type ClickHouseHook struct{}

func (h *ClickHouseHook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	ctx := context.Background()

	// Формируем запрос для записи лога в ClickHouse
	query := `
        INSERT INTO logs (timestamp, level, message)
        VALUES (?, ?, ?)
    `

	// Выполняем запрос
	if err := clickhouseConn.Exec(ctx, query,
		time.Now(),
		level.String(),
		message,
	); err != nil {
		fmt.Println("Failed to insert log into ClickHouse:", err)
	}
}

// WriteLog записывает информационное сообщение
func WriteLog(msg string) {
	Logger.Info().Msg(msg)
}

// WriteLogf записывает форматированное информационное сообщение
func WriteLogf(format string, args ...interface{}) {
	Logger.Info().Msgf(format, args...)
}

// Error записывает сообщение об ошибке с уровнем "error"
func Error(err error) {
	if err == nil {
		return
	}
	Logger.Error().Err(err).Msg(err.Error())
}

// WrapErr оборачивает ошибку, записывает её в лог и возвращает
func WrapErr(err error) error {
	if err == nil {
		return nil
	}
	Logger.Error().Err(err).Msg(err.Error())
	return err
}

// Wrap оборачивает ошибку без записи в лог
func Wrap(err error) error {
	return fmt.Errorf("%w", err)
}

func LogCommand(tgID int64, command string) {
	ctx := context.Background()

	query := `
        INSERT INTO history (tg_id, command, created_at)
        VALUES (?, ?, ?)
    `

	// Выполняем запрос
	if err := clickhouseConn.Exec(ctx, query,
		tgID,
		command,
		time.Now(),
	); err != nil {
		fmt.Println("Failed to insert command log into ClickHouse:", err)
	}
}
