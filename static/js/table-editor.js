// Table Editor JavaScript Module
(function() {
    'use strict';

    // Global state
    let tableData = {
        table: null,
        records: [],
        schema: null,
        currentPage: 1,
        isLoading: false,
        hasMore: true,
        editingCell: null
    };

    // Initialize on page load
    document.addEventListener('DOMContentLoaded', initializeTableEditor);

    // Public functions (exposed to window)
    window.openAddRecordModal = openAddRecordModal;
    window.closeAddRecordModal = closeAddRecordModal;
    window.submitNewRecord = submitNewRecord;
    window.openEditCellModal = openEditCellModal;
    window.closeEditCellModal = closeEditCellModal;
    window.saveEditedCell = saveEditedCell;
    window.openExportModal = openExportModal;
    window.closeExportModal = closeExportModal;
    window.exportData = exportData;
    window.openImportModal = openImportModal;
    window.closeImportModal = closeImportModal;
    window.processImport = processImport;
    window.toggleBulkEditMode = toggleBulkEditMode;

    // Initialize table editor
    async function initializeTableEditor() {
        const tableId = extractTableIdFromURL();
        if (!tableId) {
            showError('Table ID not found');
            return;
        }

        tableData.tableId = tableId;
        
        // Setup event listeners
        setupEventListeners();
        
        // Load initial data
        await loadTableData(1, false);
        
        // Setup infinite scroll
        setupInfiniteScroll();
    }

    // Extract table ID from URL
    function extractTableIdFromURL() {
        const pathParts = window.location.pathname.split('/');
        const tableIndex = pathParts.indexOf('table');
        if (tableIndex !== -1 && pathParts[tableIndex + 1]) {
            return pathParts[tableIndex + 1];
        }
        return null;
    }

    // Setup event listeners
    function setupEventListeners() {
        // Search input
        const searchInput = document.getElementById('search-input');
        if (searchInput) {
            searchInput.addEventListener('input', debounce(handleSearch, 300));
        }

        // Modal close on backdrop click
        ['add-record-modal', 'edit-cell-modal', 'export-modal', 'import-modal'].forEach(modalId => {
            const modal = document.getElementById(modalId);
            if (modal) {
                modal.addEventListener('click', function(e) {
                    if (e.target === this) {
                        closeModal(modalId);
                    }
                });
            }
        });

        // Export option buttons
        document.addEventListener('click', function(e) {
            if (e.target.closest('.export-option')) {
                const button = e.target.closest('.export-option');
                const format = button.dataset.format;
                if (format) {
                    exportData(format);
                }
            }
        });

        // Setup floating menu
        setupFloatingMenu();
    }

    // Setup floating menu interactions
    function setupFloatingMenu() {
        const fabMain = document.getElementById('fab-main');
        const fabMenu = document.getElementById('fab-menu');
        let isMenuOpen = false;

        if (fabMain && fabMenu) {
            fabMain.addEventListener('click', function() {
                isMenuOpen = !isMenuOpen;
                
                if (isMenuOpen) {
                    // Show menu
                    fabMenu.classList.remove('opacity-0', 'pointer-events-none');
                    fabMenu.classList.add('opacity-100', 'pointer-events-auto');
                    
                    // Rotate main button
                    fabMain.style.transform = 'rotate(45deg)';
                    
                    // Show labels with delay
                    setTimeout(() => {
                        const labels = fabMenu.querySelectorAll('.fab-label');
                        labels.forEach(label => label.classList.remove('opacity-0'));
                    }, 100);
                } else {
                    // Hide menu
                    fabMenu.classList.remove('opacity-100', 'pointer-events-auto');
                    fabMenu.classList.add('opacity-0', 'pointer-events-none');
                    
                    // Reset main button
                    fabMain.style.transform = 'rotate(0deg)';
                    
                    // Hide labels
                    const labels = fabMenu.querySelectorAll('.fab-label');
                    labels.forEach(label => label.classList.add('opacity-0'));
                }
            });

            // Close menu when clicking outside
            document.addEventListener('click', function(e) {
                if (!e.target.closest('#fab-main') && !e.target.closest('#fab-menu') && isMenuOpen) {
                    fabMain.click(); // Toggle menu closed
                }
            });

            // Handle FAB action clicks
            fabMenu.addEventListener('click', function(e) {
                const actionButton = e.target.closest('.fab-action');
                if (actionButton) {
                    const action = actionButton.dataset.action;
                    if (action) {
                        // Execute the action function
                        try {
                            eval(action);
                        } catch (error) {
                            console.error('Error executing FAB action:', error);
                        }
                        // Close menu after action
                        if (isMenuOpen) {
                            fabMain.click();
                        }
                    }
                }
            });
        }
    }

    // Load table data
    async function loadTableData(page = 1, append = false) {
        if (tableData.isLoading || (!tableData.hasMore && append)) {
            return;
        }

        tableData.isLoading = true;
        showLoadingState(!append);

        try {
            const response = await fetch(`/api/table/${tableData.tableId}?page=${page}&limit=20`);
            if (!response.ok) {
                throw new Error('Failed to fetch table data');
            }

            const data = await response.json();
            
            // Update table info
            if (!append || !tableData.table) {
                tableData.table = data.table;
                tableData.schema = data.table.schema;
                updateTableHeader();
                createGridHeaders();
            }

            // Update records
            if (append) {
                tableData.records = [...tableData.records, ...data.records];
            } else {
                tableData.records = data.records;
            }

            tableData.currentPage = data.pagination.page;
            tableData.hasMore = data.pagination.has_more;

            // Render data
            renderDataRows(data.records, append);
            updateRecordCounts();
            
            hideLoadingState();
            
            if (tableData.records.length === 0) {
                showEmptyState();
            }
        } catch (error) {
            console.error('Error loading table data:', error);
            showError('Failed to load table data');
        } finally {
            tableData.isLoading = false;
        }
    }

    // Setup infinite scroll
    function setupInfiniteScroll() {
        const bodyContainer = document.getElementById('grid-body');
        if (!bodyContainer) return;

        bodyContainer.addEventListener('scroll', () => {
            const { scrollTop, scrollHeight, clientHeight } = bodyContainer;
            
            if (scrollTop + clientHeight >= scrollHeight - 200) {
                if (!tableData.isLoading && tableData.hasMore) {
                    loadTableData(tableData.currentPage + 1, true);
                }
            }
        });
    }

    // Update table header
    function updateTableHeader() {
        const titleElement = document.getElementById('table-title');
        if (titleElement && tableData.table) {
            titleElement.textContent = tableData.table.name || '테이블 편집기';
        }
    }

    // Create grid headers
    function createGridHeaders() {
        const headerContainer = document.getElementById('grid-header');
        if (!headerContainer || !tableData.schema) return;

        const properties = tableData.schema.properties || {};
        const headers = ['ID', ...Object.keys(properties), 'Actions'];
        
        headerContainer.innerHTML = headers.map((header, index) => {
            const isLast = index === headers.length - 1;
            return `
                <div class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider ${!isLast ? 'border-r border-gray-200' : ''}">
                    ${header === 'ID' ? '_ID' : header === 'Actions' ? '작업' : header}
                </div>
            `;
        }).join('');
    }

    // Render data rows
    function renderDataRows(records, append = false) {
        const bodyContainer = document.getElementById('grid-body');
        if (!bodyContainer || !tableData.schema) return;

        const properties = Object.keys(tableData.schema.properties || {});
        
        // Remove loading/empty states if exists
        const loadingState = document.getElementById('loading-state');
        const emptyState = document.getElementById('empty-state');
        if (loadingState) loadingState.remove();
        if (emptyState) emptyState.classList.add('hidden');

        const rowsHTML = records.map(record => {
            const cells = [
                `<div class="px-6 py-4 text-sm text-gray-900 border-r border-gray-200">${record._id || ''}</div>`,
                ...properties.map(prop => {
                    const value = record[prop] || '';
                    const displayValue = formatCellValue(value, tableData.schema.properties[prop]);
                    return `
                        <div class="px-6 py-4 text-sm text-gray-900 border-r border-gray-200 cursor-pointer hover:bg-gray-50"
                             onclick="openEditCellModal('${record._id}', '${prop}', '${escapeHtml(JSON.stringify(value))}')">
                            ${displayValue}
                        </div>
                    `;
                }),
                `<div class="px-6 py-4 text-sm text-gray-900">
                    <button onclick="deleteRecord('${record._id}')" class="text-red-600 hover:text-red-900">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                        </svg>
                    </button>
                </div>`
            ];

            return `
                <div class="grid auto-cols-fr grid-flow-col gap-0 border-b border-gray-200 hover:bg-gray-50">
                    ${cells.join('')}
                </div>
            `;
        }).join('');

        if (append) {
            // Find the last infinite scroll loader if exists and remove it
            const loader = bodyContainer.querySelector('.infinite-scroll-loader');
            if (loader) loader.remove();
            
            bodyContainer.insertAdjacentHTML('beforeend', rowsHTML);
        } else {
            bodyContainer.innerHTML = rowsHTML;
        }

        // Add infinite scroll loader if there's more data
        if (tableData.hasMore) {
            bodyContainer.insertAdjacentHTML('beforeend', `
                <div class="infinite-scroll-loader p-4 text-center">
                    <svg class="mx-auto h-6 w-6 animate-spin text-gray-400" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"/>
                    </svg>
                </div>
            `);
        }
    }

    // Format cell value for display
    function formatCellValue(value, propertySchema) {
        if (value === null || value === undefined) {
            return '<span class="text-gray-400">null</span>';
        }

        if (propertySchema) {
            switch (propertySchema.type) {
                case 'boolean':
                    return value ? '✓' : '✗';
                case 'number':
                case 'integer':
                    return value.toLocaleString();
                case 'string':
                    if (propertySchema.format === 'date') {
                        return new Date(value).toLocaleDateString('ko-KR');
                    }
                    if (propertySchema.format === 'date-time') {
                        return new Date(value).toLocaleString('ko-KR');
                    }
                    return escapeHtml(value);
                case 'array':
                    return `[${value.length} items]`;
                case 'object':
                    return '{...}';
                default:
                    return escapeHtml(String(value));
            }
        }

        return escapeHtml(String(value));
    }

    // Update record counts
    function updateRecordCounts() {
        const totalElement = document.getElementById('total-records');
        const footerTotalElement = document.getElementById('total-records-footer');
        const shownElement = document.getElementById('shown-records');

        if (totalElement) totalElement.textContent = tableData.table?.record_count || 0;
        if (footerTotalElement) footerTotalElement.textContent = tableData.table?.record_count || 0;
        if (shownElement) shownElement.textContent = tableData.records.length;
    }

    // Search handler
    async function handleSearch(event) {
        const searchTerm = event.target.value.trim();
        if (searchTerm.length < 2) {
            await loadTableData(1, false);
            return;
        }

        // Filter records locally for now
        const filtered = tableData.records.filter(record => {
            return Object.values(record).some(value => 
                String(value).toLowerCase().includes(searchTerm.toLowerCase())
            );
        });

        renderFilteredRows(filtered);
    }

    // Render filtered rows
    function renderFilteredRows(records) {
        renderDataRows(records, false);
        document.getElementById('shown-records').textContent = records.length;
    }

    // Modal functions
    function openAddRecordModal() {
        if (!tableData.schema) return;

        const form = document.getElementById('add-record-form');
        if (!form) return;

        // Generate form fields based on schema
        const properties = tableData.schema.properties || {};
        const required = tableData.schema.required || [];

        form.innerHTML = Object.entries(properties).map(([key, prop]) => {
            const isRequired = required.includes(key);
            return createFormField(key, prop, isRequired);
        }).join('');

        document.getElementById('add-record-modal').classList.remove('hidden');
    }

    function closeAddRecordModal() {
        document.getElementById('add-record-modal').classList.add('hidden');
    }

    async function submitNewRecord() {
        const form = document.getElementById('add-record-form');
        if (!form) return;

        const formData = new FormData(form);
        const record = {};

        for (const [key, value] of formData.entries()) {
            const prop = tableData.schema.properties[key];
            if (prop) {
                record[key] = parseFormValue(value, prop.type);
            }
        }

        try {
            const response = await fetch(`/api/table/${tableData.tableId}/record`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(record)
            });

            if (!response.ok) throw new Error('Failed to add record');

            closeAddRecordModal();
            await loadTableData(1, false);
            showSuccess('레코드가 추가되었습니다');
        } catch (error) {
            console.error('Error adding record:', error);
            showError('레코드 추가에 실패했습니다');
        }
    }

    function openEditCellModal(recordId, field, value) {
        tableData.editingCell = { recordId, field, value: JSON.parse(value) };

        const label = document.getElementById('edit-cell-label');
        const container = document.getElementById('edit-cell-input-container');
        
        if (label) label.textContent = `${field} 편집`;
        
        if (container && tableData.schema) {
            const prop = tableData.schema.properties[field];
            container.innerHTML = createEditInput(field, prop, tableData.editingCell.value);
        }

        document.getElementById('edit-cell-modal').classList.remove('hidden');
    }

    function closeEditCellModal() {
        document.getElementById('edit-cell-modal').classList.add('hidden');
        tableData.editingCell = null;
    }

    async function saveEditedCell() {
        if (!tableData.editingCell) return;

        const input = document.querySelector('#edit-cell-input-container input, #edit-cell-input-container textarea, #edit-cell-input-container select');
        if (!input) return;

        const prop = tableData.schema.properties[tableData.editingCell.field];
        const newValue = parseFormValue(input.value, prop.type);

        try {
            const response = await fetch(`/api/table/${tableData.tableId}/record/${tableData.editingCell.recordId}`, {
                method: 'PATCH',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ [tableData.editingCell.field]: newValue })
            });

            if (!response.ok) throw new Error('Failed to update cell');

            closeEditCellModal();
            
            // Update local data
            const record = tableData.records.find(r => r._id === tableData.editingCell.recordId);
            if (record) {
                record[tableData.editingCell.field] = newValue;
                renderDataRows(tableData.records, false);
            }

            showSuccess('셀이 업데이트되었습니다');
        } catch (error) {
            console.error('Error updating cell:', error);
            showError('셀 업데이트에 실패했습니다');
        }
    }

    function openExportModal() {
        document.getElementById('export-modal').classList.remove('hidden');
    }

    function closeExportModal() {
        document.getElementById('export-modal').classList.add('hidden');
    }

    async function exportData(format) {
        try {
            const response = await fetch(`/api/table/${tableData.tableId}/export?format=${format}`);
            if (!response.ok) throw new Error('Export failed');

            const blob = await response.blob();
            const url = URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = `${tableData.table.name}.${format}`;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            URL.revokeObjectURL(url);

            closeExportModal();
            showSuccess(`${format.toUpperCase()} 형식으로 내보내기 완료`);
        } catch (error) {
            console.error('Error exporting data:', error);
            showError('데이터 내보내기에 실패했습니다');
        }
    }

    // Helper functions
    function createFormField(key, prop, isRequired) {
        let input = '';
        const baseClasses = 'mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm';

        switch (prop.type) {
            case 'boolean':
                input = `
                    <input type="checkbox" name="${key}" id="${key}" class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded">
                `;
                break;
            case 'integer':
            case 'number':
                input = `<input type="number" name="${key}" id="${key}" ${isRequired ? 'required' : ''} class="${baseClasses}">`;
                break;
            case 'string':
                if (prop.enum) {
                    input = `
                        <select name="${key}" id="${key}" ${isRequired ? 'required' : ''} class="${baseClasses}">
                            <option value="">선택하세요</option>
                            ${prop.enum.map(opt => `<option value="${opt}">${opt}</option>`).join('')}
                        </select>
                    `;
                } else if (prop.format === 'date') {
                    input = `<input type="date" name="${key}" id="${key}" ${isRequired ? 'required' : ''} class="${baseClasses}">`;
                } else if (prop.format === 'date-time') {
                    input = `<input type="datetime-local" name="${key}" id="${key}" ${isRequired ? 'required' : ''} class="${baseClasses}">`;
                } else if (prop.format === 'email') {
                    input = `<input type="email" name="${key}" id="${key}" ${isRequired ? 'required' : ''} class="${baseClasses}">`;
                } else {
                    input = `<input type="text" name="${key}" id="${key}" ${isRequired ? 'required' : ''} class="${baseClasses}">`;
                }
                break;
            default:
                input = `<textarea name="${key}" id="${key}" rows="3" ${isRequired ? 'required' : ''} class="${baseClasses}"></textarea>`;
        }

        return `
            <div class="mb-4">
                <label for="${key}" class="block text-sm font-medium text-gray-700">
                    ${prop.title || key}
                    ${isRequired ? '<span class="text-red-500">*</span>' : ''}
                </label>
                ${input}
                ${prop.description ? `<p class="mt-1 text-sm text-gray-500">${prop.description}</p>` : ''}
            </div>
        `;
    }

    function createEditInput(key, prop, value) {
        const baseClasses = 'mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm';

        switch (prop.type) {
            case 'boolean':
                return `<input type="checkbox" ${value ? 'checked' : ''} class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded">`;
            case 'integer':
            case 'number':
                return `<input type="number" value="${value || ''}" class="${baseClasses}">`;
            case 'string':
                if (prop.enum) {
                    return `
                        <select class="${baseClasses}">
                            ${prop.enum.map(opt => `<option value="${opt}" ${value === opt ? 'selected' : ''}>${opt}</option>`).join('')}
                        </select>
                    `;
                } else if (prop.format === 'date' || prop.format === 'date-time') {
                    return `<input type="${prop.format === 'date' ? 'date' : 'datetime-local'}" value="${value || ''}" class="${baseClasses}">`;
                } else {
                    return `<input type="text" value="${escapeHtml(value || '')}" class="${baseClasses}">`;
                }
            default:
                return `<textarea rows="3" class="${baseClasses}">${escapeHtml(value || '')}</textarea>`;
        }
    }

    function parseFormValue(value, type) {
        switch (type) {
            case 'boolean':
                return value === 'on' || value === true;
            case 'integer':
                return parseInt(value, 10) || 0;
            case 'number':
                return parseFloat(value) || 0;
            default:
                return value;
        }
    }

    function closeModal(modalId) {
        document.getElementById(modalId).classList.add('hidden');
    }

    function showLoadingState(show = true) {
        const loadingState = document.getElementById('loading-state');
        if (loadingState) {
            loadingState.style.display = show ? 'block' : 'none';
        }
    }

    function hideLoadingState() {
        showLoadingState(false);
    }

    function showEmptyState() {
        const emptyState = document.getElementById('empty-state');
        if (emptyState) {
            emptyState.classList.remove('hidden');
        }
    }

    function showError(message) {
        console.error(message);
        // TODO: Implement toast notification
    }

    function showSuccess(message) {
        console.log(message);
        // TODO: Implement toast notification
    }

    function escapeHtml(str) {
        const div = document.createElement('div');
        div.textContent = str;
        return div.innerHTML;
    }

    function debounce(func, wait) {
        let timeout;
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    }

    // Import modal functions
    let importData = {
        file: null,
        parsedData: null,
        mode: 'replace'
    };

    function openImportModal() {
        document.getElementById('import-modal').classList.remove('hidden');
        setupImportHandlers();
        resetImportModal();
    }

    function closeImportModal() {
        document.getElementById('import-modal').classList.add('hidden');
        resetImportModal();
    }

    function resetImportModal() {
        importData.file = null;
        importData.parsedData = null;
        importData.mode = 'replace';
        
        // Reset file input
        document.getElementById('import-file-input').value = '';
        document.getElementById('selected-file-info').classList.add('hidden');
        document.getElementById('import-preview').classList.add('hidden');
        document.getElementById('import-error').classList.add('hidden');
        document.getElementById('import-submit-btn').disabled = true;
        
        // Reset radio buttons
        document.querySelector('input[name="import-mode"][value="replace"]').checked = true;
    }

    function setupImportHandlers() {
        const fileInput = document.getElementById('import-file-input');
        const dropZone = document.getElementById('file-drop-zone');
        const removeBtn = document.getElementById('remove-file-btn');

        // File input change
        fileInput.addEventListener('change', handleFileSelect);
        
        // Drop zone click
        dropZone.addEventListener('click', () => fileInput.click());
        
        // Drag and drop
        dropZone.addEventListener('dragover', (e) => {
            e.preventDefault();
            dropZone.classList.add('border-blue-400', 'bg-blue-50');
        });
        
        dropZone.addEventListener('dragleave', (e) => {
            e.preventDefault();
            dropZone.classList.remove('border-blue-400', 'bg-blue-50');
        });
        
        dropZone.addEventListener('drop', (e) => {
            e.preventDefault();
            dropZone.classList.remove('border-blue-400', 'bg-blue-50');
            
            const files = e.dataTransfer.files;
            if (files.length > 0) {
                handleFileSelect({ target: { files } });
            }
        });
        
        // Remove file button
        if (removeBtn) {
            removeBtn.addEventListener('click', () => {
                resetImportModal();
            });
        }

        // Import mode change
        document.querySelectorAll('input[name="import-mode"]').forEach(radio => {
            radio.addEventListener('change', (e) => {
                importData.mode = e.target.value;
            });
        });
    }

    function handleFileSelect(event) {
        const file = event.target.files[0];
        if (!file) return;

        // Validate file
        const maxSize = 10 * 1024 * 1024; // 10MB
        if (file.size > maxSize) {
            showImportError('파일 크기가 10MB를 초과합니다.');
            return;
        }

        const validTypes = ['text/csv', 'application/json', '.csv', '.json'];
        const fileExtension = '.' + file.name.split('.').pop().toLowerCase();
        if (!validTypes.includes(file.type) && !validTypes.includes(fileExtension)) {
            showImportError('CSV 또는 JSON 파일만 지원됩니다.');
            return;
        }

        importData.file = file;
        showSelectedFile(file);
        parseFile(file);
    }

    function showSelectedFile(file) {
        document.getElementById('selected-file-name').textContent = file.name;
        document.getElementById('selected-file-size').textContent = formatFileSize(file.size);
        document.getElementById('selected-file-info').classList.remove('hidden');
        document.getElementById('import-error').classList.add('hidden');
    }

    function parseFile(file) {
        const reader = new FileReader();
        reader.onload = function(e) {
            const content = e.target.result;
            
            try {
                if (file.name.toLowerCase().endsWith('.json')) {
                    parseJSONFile(content);
                } else if (file.name.toLowerCase().endsWith('.csv')) {
                    parseCSVFile(content);
                }
            } catch (error) {
                showImportError('파일 파싱 오류: ' + error.message);
            }
        };
        
        reader.onerror = function() {
            showImportError('파일을 읽을 수 없습니다.');
        };
        
        reader.readAsText(file);
    }

    function parseJSONFile(content) {
        try {
            const data = JSON.parse(content);
            
            if (!Array.isArray(data)) {
                throw new Error('JSON 파일은 배열 형태여야 합니다.');
            }
            
            if (data.length === 0) {
                throw new Error('빈 데이터입니다.');
            }
            
            // Validate against schema if available
            if (tableData.schema) {
                validateDataAgainstSchema(data);
            }
            
            importData.parsedData = data;
            showImportPreview(data);
            
        } catch (error) {
            showImportError('JSON 파싱 오류: ' + error.message);
        }
    }

    function parseCSVFile(content) {
        try {
            const lines = content.trim().split('\n');
            if (lines.length < 2) {
                throw new Error('CSV 파일에는 최소한 헤더와 한 줄의 데이터가 있어야 합니다.');
            }
            
            // Parse header
            const headers = parseCSVLine(lines[0]);
            
            // Parse data rows
            const data = [];
            for (let i = 1; i < lines.length; i++) {
                if (lines[i].trim()) {
                    const values = parseCSVLine(lines[i]);
                    if (values.length === headers.length) {
                        const row = {};
                        headers.forEach((header, index) => {
                            row[header.trim()] = values[index]?.trim() || '';
                        });
                        data.push(row);
                    }
                }
            }
            
            if (data.length === 0) {
                throw new Error('유효한 데이터 행이 없습니다.');
            }
            
            // Validate against schema if available
            if (tableData.schema) {
                validateDataAgainstSchema(data);
            }
            
            importData.parsedData = data;
            showImportPreview(data);
            
        } catch (error) {
            showImportError('CSV 파싱 오류: ' + error.message);
        }
    }

    function parseCSVLine(line) {
        const result = [];
        let current = '';
        let inQuotes = false;
        
        for (let i = 0; i < line.length; i++) {
            const char = line[i];
            
            if (char === '"') {
                inQuotes = !inQuotes;
            } else if (char === ',' && !inQuotes) {
                result.push(current);
                current = '';
            } else {
                current += char;
            }
        }
        
        result.push(current);
        return result;
    }

    function validateDataAgainstSchema(data) {
        const properties = tableData.schema.properties || {};
        const required = tableData.schema.required || [];
        
        // Check if all required fields are present in at least one row
        const sampleRow = data[0];
        for (const field of required) {
            if (!(field in sampleRow)) {
                throw new Error(`필수 필드 '${field}'가 누락되었습니다.`);
            }
        }
        
        // Type validation can be added here if needed
    }

    function showImportPreview(data) {
        const preview = data.slice(0, 5); // Show first 5 rows
        const container = document.getElementById('import-preview-content');
        
        if (preview.length === 0) return;
        
        // Create table for preview
        const headers = Object.keys(preview[0]);
        let html = `
            <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                    <tr>
                        ${headers.map(header => `<th class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase">${escapeHtml(header)}</th>`).join('')}
                    </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                    ${preview.map(row => `
                        <tr>
                            ${headers.map(header => `<td class="px-3 py-2 text-sm text-gray-900">${escapeHtml(String(row[header] || ''))}</td>`).join('')}
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;
        
        container.innerHTML = html;
        document.getElementById('import-total-rows').textContent = data.length;
        document.getElementById('import-preview').classList.remove('hidden');
        document.getElementById('import-submit-btn').disabled = false;
    }

    function showImportError(message) {
        document.getElementById('import-error-message').textContent = message;
        document.getElementById('import-error').classList.remove('hidden');
        document.getElementById('import-preview').classList.add('hidden');
        document.getElementById('import-submit-btn').disabled = true;
    }

    async function processImport() {
        if (!importData.parsedData) {
            showImportError('가져올 데이터가 없습니다.');
            return;
        }

        const submitBtn = document.getElementById('import-submit-btn');
        const originalText = submitBtn.textContent;
        
        try {
            submitBtn.disabled = true;
            submitBtn.textContent = '처리 중...';

            const response = await fetch(`/api/table/${tableData.tableId}/import`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    data: importData.parsedData,
                    mode: importData.mode
                })
            });

            if (!response.ok) {
                const error = await response.text();
                throw new Error(error || 'Import failed');
            }

            const result = await response.json();
            
            closeImportModal();
            showSuccess(`성공적으로 ${result.imported || importData.parsedData.length}개의 레코드를 가져왔습니다.`);
            
            // Reload table data
            await loadTableData(1, false);
            
        } catch (error) {
            console.error('Import error:', error);
            showImportError('가져오기 실패: ' + error.message);
        } finally {
            submitBtn.disabled = false;
            submitBtn.textContent = originalText;
        }
    }

    function formatFileSize(bytes) {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }

    // Bulk edit functions
    function toggleBulkEditMode() {
        // TODO: Implement bulk edit mode
        showSuccess('일괄 편집 모드는 곧 추가될 예정입니다');
    }

})();