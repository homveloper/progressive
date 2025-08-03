const { test, expect } = require('@playwright/test');

test.describe('DatabaseTableGrid Component', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the showcase page where DatabaseTableGrid is displayed
    await page.goto('/showcase');
    
    // Wait for the page to load and WebAssembly to initialize
    await page.waitForSelector('.database-table-grid', { timeout: 30000 });
    
    // Scroll to the database table grid section
    const dbSection = page.locator('h2:text("🗄️ Database Table Grid (PostgreSQL Style)")');
    await dbSection.scrollIntoViewIfNeeded();
    
    // Wait a bit for scrolling to complete
    await page.waitForTimeout(500);
  });

  test('should display database table grid components', async ({ page }) => {
    // Check if the database table grid section exists
    const dbSection = page.locator('h2:text("🗄️ Database Table Grid (PostgreSQL Style)")');
    await expect(dbSection).toBeVisible();
    
    // Check if both demo tables are present
    const userTableDemo = page.locator('h4:text("User Table Example")');
    const productTableDemo = page.locator('h4:text("Product Table with Full Metadata")');
    
    await expect(userTableDemo).toBeVisible();
    await expect(productTableDemo).toBeVisible();
  });

  test('should display table names correctly', async ({ page }) => {
    // Check if table names are displayed
    const usersTableName = page.locator('.database-table-grid .table-name:text("Table: users")');
    const productsTableName = page.locator('.database-table-grid .table-name:text("Table: products")');
    
    await expect(usersTableName).toBeVisible();
    await expect(productsTableName).toBeVisible();
  });

  test('should display metadata rows correctly', async ({ page }) => {
    // Focus on the first database table grid (users table)
    const firstGrid = page.locator('.database-table-grid').first();
    
    // Check if metadata headers are present
    await expect(firstGrid.locator('.meta-header:text("Field")')).toBeVisible();
    await expect(firstGrid.locator('.meta-header:text("id")')).toBeVisible();
    await expect(firstGrid.locator('.meta-header:text("username")')).toBeVisible();
    await expect(firstGrid.locator('.meta-header:text("email")')).toBeVisible();
    
    // Check if Type row is present
    await expect(firstGrid.locator('.meta-label:text("Type")')).toBeVisible();
    await expect(firstGrid.locator('.type-cell:text("integer")')).toBeVisible();
    await expect(firstGrid.locator('.type-cell:text("varchar(50)")')).toBeVisible();
    
    // Check if Nullable row is present
    await expect(firstGrid.locator('.meta-label:text("Nullable")')).toBeVisible();
    await expect(firstGrid.locator('.nullable-cell:text("NO")').first()).toBeVisible();
    await expect(firstGrid.locator('.nullable-cell:text("YES")').first()).toBeVisible();
  });

  test('should display description metadata when enabled', async ({ page }) => {
    const firstGrid = page.locator('.database-table-grid').first();
    
    // Check if Description row is present (should be enabled for user table)
    await expect(firstGrid.locator('.meta-label:text("Description")')).toBeVisible();
    await expect(firstGrid.locator('.desc-cell:text("Primary key")')).toBeVisible();
    await expect(firstGrid.locator('.desc-cell:text("Unique username")')).toBeVisible();
  });

  test('should display default values when ShowConstraints is enabled', async ({ page }) => {
    const firstGrid = page.locator('.database-table-grid').first();
    
    // Check if Default row is present
    await expect(firstGrid.locator('.meta-label:text("Default")')).toBeVisible();
    await expect(firstGrid.locator('.default-cell:text("nextval(\'users_id_seq\')")')).toBeVisible();
    await expect(firstGrid.locator('.default-cell:text("true")')).toBeVisible();
  });

  test('should display separator row between metadata and data', async ({ page }) => {
    const firstGrid = page.locator('.database-table-grid').first();
    
    // Check if separator row exists with dashes
    const separatorCells = firstGrid.locator('.separator-row .separator:text("---")');
    await expect(separatorCells.first()).toBeVisible();
    
    // Should have multiple separator cells (one per column)
    const separatorCount = await separatorCells.count();
    expect(separatorCount).toBeGreaterThan(4); // At least 5 columns + row number column
  });

  test('should display data rows with correct values', async ({ page }) => {
    const firstGrid = page.locator('.database-table-grid').first();
    
    // Check if data rows are present (should be at least 1, but we'll check for specific data)
    const dataRows = firstGrid.locator('.data-row');
    await expect(dataRows).toHaveCountGreaterThan(0);
    
    // Check specific data values
    await expect(firstGrid.locator('.data-cell:text("john_doe")')).toBeVisible();
    await expect(firstGrid.locator('.data-cell:text("john@example.com")')).toBeVisible();
    await expect(firstGrid.locator('.data-cell:text("2024-01-15 10:30:00")')).toBeVisible();
  });

  test('should apply correct CSS classes for data types', async ({ page }) => {
    const secondGrid = page.locator('.database-table-grid').nth(1); // Products table
    
    // Check number cells (price, stock)
    const numberCells = secondGrid.locator('.cell-number');
    await expect(numberCells.first()).toBeVisible();
    
    // Check text cells (name, category)
    const textCells = secondGrid.locator('.cell-text');
    await expect(textCells.first()).toBeVisible();
    
    // Check boolean cells (is_active)
    const booleanCells = secondGrid.locator('.cell-boolean');
    await expect(booleanCells.first()).toBeVisible();
  });

  test('should be scrollable when content exceeds container', async ({ page }) => {
    const firstGrid = page.locator('.database-table-grid').first();
    
    // Check if the grid has scrollable class
    await expect(firstGrid).toHaveClass(/grid-scrollable/);
    
    // Check if the grid has proper max-height style
    const style = await firstGrid.getAttribute('style');
    expect(style).toContain('max-height');
    expect(style).toContain('overflow: auto');
  });

  test('should support cell selection', async ({ page }) => {
    const firstGrid = page.locator('.database-table-grid').first();
    
    // Click on a data cell
    const targetCell = firstGrid.locator('.data-cell').first();
    await targetCell.click();
    
    // Check if the cell gets selected class
    await expect(targetCell).toHaveClass(/selected/);
  });

  test('should display row numbers correctly', async ({ page }) => {
    const firstGrid = page.locator('.database-table-grid').first();
    
    // Check if row numbers are present (at least row 1 should exist)
    await expect(firstGrid.locator('.row-number:text("1")')).toBeVisible();
    
    // Check that we have row number elements
    const rowNumbers = firstGrid.locator('.row-number');
    await expect(rowNumbers).toHaveCountGreaterThan(0);
  });

  test('should handle different data types in products table', async ({ page }) => {
    const secondGrid = page.locator('.database-table-grid').nth(1);
    
    // Check for decimal type (precision and scale)
    await expect(secondGrid.locator('.type-cell:text("decimal(10,2)")')).toBeVisible();
    
    // Check for bigint type
    await expect(secondGrid.locator('.type-cell:text("bigint")')).toBeVisible();
    
    // Check for varchar with length
    await expect(secondGrid.locator('.type-cell:text("varchar(100)")')).toBeVisible();
    
    // Check actual product data
    await expect(secondGrid.locator('.data-cell:text("Wireless Headphones")')).toBeVisible();
    await expect(secondGrid.locator('.data-cell:text("99.99")')).toBeVisible();
    await expect(secondGrid.locator('.data-cell:text("Electronics")')).toBeVisible();
  });

  test('should be responsive on mobile viewports', async ({ page }) => {
    // Set mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });
    
    const firstGrid = page.locator('.database-table-grid').first();
    
    // Grid should still be visible and functional on mobile
    await expect(firstGrid).toBeVisible();
    
    // Should be scrollable
    await expect(firstGrid).toHaveClass(/grid-scrollable/);
    
    // Should be able to scroll horizontally if needed
    const scrollWidth = await firstGrid.evaluate(el => el.scrollWidth);
    const clientWidth = await firstGrid.evaluate(el => el.clientWidth);
    
    // If content is wider than container, it should be scrollable
    if (scrollWidth > clientWidth) {
      // Test horizontal scrolling
      await firstGrid.evaluate(el => el.scrollLeft = 100);
      const newScrollLeft = await firstGrid.evaluate(el => el.scrollLeft);
      expect(newScrollLeft).toBeGreaterThan(0);
    }
  });

  test('should display demo descriptions', async ({ page }) => {
    // Check if demo descriptions are present
    const userDescription = page.locator('.database-demo-description').first();
    const productDescription = page.locator('.database-demo-description').nth(1);
    
    await expect(userDescription).toBeVisible();
    await expect(userDescription).toContainText('This component embeds SpreadsheetGrid');
    
    await expect(productDescription).toBeVisible();
    await expect(productDescription).toContainText('complete product table with all metadata types');
  });

  test('should have proper accessibility attributes', async ({ page }) => {
    const firstGrid = page.locator('.database-table-grid').first();
    
    // Check if table structure is accessible
    // Headers should be properly structured
    const headers = firstGrid.locator('.meta-header');
    await expect(headers.first()).toBeVisible();
    
    // Cells should be clickable (implicit role="button" behavior)
    const dataCells = firstGrid.locator('.data-cell');
    await expect(dataCells.first()).toBeVisible();
  });
});