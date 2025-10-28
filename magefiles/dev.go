package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"chef/core/auth"
	"chef/core/cli"
	"chef/core/randx"
	"chef/data"

	"github.com/jackc/pgx/v5/pgtype"
)

// Bootstrap creates a user in the specified environment
func Bootstrap() error {
	ctx := context.Background()

	d, err := connect()
	if err != nil {
		return err
	}

	var email string
	cli.InteractiveStr("Email: ", &email)

	var isAdmin bool
	cli.InteractiveYesNo("Is Admin (y/n): ", &isAdmin)

	var isRoot bool
	cli.InteractiveYesNo("Is Root (y/n): ", &isRoot)

	account, err := d.q.AccountGetByEmail(ctx, email)
	if err != nil {
		return err
	}

	if isRoot {
		if err := d.q.AccountSetRoot(ctx, data.AccountSetRootParams{
			ID:     account.ID,
			IsRoot: isRoot,
		}); err != nil {
			return err
		}
	} else if isAdmin {
		if err := d.q.AccountSetAdmin(ctx, data.AccountSetAdminParams{
			ID:      account.ID,
			IsAdmin: isAdmin,
		}); err != nil {
			return err
		}
	}

	fmt.Println("done.")

	return nil
}

func GenBearer() error {
	ctx := context.Background()
	now := time.Now().UTC()

	d, err := connect()
	if err != nil {
		return err
	}

	var email string
	cli.InteractiveStr("Email: ", &email)

	acc, err := d.q.AccountGetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("Account not found.")
			return nil
		}
		return err
	}

	sessionID := randx.UID()

	expAt := now.Add(auth.SessionDuration)
	randToken := randx.UID()

	tokenHash := d.hasher.Hash(randToken)

	if err := d.q.SessionCreate(ctx, data.SessionCreateParams{
		ID:        sessionID,
		AccountID: acc.ID,
		ExpiresAt: pgtype.Timestamp{Time: expAt, Valid: true},
		TokenHash: tokenHash,
		CreatedAt: pgtype.Timestamp{Time: now, Valid: true},
	}); err != nil {
		return err
	}

	fmt.Println(randToken)

	return nil
}
