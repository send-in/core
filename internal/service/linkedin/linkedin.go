package service

import (
	model "core/internal/model"
	logger "core/pkg/log"

	"fmt"
	"net/http"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EnrichmentRequest struct {
	Token      string
	Agent      string
	Profile    string
	
	// TODO(v1.1):
	// Add Account.JSession migration and use it
	// for csrf-token + JSESSIONID cookie.
	JSession   string
	AccountID  uuid.UUID 
}

var EnrichmentJobs = make(
	chan EnrichmentRequest,
	100,
)

func LinkedinService(db *gorm.DB) {
	go func() {
		for job := range EnrichmentJobs {
			logger.Info("🔎 Enriching %s", job.Profile)

			start := 0

			for {
				request, err := http.NewRequest(
					http.MethodGet,
					fmt.Sprintf(
						"https://www.linkedin.com/voyager/api/relationships/connections?count=100&start=%d",
						start,
					), nil,
				)
				
				if err != nil { break }

				request.Header = http.Header{
					"User-Agent": []string{
						job.Agent,
					},
					"csrf-token": []string{
						job.JSession,
					},
					"Cookie": []string{
						fmt.Sprintf(
							"li_at=%s; JSESSIONID=%s",
							job.Token,
							job.JSession,
						),
					},
				}

				response, err := http.DefaultClient.Do(request)
				if err != nil { break }

				var payload ConnectionsResponse
				err = json.
					  NewDecoder(response.Body).
					  Decode(&payload)

				response.Body.Close()
				if err != nil { break }

				if len(payload.Elements) == 0 {
					logger.Info("All Connections Fetched")
					break
				}

				for _, element := range payload.Elements {
					connection := LinkedinConnection(
						element,
					)

					err := db.Transaction(
						func(tx *gorm.DB) error {
							if err := tx.
								Where(
									"public_id = ?",
									connection.PublicID,
								).
								Assign(model.Connection{
									FirstName: connection.FirstName,
									LastName:  connection.LastName,
									Bio:       connection.Bio,
									Picture:   connection.Picture,
									Company:   connection.Company,
									Country:   connection.Country,
									Timezone:  connection.Timezone,
								}).
								FirstOrCreate(
									&connection,
								).
								Error; err != nil {
								return err
							}

							return tx.
								FirstOrCreate(
									&model.AccountConnection{},
									model.AccountConnection{
										AccountID:    job.AccountID,
										ConnectionID: connection.ID,
									},
								).
								Error
						},
					)

					if err != nil {
						logger.Error(
							"Failed to sync %s: %v",
							connection.PublicID,
							err,
						)
					}
				}

				start += payload.Paging.Count
			}

			logger.Success(
				"📨 Enrichment completed for %s",
				job.Profile,
			)
		}
	}()
}