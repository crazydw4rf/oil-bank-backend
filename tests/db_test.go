package tests

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/crazydw4rf/oil-bank-backend/internal/services/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestDB(t *testing.T) {
	cfg, err := config.InitConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error woyy: %#v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*150)
	defer cancel()

	conn, err := pgxpool.New(ctx, cfg.DATABASE_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "1: Error woyy: %#v\n", err)
		return
	}
	defer conn.Close()

	rows, err := conn.Query(ctx, fmt.Sprintf(`
		SELECT '%s'::uuid AS id, 1 AS num
		UNION ALL
		SELECT '%s'::uuid AS id, 2 AS num
		UNION ALL
		SELECT '%s'::uuid AS id, 5 AS num
	`, uuid.New(), uuid.New(), uuid.New()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "2: Error woyy: %#v\n", err)
	}

	// var ID uuid.UUID
	// var num int

	type FooStruct struct {
		Id   uuid.UUID `db:"id"`
		Num  int       `db:"num"`
		Num2 int       `db:"-"`
	}

	res, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[FooStruct])
	if err != nil {
		fmt.Fprintf(os.Stderr, "3: Error woyy: %#v\n", err)
		return
	}

	fmt.Printf("%v\n", res)

	// for row.Next() {
	// 	if err = row.Err(); err != nil {
	// 		fmt.Fprintf(os.Stderr, "4: Error woyy: %#v\n", err)
	// 		return
	// 	}

	// 	err = row.Scan(&ID, &num)
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stderr, "5: Error woyy: %#v\n", err)
	// 		return
	// 	}

	// 	fmt.Printf("ID: %v, num: %d\n", ID, num)
	// }
}
