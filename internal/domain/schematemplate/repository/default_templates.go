package repository

import (
	"encoding/json"
	"time"

	"progressive/internal/domain/schematemplate"
)

// TemplateDefinition represents a type-safe template definition
type TemplateDefinition struct {
	ID          string
	Name        string
	Description string
	Category    string
	Icon        string
	Schema      SchemaDefinition
}

// SchemaDefinition represents the JSON schema structure
type SchemaDefinition struct {
	Type       string                 `json:"type"`
	Properties map[string]PropertyDef `json:"properties"`
	Required   []string               `json:"required,omitempty"`
}

// PropertyDef represents individual property definitions
type PropertyDef struct {
	Type        string   `json:"type"`
	Title       string   `json:"title,omitempty"`
	Format      string   `json:"format,omitempty"`
	Pattern     string   `json:"pattern,omitempty"`
	MinLength   int      `json:"minLength,omitempty"`
	Minimum     int      `json:"minimum,omitempty"`
	Maximum     int      `json:"maximum,omitempty"`
	Enum        []string `json:"enum,omitempty"`
	Description string   `json:"description,omitempty"`
}

// GetDefaultTemplateDefinitions returns compile-time safe template definitions
func GetDefaultTemplateDefinitions() []TemplateDefinition {
	return []TemplateDefinition{
		// Business Templates
		{
			ID:          "customer",
			Name:        "customer_management",
			Description: "고객 정보와 연락처를 관리하는 템플릿",
			Category:    schematemplate.CategoryBusiness,
			Icon:        "👥",
			Schema: SchemaDefinition{
				Type: "object",
				Properties: map[string]PropertyDef{
					"name":              {Type: "string", Title: "이름", MinLength: 1},
					"email":             {Type: "string", Format: "email", Title: "이메일"},
					"company":           {Type: "string", Title: "회사명"},
					"phone":             {Type: "string", Title: "연락처", Pattern: "^[0-9-+()\\s]+$"},
					"interest_level":    {Type: "string", Title: "관심도", Enum: []string{"높음", "중간", "낮음"}},
					"registration_date": {Type: "string", Format: "date", Title: "등록일"},
				},
				Required: []string{"name", "email"},
			},
		},
		{
			ID:          "project",
			Name:        "project_management",
			Description: "프로젝트 진행 상황과 일정을 추적하는 템플릿",
			Category:    schematemplate.CategoryBusiness,
			Icon:        "📋",
			Schema: SchemaDefinition{
				Type: "object",
				Properties: map[string]PropertyDef{
					"title":       {Type: "string", Title: "제목", MinLength: 1},
					"description": {Type: "string", Title: "설명"},
					"assignee":    {Type: "string", Title: "담당자"},
					"status":      {Type: "string", Title: "상태", Enum: []string{"TODO", "진행중", "완료"}},
					"priority":    {Type: "integer", Title: "우선순위", Minimum: 1, Maximum: 5},
					"due_date":    {Type: "string", Format: "date", Title: "마감일"},
				},
				Required: []string{"title", "status"},
			},
		},
		{
			ID:          "inventory",
			Name:        "inventory_management",
			Description: "상품 재고와 공급업체 정보를 관리하는 템플릿",
			Category:    schematemplate.CategoryBusiness,
			Icon:        "📦",
			Schema: SchemaDefinition{
				Type: "object",
				Properties: map[string]PropertyDef{
					"product_name":  {Type: "string", Title: "상품명", MinLength: 1},
					"category":      {Type: "string", Title: "카테고리"},
					"quantity":      {Type: "integer", Title: "수량", Minimum: 0},
					"price":         {Type: "number", Title: "가격", Minimum: 0},
					"supplier":      {Type: "string", Title: "공급업체"},
					"last_updated":  {Type: "string", Format: "date-time", Title: "최종 업데이트"},
				},
				Required: []string{"product_name", "quantity", "price"},
			},
		},
		{
			ID:          "event",
			Name:        "event_management",
			Description: "이벤트와 일정을 관리하는 템플릿",
			Category:    schematemplate.CategoryBusiness,
			Icon:        "📅",
			Schema: SchemaDefinition{
				Type: "object",
				Properties: map[string]PropertyDef{
					"title":     {Type: "string", Title: "제목", MinLength: 1},
					"date":      {Type: "string", Format: "date", Title: "날짜"},
					"time":      {Type: "string", Title: "시간", Pattern: "^([01]?[0-9]|2[0-3]):[0-5][0-9]$"},
					"location":  {Type: "string", Title: "장소"},
					"attendees": {Type: "integer", Title: "참석자 수", Minimum: 0},
					"type":      {Type: "string", Title: "이벤트 유형", Enum: []string{"회의", "워크샵", "세미나", "파티", "기타"}},
				},
				Required: []string{"title", "date", "time"},
			},
		},
		// Game Templates
		{
			ID:          "quest",
			Name:        "quest_management",
			Description: "게임 퀘스트와 보상을 관리하는 템플릿",
			Category:    schematemplate.CategoryGame,
			Icon:        "⚔️",
			Schema: SchemaDefinition{
				Type: "object",
				Properties: map[string]PropertyDef{
					"quest_name":           {Type: "string", Title: "퀘스트명", MinLength: 1},
					"description":          {Type: "string", Title: "설명"},
					"quest_type":           {Type: "string", Title: "퀘스트 유형", Enum: []string{"메인", "서브", "일일", "주간", "이벤트"}},
					"difficulty":           {Type: "string", Title: "난이도", Enum: []string{"쉬움", "보통", "어려움", "매우어려움"}},
					"level_requirement":    {Type: "integer", Title: "필요 레벨", Minimum: 1, Maximum: 100},
					"reward_exp":           {Type: "integer", Title: "보상 경험치", Minimum: 0},
					"reward_gold":          {Type: "integer", Title: "보상 골드", Minimum: 0},
					"reward_items":         {Type: "string", Title: "보상 아이템"},
					"completion_condition": {Type: "string", Title: "완료 조건"},
					"status":               {Type: "string", Title: "상태", Enum: []string{"활성", "비활성", "테스트중"}},
				},
				Required: []string{"quest_name", "quest_type", "difficulty", "level_requirement"},
			},
		},
		{
			ID:          "shop_item",
			Name:        "shop_item_management",
			Description: "게임 내 상점 아이템을 관리하는 템플릿",
			Category:    schematemplate.CategoryGame,
			Icon:        "🛍️",
			Schema: SchemaDefinition{
				Type: "object",
				Properties: map[string]PropertyDef{
					"item_name":         {Type: "string", Title: "상품명", MinLength: 1},
					"description":       {Type: "string", Title: "설명"},
					"category":          {Type: "string", Title: "카테고리", Enum: []string{"무기", "방어구", "소모품", "장식품", "재료", "기타"}},
					"rarity":            {Type: "string", Title: "등급", Enum: []string{"일반", "고급", "희귀", "영웅", "전설"}},
					"price_gold":        {Type: "integer", Title: "골드 가격", Minimum: 0},
					"price_gem":         {Type: "integer", Title: "보석 가격", Minimum: 0},
					"stock":             {Type: "integer", Title: "재고", Minimum: -1},
					"level_requirement": {Type: "integer", Title: "필요 레벨", Minimum: 1, Maximum: 100},
					"is_limited":        {Type: "boolean", Title: "한정 상품"},
					"sale_start_date":   {Type: "string", Format: "date", Title: "판매 시작일"},
					"sale_end_date":     {Type: "string", Format: "date", Title: "판매 종료일"},
				},
				Required: []string{"item_name", "category", "rarity"},
			},
		},
		{
			ID:          "game_item",
			Name:        "game_item_management",
			Description: "게임 아이템의 상세 정보를 관리하는 템플릿",
			Category:    schematemplate.CategoryGame,
			Icon:        "🗡️",
			Schema: SchemaDefinition{
				Type: "object",
				Properties: map[string]PropertyDef{
					"item_name":           {Type: "string", Title: "아이템명", MinLength: 1},
					"description":         {Type: "string", Title: "설명"},
					"item_type":           {Type: "string", Title: "아이템 유형", Enum: []string{"무기", "방어구", "악세서리", "소모품", "재료", "퀘스트", "기타"}},
					"rarity":              {Type: "string", Title: "등급", Enum: []string{"일반", "고급", "희귀", "영웅", "전설", "신화"}},
					"level_requirement":   {Type: "integer", Title: "필요 레벨", Minimum: 1, Maximum: 100},
					"attack_power":        {Type: "integer", Title: "공격력", Minimum: 0},
					"defense_power":       {Type: "integer", Title: "방어력", Minimum: 0},
					"hp_bonus":            {Type: "integer", Title: "체력 보너스", Minimum: 0},
					"mp_bonus":            {Type: "integer", Title: "마나 보너스", Minimum: 0},
					"special_effect":      {Type: "string", Title: "특수 효과"},
					"durability":          {Type: "integer", Title: "내구도", Minimum: 0, Maximum: 100},
					"max_stack":           {Type: "integer", Title: "최대 중첩", Minimum: 1, Maximum: 999},
					"drop_location":       {Type: "string", Title: "획득 장소"},
					"crafting_materials":  {Type: "string", Title: "제작 재료"},
				},
				Required: []string{"item_name", "item_type", "rarity"},
			},
		},
	}
}

// toDomainModel converts a TemplateDefinition to a SchemaTemplate domain model
func (def *TemplateDefinition) toDomainModel() (*schematemplate.SchemaTemplate, error) {
	schemaJSON, err := json.Marshal(def.Schema)
	if err != nil {
		return nil, err
	}

	return &schematemplate.SchemaTemplate{
		ID:          def.ID,
		Name:        def.Name,
		Description: def.Description,
		Category:    def.Category,
		Icon:        def.Icon,
		Schema:      json.RawMessage(schemaJSON),
		SampleData:  nil, // No sample data for now
		CreatedAt:   time.Now(),
	}, nil
}