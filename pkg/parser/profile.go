package parser

import (
	"errors"
	"strings"

	logger "core/pkg/log"

	"github.com/PuerkitoBio/goquery"
)

type Profile struct {
	Name     string
	Location string
	Country  string
	URL      string
}

func ParseProfile(html string) (*Profile, error) {
	document, err := goquery.
		NewDocumentFromReader(
			strings.NewReader(html),
		)

	if err != nil {
		return nil, err
	}

	profile := &Profile{}

	document.
		Find("a").
		EachWithBreak(func(_ int, selection *goquery.Selection ) bool {
			if !strings.Contains(
				strings.ToLower(
					strings.TrimSpace(
						selection.Text(),
					),
				),
				"contact info",
			) {
				return true
			}

			container := selection.Parent().Parent()
			var parts []string

			container.
				ChildrenFiltered("p").
				Each(func(_ int, p *goquery.Selection) {
					text := strings.TrimSpace(
						p.Text(),
					)

					if text == "" ||
						text == "·" ||
						strings.Contains(
							strings.ToLower(text),
							"contact info",
						) {
						return
					}

					parts = append(
						parts,
						text,
					)
				})

			if len(parts) > 0 {
				profile.Location = strings.Join(parts, ", ")
				logger.Info(
					"Found location: %s",
					profile.Location,
				)

				return false
			}

			return true
		})

	if profile.Name == "" {
		profile.Name = strings.TrimSpace(
			document.
				Find("h1,h2").
				First().
				Text(),
		)
	}

	if profile.Location == "" {
		logger.Error(
			"Failed to parse location",
		)

		return nil, errors.New(
			"location not found",
		)
	}

	logger.Success(
		"Parsed profile: %s (%s)",
		profile.Name,
		profile.Location,
	)

	return profile, nil
}