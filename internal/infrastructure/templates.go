package infrastructure

import (
	"encoding/json"
	"log"

	"github.com/jmoiron/sqlx"
)

// Template represents a table template
type Template struct {
	ID          string          `db:"id" json:"id"`
	Name        string          `db:"name" json:"name"`
	Description string          `db:"description" json:"description"`
	Category    string          `db:"category" json:"category"`
	Icon        string          `db:"icon" json:"icon"`
	Schema      json.RawMessage `db:"schema" json:"schema"`
	SampleData  json.RawMessage `db:"sample_data" json:"sample_data,omitempty"`
}

// InitializeTemplates inserts default templates into the database
func InitializeTemplates(db *sqlx.DB) error {
	log.Println("📝 Initializing default templates...")

	templates := getDefaultTemplates()

	for _, tmpl := range templates {
		// Check if template already exists
		var count int
		err := db.Get(&count, "SELECT COUNT(*) FROM templates WHERE id = $1", tmpl.ID)
		if err != nil {
			return err
		}

		if count > 0 {
			log.Printf("⏭️  Template '%s' already exists, skipping", tmpl.Name)
			continue
		}

		// Insert template
		_, err = db.NamedExec(`
			INSERT INTO templates (id, name, description, category, icon, schema, sample_data)
			VALUES (:id, :name, :description, :category, :icon, :schema, :sample_data)
		`, tmpl)

		if err != nil {
			return err
		}

		log.Printf("✅ Added template: %s", tmpl.Name)
	}

	return nil
}

func getDefaultTemplates() []Template {
	return []Template{
		// Business Templates
		{
			ID:          "customer",
			Name:        "고객 관리",
			Description: "고객 정보와 연락처를 관리하는 템플릿",
			Category:    "business",
			Icon:        "👥",
			Schema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"name": { "type": "string", "title": "이름", "minLength": 1 },
					"email": { "type": "string", "format": "email", "title": "이메일" },
					"company": { "type": "string", "title": "회사명" },
					"phone": { "type": "string", "title": "연락처", "pattern": "^[0-9-+()\\s]+$" },
					"interest_level": { "type": "string", "title": "관심도", "enum": ["높음", "중간", "낮음"] },
					"registration_date": { "type": "string", "format": "date", "title": "등록일" }
				},
				"required": ["name", "email"]
			}`),
		},
		{
			ID:          "project",
			Name:        "프로젝트 관리",
			Description: "프로젝트 진행 상황과 일정을 추적하는 템플릿",
			Category:    "business",
			Icon:        "📋",
			Schema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"title": { "type": "string", "title": "제목", "minLength": 1 },
					"description": { "type": "string", "title": "설명" },
					"assignee": { "type": "string", "title": "담당자" },
					"status": { "type": "string", "title": "상태", "enum": ["TODO", "진행중", "완료"] },
					"priority": { "type": "integer", "title": "우선순위", "minimum": 1, "maximum": 5 },
					"due_date": { "type": "string", "format": "date", "title": "마감일" }
				},
				"required": ["title", "status"]
			}`),
		},
		{
			ID:          "inventory",
			Name:        "재고 관리",
			Description: "상품 재고와 공급업체 정보를 관리하는 템플릿",
			Category:    "business",
			Icon:        "📦",
			Schema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"product_name": { "type": "string", "title": "상품명", "minLength": 1 },
					"category": { "type": "string", "title": "카테고리" },
					"quantity": { "type": "integer", "title": "수량", "minimum": 0 },
					"price": { "type": "number", "title": "가격", "minimum": 0 },
					"supplier": { "type": "string", "title": "공급업체" },
					"last_updated": { "type": "string", "format": "date-time", "title": "최종 업데이트" }
				},
				"required": ["product_name", "quantity", "price"]
			}`),
		},
		{
			ID:          "event",
			Name:        "이벤트 관리",
			Description: "이벤트와 일정을 관리하는 템플릿",
			Category:    "business",
			Icon:        "📅",
			Schema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"title": { "type": "string", "title": "제목", "minLength": 1 },
					"date": { "type": "string", "format": "date", "title": "날짜" },
					"time": { "type": "string", "title": "시간", "pattern": "^([01]?[0-9]|2[0-3]):[0-5][0-9]$" },
					"location": { "type": "string", "title": "장소" },
					"attendees": { "type": "integer", "title": "참석자 수", "minimum": 0 },
					"type": { "type": "string", "title": "이벤트 유형", "enum": ["회의", "워크샵", "세미나", "파티", "기타"] }
				},
				"required": ["title", "date", "time"]
			}`),
		},
		// Game Templates
		{
			ID:          "quest",
			Name:        "퀘스트 관리",
			Description: "게임 퀘스트와 보상을 관리하는 템플릿",
			Category:    "game",
			Icon:        "⚔️",
			Schema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"quest_name": { "type": "string", "title": "퀘스트명", "minLength": 1 },
					"description": { "type": "string", "title": "설명" },
					"quest_type": { "type": "string", "title": "퀘스트 유형", "enum": ["메인", "서브", "일일", "주간", "이벤트"] },
					"difficulty": { "type": "string", "title": "난이도", "enum": ["쉬움", "보통", "어려움", "매우어려움"] },
					"level_requirement": { "type": "integer", "title": "필요 레벨", "minimum": 1, "maximum": 100 },
					"reward_exp": { "type": "integer", "title": "보상 경험치", "minimum": 0 },
					"reward_gold": { "type": "integer", "title": "보상 골드", "minimum": 0 },
					"reward_items": { "type": "string", "title": "보상 아이템" },
					"completion_condition": { "type": "string", "title": "완료 조건" },
					"status": { "type": "string", "title": "상태", "enum": ["활성", "비활성", "테스트중"] }
				},
				"required": ["quest_name", "quest_type", "difficulty", "level_requirement"]
			}`),
		},
		{
			ID:          "shop_item",
			Name:        "상품 관리",
			Description: "게임 내 상점 아이템을 관리하는 템플릿",
			Category:    "game",
			Icon:        "🛍️",
			Schema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"item_name": { "type": "string", "title": "상품명", "minLength": 1 },
					"description": { "type": "string", "title": "설명" },
					"category": { "type": "string", "title": "카테고리", "enum": ["무기", "방어구", "소모품", "장식품", "재료", "기타"] },
					"rarity": { "type": "string", "title": "등급", "enum": ["일반", "고급", "희귀", "영웅", "전설"] },
					"price_gold": { "type": "integer", "title": "골드 가격", "minimum": 0 },
					"price_gem": { "type": "integer", "title": "보석 가격", "minimum": 0 },
					"stock": { "type": "integer", "title": "재고", "minimum": -1 },
					"level_requirement": { "type": "integer", "title": "필요 레벨", "minimum": 1, "maximum": 100 },
					"is_limited": { "type": "boolean", "title": "한정 상품" },
					"sale_start_date": { "type": "string", "format": "date", "title": "판매 시작일" },
					"sale_end_date": { "type": "string", "format": "date", "title": "판매 종료일" }
				},
				"required": ["item_name", "category", "rarity"]
			}`),
		},
		{
			ID:          "game_item",
			Name:        "게임 아이템",
			Description: "게임 아이템의 상세 정보를 관리하는 템플릿",
			Category:    "game",
			Icon:        "🗡️",
			Schema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"item_name": { "type": "string", "title": "아이템명", "minLength": 1 },
					"description": { "type": "string", "title": "설명" },
					"item_type": { "type": "string", "title": "아이템 유형", "enum": ["무기", "방어구", "악세서리", "소모품", "재료", "퀘스트", "기타"] },
					"rarity": { "type": "string", "title": "등급", "enum": ["일반", "고급", "희귀", "영웅", "전설", "신화"] },
					"level_requirement": { "type": "integer", "title": "필요 레벨", "minimum": 1, "maximum": 100 },
					"attack_power": { "type": "integer", "title": "공격력", "minimum": 0 },
					"defense_power": { "type": "integer", "title": "방어력", "minimum": 0 },
					"hp_bonus": { "type": "integer", "title": "체력 보너스", "minimum": 0 },
					"mp_bonus": { "type": "integer", "title": "마나 보너스", "minimum": 0 },
					"special_effect": { "type": "string", "title": "특수 효과" },
					"durability": { "type": "integer", "title": "내구도", "minimum": 0, "maximum": 100 },
					"max_stack": { "type": "integer", "title": "최대 중첩", "minimum": 1, "maximum": 999 },
					"drop_location": { "type": "string", "title": "획득 장소" },
					"crafting_materials": { "type": "string", "title": "제작 재료" }
				},
				"required": ["item_name", "item_type", "rarity"]
			}`),
		},
	}
}

// GetTemplates retrieves all templates from the database
func GetTemplates(db *sqlx.DB) ([]Template, error) {
	var templates []Template
	err := db.Select(&templates, "SELECT * FROM templates ORDER BY category, name")
	return templates, err
}

// GetTemplateByID retrieves a template by ID
func GetTemplateByID(db *sqlx.DB, id string) (*Template, error) {
	var template Template
	err := db.Get(&template, "SELECT * FROM templates WHERE id = $1", id)
	return &template, err
}
