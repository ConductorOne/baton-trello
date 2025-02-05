package connector

import (
	"fmt"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"strconv"
)

func getToken(pToken *pagination.Token, resourceType *v2.ResourceType) (*pagination.Bag, int, error) {
	var pageToken int
	_, bag, err := unmarshalSkipToken(pToken)
	if err != nil {
		return bag, 0, err
	}

	if bag.Current() == nil {
		bag.Push(pagination.PageState{
			ResourceTypeID: resourceType.Id,
		})
	}

	if bag.Current().Token != "" {
		pageToken, err = strconv.Atoi(bag.Current().Token)
		if err != nil {
			return bag, 0, err
		}
	}

	return bag, pageToken, nil
}

func unmarshalSkipToken(token *pagination.Token) (int32, *pagination.Bag, error) {
	b := &pagination.Bag{}
	err := b.Unmarshal(token.Token)
	if err != nil {
		return 0, nil, err
	}
	current := b.Current()
	skip := int32(0)
	if current != nil && current.Token != "" {
		skip64, err := strconv.ParseInt(current.Token, 10, 32)
		if err != nil {
			return 0, nil, err
		}
		skip = int32(skip64)
	}
	return skip, b, nil
}

// extractRateLimitData returns a set of annotations for rate limiting given the rate limit headers provided by Trello API.
func extractRateLimitData(response *http.Response) (*v2.RateLimitDescription, error) {
	if response == nil {
		return nil, fmt.Errorf("trello-connector: passed nil response")
	}
	var err error

	var r int64
	remaining := response.Header.Get("X-Ratelimit-Remaining")
	if remaining != "" {
		r, err = strconv.ParseInt(remaining, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ratelimit-remaining: %w", err)
		}
	}

	var l int64
	limit := response.Header.Get("X-Ratelimit-Limit")
	if limit != "" {
		l, err = strconv.ParseInt(limit, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ratelimit-limit: %w", err)
		}
	}

	var ra *timestamppb.Timestamp
	resetAt := response.Header.Get("X-Ratelimit-Reset")
	if resetAt != "" {
		ts, err := strconv.ParseInt(resetAt, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ratelimit-reset: %w", err)
		}
		ra = &timestamppb.Timestamp{Seconds: ts}
	}

	return &v2.RateLimitDescription{
		Limit:     l,
		Remaining: r,
		ResetAt:   ra,
	}, nil
}
