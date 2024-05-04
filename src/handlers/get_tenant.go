package handlers

import (
	"context"
	"errors"
	"fmt"

	"gettenant/db"
	"gettenant/models"
)

func HandleRequest(ctx context.Context, event *models.Request) (*models.Response, error) {
	if event == nil {
		return nil, fmt.Errorf("received nil event")
	}

	if event.Id == "" {
		return nil, errors.New("id is required")
	}

	item, err := db.GetTenant(ctx, event)

	return item, err
}
