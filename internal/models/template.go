package models

// TableTemplate represents a predefined table template
type TableTemplate struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Category    string      `json:"category"`
	Icon        string      `json:"icon"`
	Schema      interface{} `json:"schema"`
}

// GetTableTemplates returns all available table templates
func GetTableTemplates() []TableTemplate {
	return []TableTemplate{
		// Business Templates
		{
			ID:          "customer",
			Name:        "고객 관리",
			Description: "이름, 이메일, 회사명, 연락처 필드를 포함한 고객 정보 테이블",
			Category:    "business",
			Icon:        "users",
			Schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type":      "string",
						"title":     "이름",
						"minLength": 1,
					},
					"email": map[string]interface{}{
						"type":   "string",
						"format": "email",
						"title":  "이메일",
					},
					"company": map[string]interface{}{
						"type":  "string",
						"title": "회사명",
					},
					"phone": map[string]interface{}{
						"type":    "string",
						"title":   "연락처",
						"pattern": "^[0-9-+()\\s]+$",
					},
					"interest_level": map[string]interface{}{
						"type":  "string",
						"title": "관심도",
						"enum":  []string{"높음", "중간", "낮음"},
					},
					"registration_date": map[string]interface{}{
						"type":   "string",
						"format": "date",
						"title":  "등록일",
					},
				},
				"required": []string{"name", "email"},
			},
		},
		{
			ID:          "project",
			Name:        "프로젝트 관리",
			Description: "제목, 설명, 담당자, 상태, 마감일 필드를 포함한 프로젝트 추적 테이블",
			Category:    "business",
			Icon:        "briefcase",
			Schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{
						"type":      "string",
						"title":     "제목",
						"minLength": 1,
					},
					"description": map[string]interface{}{
						"type":  "string",
						"title": "설명",
					},
					"assignee": map[string]interface{}{
						"type":  "string",
						"title": "담당자",
					},
					"status": map[string]interface{}{
						"type":  "string",
						"title": "상태",
						"enum":  []string{"TODO", "진행중", "완료"},
					},
					"priority": map[string]interface{}{
						"type":    "integer",
						"title":   "우선순위",
						"minimum": 1,
						"maximum": 5,
					},
					"due_date": map[string]interface{}{
						"type":   "string",
						"format": "date",
						"title":  "마감일",
					},
				},
				"required": []string{"title", "status"},
			},
		},
		{
			ID:          "inventory",
			Name:        "재고 관리",
			Description: "상품명, 카테고리, 수량, 가격, 공급업체 필드를 포함한 재고 테이블",
			Category:    "business",
			Icon:        "package",
			Schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"product_name": map[string]interface{}{
						"type":      "string",
						"title":     "상품명",
						"minLength": 1,
					},
					"category": map[string]interface{}{
						"type":  "string",
						"title": "카테고리",
					},
					"quantity": map[string]interface{}{
						"type":    "integer",
						"title":   "수량",
						"minimum": 0,
					},
					"price": map[string]interface{}{
						"type":    "number",
						"title":   "가격",
						"minimum": 0,
					},
					"supplier": map[string]interface{}{
						"type":  "string",
						"title": "공급업체",
					},
					"last_updated": map[string]interface{}{
						"type":   "string",
						"format": "date-time",
						"title":  "최종 업데이트",
					},
				},
				"required": []string{"product_name", "quantity", "price"},
			},
		},
		{
			ID:          "event",
			Name:        "이벤트 관리",
			Description: "제목, 날짜, 시간, 장소, 참석자 필드를 포함한 이벤트 테이블",
			Category:    "business",
			Icon:        "calendar",
			Schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{
						"type":      "string",
						"title":     "제목",
						"minLength": 1,
					},
					"date": map[string]interface{}{
						"type":   "string",
						"format": "date",
						"title":  "날짜",
					},
					"time": map[string]interface{}{
						"type":    "string",
						"title":   "시간",
						"pattern": "^([01]?[0-9]|2[0-3]):[0-5][0-9]$",
					},
					"location": map[string]interface{}{
						"type":  "string",
						"title": "장소",
					},
					"attendees": map[string]interface{}{
						"type":    "integer",
						"title":   "참석자 수",
						"minimum": 0,
					},
					"type": map[string]interface{}{
						"type":  "string",
						"title": "이벤트 유형",
						"enum":  []string{"회의", "워크샵", "세미나", "파티", "기타"},
					},
				},
				"required": []string{"title", "date", "time"},
			},
		},

		// Game Templates
		{
			ID:          "quest",
			Name:        "퀘스트 관리",
			Description: "제목, 설명, 보상, 난이도, 완료 조건을 포함한 게임 퀘스트 테이블",
			Category:    "game",
			Icon:        "sword",
			Schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"quest_name": map[string]interface{}{
						"type":      "string",
						"title":     "퀘스트명",
						"minLength": 1,
					},
					"description": map[string]interface{}{
						"type":  "string",
						"title": "설명",
					},
					"quest_type": map[string]interface{}{
						"type":  "string",
						"title": "퀘스트 유형",
						"enum":  []string{"메인", "서브", "일일", "주간", "이벤트"},
					},
					"difficulty": map[string]interface{}{
						"type":  "string",
						"title": "난이도",
						"enum":  []string{"쉬움", "보통", "어려움", "매우어려움"},
					},
					"level_requirement": map[string]interface{}{
						"type":    "integer",
						"title":   "필요 레벨",
						"minimum": 1,
						"maximum": 100,
					},
					"reward_exp": map[string]interface{}{
						"type":    "integer",
						"title":   "보상 경험치",
						"minimum": 0,
					},
					"reward_gold": map[string]interface{}{
						"type":    "integer",
						"title":   "보상 골드",
						"minimum": 0,
					},
					"reward_items": map[string]interface{}{
						"type":  "string",
						"title": "보상 아이템",
					},
					"completion_condition": map[string]interface{}{
						"type":  "string",
						"title": "완료 조건",
					},
					"status": map[string]interface{}{
						"type":  "string",
						"title": "상태",
						"enum":  []string{"활성", "비활성", "테스트중"},
					},
				},
				"required": []string{"quest_name", "quest_type", "difficulty", "level_requirement"},
			},
		},
		{
			ID:          "shop_item",
			Name:        "상품 관리",
			Description: "상품명, 가격, 카테고리, 재고를 포함한 게임 내 상점 아이템 테이블",
			Category:    "game",
			Icon:        "shopping-cart",
			Schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"item_name": map[string]interface{}{
						"type":      "string",
						"title":     "상품명",
						"minLength": 1,
					},
					"description": map[string]interface{}{
						"type":  "string",
						"title": "설명",
					},
					"category": map[string]interface{}{
						"type":  "string",
						"title": "카테고리",
						"enum":  []string{"무기", "방어구", "소모품", "장식품", "재료", "기타"},
					},
					"rarity": map[string]interface{}{
						"type":  "string",
						"title": "등급",
						"enum":  []string{"일반", "고급", "희귀", "영웅", "전설"},
					},
					"price_gold": map[string]interface{}{
						"type":    "integer",
						"title":   "골드 가격",
						"minimum": 0,
					},
					"price_gem": map[string]interface{}{
						"type":    "integer",
						"title":   "보석 가격",
						"minimum": 0,
					},
					"stock": map[string]interface{}{
						"type":    "integer",
						"title":   "재고",
						"minimum": -1,
					},
					"level_requirement": map[string]interface{}{
						"type":    "integer",
						"title":   "필요 레벨",
						"minimum": 1,
						"maximum": 100,
					},
					"is_limited": map[string]interface{}{
						"type":  "boolean",
						"title": "한정 상품",
					},
					"sale_start_date": map[string]interface{}{
						"type":   "string",
						"format": "date",
						"title":  "판매 시작일",
					},
					"sale_end_date": map[string]interface{}{
						"type":   "string",
						"format": "date",
						"title":  "판매 종료일",
					},
				},
				"required": []string{"item_name", "category", "rarity"},
			},
		},
		{
			ID:          "game_item",
			Name:        "게임 아이템",
			Description: "아이템명, 능력치, 등급, 획득 방법을 포함한 게임 아이템 데이터 테이블",
			Category:    "game",
			Icon:        "shield",
			Schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"item_name": map[string]interface{}{
						"type":      "string",
						"title":     "아이템명",
						"minLength": 1,
					},
					"description": map[string]interface{}{
						"type":  "string",
						"title": "설명",
					},
					"item_type": map[string]interface{}{
						"type":  "string",
						"title": "아이템 유형",
						"enum":  []string{"무기", "방어구", "악세서리", "소모품", "재료", "퀘스트", "기타"},
					},
					"rarity": map[string]interface{}{
						"type":  "string",
						"title": "등급",
						"enum":  []string{"일반", "고급", "희귀", "영웅", "전설", "신화"},
					},
					"level_requirement": map[string]interface{}{
						"type":    "integer",
						"title":   "필요 레벨",
						"minimum": 1,
						"maximum": 100,
					},
					"attack_power": map[string]interface{}{
						"type":    "integer",
						"title":   "공격력",
						"minimum": 0,
					},
					"defense_power": map[string]interface{}{
						"type":    "integer",
						"title":   "방어력",
						"minimum": 0,
					},
					"hp_bonus": map[string]interface{}{
						"type":    "integer",
						"title":   "체력 보너스",
						"minimum": 0,
					},
					"mp_bonus": map[string]interface{}{
						"type":    "integer",
						"title":   "마나 보너스",
						"minimum": 0,
					},
					"special_effect": map[string]interface{}{
						"type":  "string",
						"title": "특수 효과",
					},
					"durability": map[string]interface{}{
						"type":    "integer",
						"title":   "내구도",
						"minimum": 0,
						"maximum": 100,
					},
					"max_stack": map[string]interface{}{
						"type":    "integer",
						"title":   "최대 중첩",
						"minimum": 1,
						"maximum": 999,
					},
					"drop_location": map[string]interface{}{
						"type":  "string",
						"title": "획득 장소",
					},
					"crafting_materials": map[string]interface{}{
						"type":  "string",
						"title": "제작 재료",
					},
				},
				"required": []string{"item_name", "item_type", "rarity"},
			},
		},
	}
}

// GetTemplateByID returns a specific template by ID
func GetTemplateByID(id string) *TableTemplate {
	templates := GetTableTemplates()
	for _, template := range templates {
		if template.ID == id {
			return &template
		}
	}
	return nil
}

// GetTemplatesByCategory returns templates filtered by category
func GetTemplatesByCategory(category string) []TableTemplate {
	templates := GetTableTemplates()
	var filtered []TableTemplate
	
	for _, template := range templates {
		if template.Category == category {
			filtered = append(filtered, template)
		}
	}
	
	return filtered
}