package bilibili

import "context"

type SearchService struct {
	client *Client
}

type SearchType string

const (
	SearchTypeVideo    SearchType = "video"
	SearchTypeBangumi  SearchType = "media_bangumi"
	SearchTypeFT       SearchType = "media_ft"
	SearchTypeUser     SearchType = "bili_user"
	SearchTypeLiveRoom SearchType = "live_room"
	SearchTypeTopic    SearchType = "topic"
)

func (s *SearchService) All(ctx context.Context, keyword string, page int) (*SearchAllResult, error) {
	var out SearchAllResult
	req := s.client.NewRequest(endpointSearchAll).Param("keyword", keyword)
	if page > 0 {
		req.ParamInt("page", int64(page))
	}
	if err := req.Do(ctx, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *SearchService) ByType(ctx context.Context, typ SearchType, keyword string, page int) (*SearchTypeResult, error) {
	var out SearchTypeResult
	req := s.client.NewRequest(endpointSearchType).
		Param("search_type", string(typ)).
		Param("keyword", keyword)
	if page > 0 {
		req.ParamInt("page", int64(page))
	}
	if err := req.Do(ctx, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *SearchService) Suggest(ctx context.Context, keyword string) ([]SearchSuggestItem, error) {
	var payload map[string]SearchSuggestItem
	err := s.client.NewRequest(endpointSearchSuggest).
		Param("term", keyword).
		Do(ctx, &payload)
	if err != nil {
		return nil, err
	}

	result := make([]SearchSuggestItem, 0, len(payload))
	for _, item := range payload {
		result = append(result, item)
	}
	return result, nil
}

type SearchPage struct {
	NumResults int `json:"numResults"`
	NumPages   int `json:"numPages"`
	Page       int `json:"page"`
	PageSize   int `json:"pagesize"`
}

type SearchAllResult struct {
	SearchPage
	Result []struct {
		ResultType string                   `json:"result_type"`
		Data       []map[string]interface{} `json:"data"`
	} `json:"result"`
}

type SearchTypeResult struct {
	SearchPage
	Result []map[string]interface{} `json:"result"`
}

type SearchSuggestItem struct {
	Value string `json:"value"`
	Ref   int    `json:"ref"`
	Name  string `json:"name"`
	Term  string `json:"term"`
	Spid  int    `json:"spid"`
}
