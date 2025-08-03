package spreadsheet

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// ComponentShowcasePage - 모든 UI 컴포넌트를 개별적으로 보여주는 쇼케이스 페이지
type ComponentShowcasePage struct {
	app.Compo
}

func (c *ComponentShowcasePage) Render() app.UI {
	return app.Div().Class("showcase-container").Body(
		app.H1().Class("showcase-title").Text("UI Components Showcase"),
		
		// Button Examples Section
		c.renderButtonSection(),
		
		// Badge Examples Section
		c.renderBadgeSection(),
		
		// Card Examples Section
		c.renderCardSection(),
		
		// Input Field Examples Section
		c.renderInputSection(),
		
		// Spreadsheet Cell Examples Section
		c.renderCellSection(),
		
		// Spreadsheet Grid Examples Section
		c.renderGridSection(),
		
		// Sheet Tab Examples Section
		c.renderTabSection(),
		
		// Toolbar Button Examples Section
		c.renderToolbarButtonSection(),
		
		// Status Indicator Examples Section
		c.renderStatusSection(),
		
		// Workbook Data Model Demo Section
		c.renderWorkbookDemo(),
		
		// Database Table Grid Demo Section
		c.renderDatabaseTableDemo(),
	)
}

func (c *ComponentShowcasePage) renderWorkbookDemo() app.UI {
	return app.Div().Class("component-section").Body(
		app.H2().Text("📚 Workbook Data Model (Univer-inspired)"),
		app.Div().Class("component-demo full-width").Body(
			&WorkbookDemo{},
		),
	)
}

func (c *ComponentShowcasePage) renderButtonSection() app.UI {
	return app.Div().Class("component-section").Body(
		app.H2().Text("🔘 Buttons"),
		app.Div().Class("component-grid").Body(
			// Primary buttons
			app.Div().Class("component-demo").Body(
				app.H4().Text("Primary Variants"),
				app.Div().Class("demo-row").Body(
					&UIButton{Text: "Primary", Variant: ButtonPrimary, Size: ButtonMedium},
					&UIButton{Text: "Secondary", Variant: ButtonSecondary, Size: ButtonMedium},
					&UIButton{Text: "Success", Variant: ButtonSuccess, Size: ButtonMedium},
					&UIButton{Text: "Danger", Variant: ButtonDanger, Size: ButtonMedium},
				),
			),
			
			// Button sizes
			app.Div().Class("component-demo").Body(
				app.H4().Text("Sizes"),
				app.Div().Class("demo-row").Body(
					&UIButton{Text: "Small", Variant: ButtonPrimary, Size: ButtonSmall},
					&UIButton{Text: "Medium", Variant: ButtonPrimary, Size: ButtonMedium},
					&UIButton{Text: "Large", Variant: ButtonPrimary, Size: ButtonLarge},
				),
			),
			
			// Button states
			app.Div().Class("component-demo").Body(
				app.H4().Text("States"),
				app.Div().Class("demo-row").Body(
					&UIButton{Text: "Normal", Variant: ButtonPrimary, Size: ButtonMedium},
					&UIButton{Text: "Disabled", Variant: ButtonPrimary, Size: ButtonMedium, Disabled: true},
					&UIButton{Text: "Loading", Variant: ButtonPrimary, Size: ButtonMedium, Loading: true},
					&UIButton{Text: "With Icon", Variant: ButtonPrimary, Size: ButtonMedium, Icon: "★"},
				),
			),
		),
	)
}

func (c *ComponentShowcasePage) renderBadgeSection() app.UI {
	return app.Div().Class("component-section").Body(
		app.H2().Text("🏷️ Badges"),
		app.Div().Class("component-grid").Body(
			// Badge variants
			app.Div().Class("component-demo").Body(
				app.H4().Text("Variants"),
				app.Div().Class("demo-row").Body(
					&UIBadge{Text: "Primary", Variant: BadgePrimary},
					&UIBadge{Text: "Secondary", Variant: BadgeSecondary},
					&UIBadge{Text: "Success", Variant: BadgeSuccess},
					&UIBadge{Text: "Danger", Variant: BadgeDanger},
					&UIBadge{Text: "Warning", Variant: BadgeWarning},
					&UIBadge{Text: "Info", Variant: BadgeInfo},
				),
			),
			
			// Badge with counts
			app.Div().Class("component-demo").Body(
				app.H4().Text("Count Badges"),
				app.Div().Class("demo-row").Body(
					&UIBadge{Variant: BadgeDanger, Count: 5},
					&UIBadge{Variant: BadgeWarning, Count: 23},
					&UIBadge{Variant: BadgeInfo, Count: 99},
					&UIBadge{Variant: BadgeSuccess, Count: 150}, // Shows as 99+
				),
			),
			
			// Dot badges
			app.Div().Class("component-demo").Body(
				app.H4().Text("Dot Indicators"),
				app.Div().Class("demo-row").Body(
					&UIBadge{Variant: BadgePrimary, Dot: true},
					&UIBadge{Variant: BadgeDanger, Dot: true},
					&UIBadge{Variant: BadgeSuccess, Dot: true},
				),
			),
		),
	)
}

func (c *ComponentShowcasePage) renderCardSection() app.UI {
	return app.Div().Class("component-section").Body(
		app.H2().Text("🃏 Cards"),
		app.Div().Class("component-grid").Body(
			// Basic cards
			app.Div().Class("component-demo").Body(
				app.H4().Text("Basic Cards"),
				app.Div().Class("demo-row").Body(
					&UICard{
						Title:       "Simple Card",
						Description: "This is a basic card component with title and description.",
						Bordered:    true,
					},
					&UICard{
						Title:       "Hoverable Card",
						Description: "This card has hover effects enabled.",
						Hoverable:   true,
						Bordered:    true,
					},
				),
			),
			
			// Cards with actions
			app.Div().Class("component-demo").Body(
				app.H4().Text("Cards with Actions"),
				app.Div().Class("demo-row").Body(
					&UICard{
						Title:       "Action Card",
						Description: "This card includes action buttons.",
						Actions: []app.UI{
							&UIButton{Text: "Cancel", Variant: ButtonSecondary, Size: ButtonSmall},
							&UIButton{Text: "Save", Variant: ButtonPrimary, Size: ButtonSmall},
						},
						Bordered:  true,
						Hoverable: true,
					},
				),
			),
		),
	)
}

func (c *ComponentShowcasePage) renderInputSection() app.UI {
	return app.Div().Class("component-section").Body(
		app.H2().Text("📝 Input Fields"),
		app.Div().Class("component-grid").Body(
			// Basic inputs
			app.Div().Class("component-demo").Body(
				app.H4().Text("Basic Inputs"),
				app.Div().Class("demo-column").Body(
					&InputField{
						Label:       "Name",
						Placeholder: "Enter your name",
						Type:        "text",
					},
					&InputField{
						Label:       "Email",
						Placeholder: "Enter your email",
						Type:        "email",
						Required:    true,
					},
					&InputField{
						Label:    "Password",
						Type:     "password",
						Required: true,
					},
				),
			),
			
			// Input states
			app.Div().Class("component-demo").Body(
				app.H4().Text("States"),
				app.Div().Class("demo-column").Body(
					&InputField{
						Label:    "Disabled Field",
						Value:    "Cannot edit this",
						Disabled: true,
					},
					&InputField{
						Label: "Field with Error",
						Value: "invalid@",
						Error: "Please enter a valid email address",
					},
				),
			),
		),
	)
}

func (c *ComponentShowcasePage) renderCellSection() app.UI {
	return app.Div().Class("component-section").Body(
		app.H2().Text("📊 Spreadsheet Cells"),
		app.Div().Class("component-grid").Body(
			// Cell states
			app.Div().Class("component-demo").Body(
				app.H4().Text("Cell States"),
				app.Div().Class("demo-grid").Body(
					&SpreadsheetCell{Value: "Normal Cell"},
					&SpreadsheetCell{Value: "Selected Cell", Selected: true},
					&SpreadsheetCell{Value: "Read Only", ReadOnly: true},
				),
			),
			
			// Formatted cells
			app.Div().Class("component-demo").Body(
				app.H4().Text("Formatted Cells"),
				app.Div().Class("demo-grid").Body(
					&SpreadsheetCell{
						Value: "Bold Text",
						Format: CellFormat{Bold: true},
					},
					&SpreadsheetCell{
						Value: "Italic Text",
						Format: CellFormat{Italic: true},
					},
					&SpreadsheetCell{
						Value: "Colored Text",
						Format: CellFormat{Color: "#ff0000"},
					},
					&SpreadsheetCell{
						Value: "Background",
						Format: CellFormat{BgColor: "#ffff00"},
					},
				),
			),
		),
	)
}

func (c *ComponentShowcasePage) renderGridSection() app.UI {
	return app.Div().Class("component-section").Body(
		app.H2().Text("📊 Spreadsheet Grids"),
		app.Div().Class("component-grid").Body(
			// Small Grid (3x3)
			app.Div().Class("component-demo").Body(
				app.H4().Text("Small Grid (3x3)"),
				c.renderSmallGrid(),
			),
			
			// Medium Grid (5x5)
			app.Div().Class("component-demo").Body(
				app.H4().Text("Medium Grid (5x5)"),
				c.renderMediumGrid(),
			),
			
			// Large Grid (8x6)
			app.Div().Class("component-demo full-width").Body(
				app.H4().Text("Large Grid (8x6)"),
				c.renderLargeGrid(),
			),
			
			// Mobile Scrollable Grid
			app.Div().Class("component-demo").Body(
				app.H4().Text("Mobile Scrollable Grid"),
				c.renderMobileGrid(),
			),
		),
	)
}

func (c *ComponentShowcasePage) renderSmallGrid() app.UI {
	sampleData := map[string]string{
		"0-0": "A1", "0-1": "B1", "0-2": "C1",
		"1-0": "100", "1-1": "200", "1-2": "=A2+B2",
		"2-0": "Item", "2-1": "Price", "2-2": "$300",
	}
	
	return &SpreadsheetGrid{
		Rows: 3,
		Cols: 3,
		Data: sampleData,
		Size: "small",
	}
}

func (c *ComponentShowcasePage) renderMediumGrid() app.UI {
	sampleData := map[string]string{
		"0-0": "Product", "0-1": "Q1", "0-2": "Q2", "0-3": "Q3", "0-4": "Total",
		"1-0": "Widget A", "1-1": "100", "1-2": "120", "1-3": "110", "1-4": "=B2+C2+D2",
		"2-0": "Widget B", "2-1": "80", "2-2": "90", "2-3": "95", "2-4": "=B3+C3+D3",
		"3-0": "Widget C", "3-1": "150", "3-2": "140", "3-3": "160", "3-4": "=B4+C4+D4",
		"4-0": "Total", "4-1": "=B2+B3+B4", "4-2": "=C2+C3+C4", "4-3": "=D2+D3+D4", "4-4": "=B5+C5+D5",
	}
	
	return &SpreadsheetGrid{
		Rows: 5,
		Cols: 5,
		Data: sampleData,
		Size: "medium",
	}
}

func (c *ComponentShowcasePage) renderLargeGrid() app.UI {
	sampleData := map[string]string{
		"0-0": "Item", "0-1": "Jan", "0-2": "Feb", "0-3": "Mar", "0-4": "Apr", "0-5": "May", "0-6": "Jun", "0-7": "Total",
		"1-0": "Sales", "1-1": "1000", "1-2": "1200", "1-3": "1100", "1-4": "1300", "1-5": "1250", "1-6": "1400", "1-7": "=SUM(B2:G2)",
		"2-0": "Costs", "2-1": "800", "2-2": "850", "2-3": "900", "2-4": "920", "2-5": "880", "2-6": "950", "2-7": "=SUM(B3:G3)",
		"3-0": "Profit", "3-1": "=B2-B3", "3-2": "=C2-C3", "3-3": "=D2-D3", "3-4": "=E2-E3", "3-5": "=F2-F3", "3-6": "=G2-G3", "3-7": "=B4+C4+D4+E4+F4+G4",
		"4-0": "Growth %", "4-1": "-", "4-2": "20%", "4-3": "-8.3%", "4-4": "18.2%", "4-5": "-3.8%", "4-6": "12%", "4-7": "Average",
		"5-0": "Target", "5-1": "1000", "5-2": "1100", "5-3": "1150", "5-4": "1200", "5-5": "1250", "5-6": "1300", "5-7": "=SUM(B6:G6)",
	}
	
	return &SpreadsheetGrid{
		Rows: 6,
		Cols: 8,
		Data: sampleData,
		Size: "large",
	}
}

func (c *ComponentShowcasePage) renderMobileGrid() app.UI {
	// Create a larger dataset to demonstrate scrolling
	sampleData := map[string]string{
		// Headers
		"0-0": "ID", "0-1": "Name", "0-2": "Email", "0-3": "Department", "0-4": "Role", "0-5": "Salary", "0-6": "Start Date", "0-7": "Status",
		
		// Employee data (10 rows)
		"1-0": "001", "1-1": "John Doe", "1-2": "john@company.com", "1-3": "Engineering", "1-4": "Senior Dev", "1-5": "$120,000", "1-6": "2020-01-15", "1-7": "Active",
		"2-0": "002", "2-1": "Jane Smith", "2-2": "jane@company.com", "2-3": "Design", "2-4": "UX Designer", "2-5": "$95,000", "2-6": "2021-03-10", "2-7": "Active",
		"3-0": "003", "3-1": "Mike Johnson", "3-2": "mike@company.com", "3-3": "Marketing", "3-4": "Manager", "3-5": "$110,000", "3-6": "2019-08-22", "3-7": "Active",
		"4-0": "004", "4-1": "Sarah Wilson", "4-2": "sarah@company.com", "4-3": "Engineering", "4-4": "Lead Dev", "4-5": "$140,000", "4-6": "2018-05-30", "4-7": "Active",
		"5-0": "005", "5-1": "Tom Brown", "5-2": "tom@company.com", "5-3": "Sales", "5-4": "Sales Rep", "5-5": "$85,000", "5-6": "2022-01-12", "5-7": "Active",
		"6-0": "006", "6-1": "Lisa Davis", "6-2": "lisa@company.com", "6-3": "HR", "6-4": "HR Manager", "6-5": "$105,000", "6-6": "2020-09-15", "6-7": "Active",
		"7-0": "007", "7-1": "Chris Lee", "7-2": "chris@company.com", "7-3": "Engineering", "7-4": "Junior Dev", "7-5": "$75,000", "7-6": "2023-02-01", "7-7": "Active",
		"8-0": "008", "8-1": "Emily Chen", "8-2": "emily@company.com", "8-3": "Design", "8-4": "Designer", "8-5": "$80,000", "8-6": "2022-06-20", "8-7": "Active",
		"9-0": "009", "9-1": "David Kim", "9-2": "david@company.com", "9-3": "Operations", "9-4": "Ops Manager", "9-5": "$115,000", "9-6": "2019-12-03", "9-7": "Active",
		"10-0": "010", "10-1": "Anna Taylor", "10-2": "anna@company.com", "10-3": "Finance", "10-4": "Analyst", "10-5": "$90,000", "10-6": "2021-11-08", "10-7": "Active",
	}
	
	return app.Div().Body(
		app.P().Class("mobile-grid-description").Text("📱 This grid is optimized for mobile with horizontal and vertical scrolling. Try scrolling to see all data!"),
		&SpreadsheetGrid{
			Rows:       11, // 1 header + 10 data rows
			Cols:       8,  // 8 columns
			Data:       sampleData,
			Size:       "medium",
			Scrollable: true,
			MaxHeight:  "250px", // Mobile-friendly height
			MaxWidth:   "100%",  // Full width but with horizontal scroll
		},
	)
}

func (c *ComponentShowcasePage) renderTabSection() app.UI {
	return app.Div().Class("component-section").Body(
		app.H2().Text("📑 Sheet Tabs"),
		app.Div().Class("component-grid").Body(
			// Tab states
			app.Div().Class("component-demo").Body(
				app.H4().Text("Tab States"),
				app.Div().Class("demo-row").Body(
					&SheetTab{Name: "Sheet1", Active: true},
					&SheetTab{Name: "Sheet2"},
					&SheetTab{Name: "Sheet3", Closable: true},
				),
			),
		),
	)
}

func (c *ComponentShowcasePage) renderToolbarButtonSection() app.UI {
	return app.Div().Class("component-section").Body(
		app.H2().Text("🔧 Toolbar Buttons"),
		app.Div().Class("component-grid").Body(
			// Toolbar button states
			app.Div().Class("component-demo").Body(
				app.H4().Text("Toolbar States"),
				app.Div().Class("demo-row").Body(
					&ToolbarButton{Icon: "B", Text: "Bold", Tooltip: "Bold text"},
					&ToolbarButton{Icon: "I", Text: "Italic", Tooltip: "Italic text", Active: true},
					&ToolbarButton{Icon: "U", Text: "Underline", Tooltip: "Underline text"},
					&ToolbarButton{Icon: "S", Text: "Save", Tooltip: "Save document", Disabled: true},
				),
			),
		),
	)
}

func (c *ComponentShowcasePage) renderStatusSection() app.UI {
	return app.Div().Class("component-section").Body(
		app.H2().Text("🔴 Status Indicators"),
		app.Div().Class("component-grid").Body(
			// Status variants
			app.Div().Class("component-demo").Body(
				app.H4().Text("Status Types"),
				app.Div().Class("demo-row").Body(
					&StatusIndicator{Status: "online", Size: "md"},
					&StatusIndicator{Status: "offline", Size: "md"},
					&StatusIndicator{Status: "busy", Size: "md"},
					&StatusIndicator{Status: "away", Size: "md"},
				),
			),
			
			// Status sizes
			app.Div().Class("component-demo").Body(
				app.H4().Text("Sizes"),
				app.Div().Class("demo-row").Body(
					&StatusIndicator{Status: "online", Size: "sm"},
					&StatusIndicator{Status: "online", Size: "md"},
					&StatusIndicator{Status: "online", Size: "lg"},
				),
			),
			
			// Status with dots
			app.Div().Class("component-demo").Body(
				app.H4().Text("With Dot"),
				app.Div().Class("demo-row").Body(
					&StatusIndicator{Status: "online", Size: "md", WithDot: true},
					&StatusIndicator{Status: "busy", Size: "md", WithDot: true},
				),
			),
		),
	)
}

func (c *ComponentShowcasePage) renderDatabaseTableDemo() app.UI {
	return app.Div().Class("component-section").Body(
		app.H2().Text("🗄️ Database Table Grid (PostgreSQL Style)"),
		app.Div().Class("component-grid").Body(
			// Basic Database Table
			app.Div().Class("component-demo full-width").Body(
				app.H4().Text("User Table Example"),
				c.renderUserTableGrid(),
			),
			
			// Product Table with all metadata
			app.Div().Class("component-demo full-width").Body(
				app.H4().Text("Product Table with Full Metadata"),
				c.renderProductTableGrid(),
			),
		),
	)
}

func (c *ComponentShowcasePage) renderUserTableGrid() app.UI {
	// Define columns for a typical user table
	columns := []ColumnMetadata{
		{
			Name:        "id",
			Type:        "integer",
			Nullable:    false,
			Default:     "nextval('users_id_seq')",
			Description: "Primary key",
		},
		{
			Name:        "username",
			Type:        "varchar",
			Length:      50,
			Nullable:    false,
			Description: "Unique username",
		},
		{
			Name:        "email",
			Type:        "varchar",
			Length:      255,
			Nullable:    false,
			Description: "User email address",
		},
		{
			Name:        "active",
			Type:        "boolean",
			Nullable:    false,
			Default:     "true",
			Description: "Account status",
		},
		{
			Name:        "created_at",
			Type:        "timestamp",
			Nullable:    false,
			Default:     "CURRENT_TIMESTAMP",
			Description: "Account creation date",
		},
	}
	
	// Sample data
	sampleData := map[string]string{
		"0-0": "1", "0-1": "john_doe", "0-2": "john@example.com", "0-3": "true", "0-4": "2024-01-15 10:30:00",
		"1-0": "2", "1-1": "jane_smith", "1-2": "jane@example.com", "1-3": "true", "1-4": "2024-01-16 14:20:00",
		"2-0": "3", "2-1": "bob_wilson", "2-2": "bob@example.com", "2-3": "false", "2-4": "2024-01-17 09:15:00",
		"3-0": "4", "3-1": "alice_brown", "3-2": "alice@example.com", "3-3": "true", "3-4": "2024-01-18 16:45:00",
	}
	
	return app.Div().Body(
		app.P().Class("database-demo-description").Text("🔍 This component embeds SpreadsheetGrid to create PostgreSQL-style table views with metadata rows showing field information above the data."),
		&DatabaseTableGrid{
			TableName:       "users",
			Columns:         columns,
			Data:            sampleData,
			ShowDescription: true,
			ShowConstraints: true,
			Scrollable:      true,
			MaxHeight:       "350px",
		},
	)
}

func (c *ComponentShowcasePage) renderProductTableGrid() app.UI {
	// Define columns for a product table with various data types
	columns := []ColumnMetadata{
		{
			Name:        "product_id",
			Type:        "bigint",
			Nullable:    false,
			Default:     "nextval('products_id_seq')",
			Description: "Product identifier",
		},
		{
			Name:        "name",
			Type:        "varchar",
			Length:      100,
			Nullable:    false,
			Description: "Product name",
		},
		{
			Name:        "price",
			Type:        "decimal",
			Precision:   10,
			Scale:       2,
			Nullable:    false,
			Description: "Unit price in USD",
		},
		{
			Name:        "category",
			Type:        "varchar",
			Length:      50,
			Nullable:    true,
			Description: "Product category",
		},
		{
			Name:        "in_stock",
			Type:        "integer",
			Nullable:    false,
			Default:     "0",
			Description: "Current stock quantity",
		},
		{
			Name:        "is_active",
			Type:        "boolean",
			Nullable:    false,
			Default:     "true",
			Description: "Product availability status",
		},
	}
	
	// Sample product data
	sampleData := map[string]string{
		"0-0": "1001", "0-1": "Wireless Headphones", "0-2": "99.99", "0-3": "Electronics", "0-4": "25", "0-5": "true",
		"1-0": "1002", "1-1": "Coffee Mug", "1-2": "12.50", "1-3": "Kitchen", "1-4": "100", "1-5": "true",
		"2-0": "1003", "2-1": "Laptop Stand", "2-2": "45.00", "2-3": "Office", "2-4": "0", "2-5": "false",
		"3-0": "1004", "3-1": "Running Shoes", "3-2": "79.99", "3-3": "Sports", "3-4": "15", "3-5": "true",
		"4-0": "1005", "4-1": "Desk Lamp", "4-2": "28.75", "4-3": "Office", "4-4": "8", "4-5": "true",
	}
	
	return app.Div().Body(
		app.P().Class("database-demo-description").Text("📋 This example shows a complete product table with all metadata types including field names, types, nullability, descriptions, and default values."),
		&DatabaseTableGrid{
			TableName:       "products",
			Columns:         columns,
			Data:            sampleData,
			ShowDescription: true,
			ShowConstraints: true,
			Scrollable:      true,
			MaxHeight:       "400px",
		},
	)
}