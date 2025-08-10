// Fakeit Data Generator Module
(function() {
    'use strict';

    // Global state
    let fakeitData = {
        schema: null,
        fieldConfigs: {},
        generatedData: null
    };

    // Smart type mapping for auto-recommendations
    const smartTypeMapping = {
        // Name variations
        'name': 'name',
        'fullname': 'name',
        'full_name': 'name',
        'firstname': 'firstName',
        'first_name': 'firstName',
        'lastname': 'lastName',
        'last_name': 'lastName',
        'username': 'username',
        'user_name': 'username',
        
        // Contact information
        'email': 'email',
        'mail': 'email',
        'email_address': 'email',
        'phone': 'phone',
        'tel': 'phone',
        'telephone': 'phone',
        'mobile': 'phone',
        'phone_number': 'phone',
        
        // Location
        'address': 'address',
        'street': 'address',
        'city': 'city',
        'country': 'country',
        'nation': 'country',
        'latitude': 'latitude',
        'lat': 'latitude',
        'longitude': 'longitude',
        'lng': 'longitude',
        'lon': 'longitude',
        
        // Business
        'company': 'company',
        'company_name': 'company',
        'corporation': 'company',
        'job': 'jobTitle',
        'job_title': 'jobTitle',
        'position': 'jobTitle',
        'title': 'jobTitle',
        
        // Web & Tech
        'url': 'url',
        'website': 'url',
        'link': 'url',
        'password': 'password',
        'pwd': 'password',
        'pass': 'password',
        'uuid': 'uuid',
        'id': 'uuid',
        'guid': 'uuid',
        
        // Text content
        'description': 'paragraph',
        'desc': 'paragraph',
        'content': 'paragraph',
        'text': 'paragraph',
        'comment': 'sentence',
        'note': 'sentence',
        'message': 'sentence',
        
        // Numbers
        'age': 'age',
        'year': 'year',
        'birth_year': 'year',
        'price': 'price',
        'cost': 'price',
        'amount': 'price',
        'total': 'price',
        'quantity': 'quantity',
        'qty': 'quantity',
        'count': 'quantity',
        'stock': 'quantity',
        'rating': 'rating',
        'score': 'rating',
        'stars': 'rating',
        
        // Date/Time
        'created': 'year',
        'updated': 'year',
        'created_at': 'year',
        'updated_at': 'year',
        'date': 'year',
        'month': 'month',
        'day': 'day',
        
        // Boolean indicators
        'active': 'weighted',
        'is_active': 'weighted',
        'enabled': 'weighted',
        'is_enabled': 'weighted',
        'verified': 'weighted',
        'is_verified': 'weighted',
        'published': 'weighted',
        'is_published': 'weighted',
        'available': 'weighted',
        'is_available': 'weighted',
        'in_stock': 'weighted',
        'instock': 'weighted',
        
        // Color
        'color': 'color',
        'colour': 'color',
        'hex': 'color',
        'hex_color': 'color'
    };

    // Fake data type options
    const fakeDataTypes = {
        string: [
            { value: 'name', label: '이름', example: 'John Doe' },
            { value: 'firstName', label: '이름(성)', example: 'John' },
            { value: 'lastName', label: '이름(성)', example: 'Doe' },
            { value: 'email', label: '이메일', example: 'john@example.com' },
            { value: 'phone', label: '전화번호', example: '010-1234-5678' },
            { value: 'address', label: '주소', example: '서울시 강남구...' },
            { value: 'company', label: '회사명', example: '삼성전자' },
            { value: 'jobTitle', label: '직책', example: '소프트웨어 엔지니어' },
            { value: 'city', label: '도시', example: '서울' },
            { value: 'country', label: '국가', example: '대한민국' },
            { value: 'lorem', label: 'Lorem 텍스트', example: 'Lorem ipsum dolor...' },
            { value: 'sentence', label: '문장', example: 'This is a sentence.' },
            { value: 'paragraph', label: '문단', example: 'This is a paragraph...' },
            { value: 'uuid', label: 'UUID', example: '550e8400-e29b-41d4-a716-446655440000' },
            { value: 'url', label: 'URL', example: 'https://example.com' },
            { value: 'username', label: '사용자명', example: 'johndoe123' },
            { value: 'password', label: '비밀번호', example: 'P@ssw0rd123' },
            { value: 'color', label: '색상', example: '#FF5733' },
            { value: 'custom', label: '커스텀', example: '직접 입력' }
        ],
        integer: [
            { value: 'number', label: '숫자', example: '42' },
            { value: 'age', label: '나이', example: '25' },
            { value: 'year', label: '연도', example: '2023' },
            { value: 'month', label: '월', example: '12' },
            { value: 'day', label: '일', example: '25' },
            { value: 'price', label: '가격', example: '10000' },
            { value: 'quantity', label: '수량', example: '5' },
            { value: 'rating', label: '평점', example: '4' },
            { value: 'custom', label: '커스텀 범위', example: '1-100' }
        ],
        number: [
            { value: 'float', label: '소수', example: '3.14' },
            { value: 'price', label: '가격', example: '99.99' },
            { value: 'latitude', label: '위도', example: '37.5665' },
            { value: 'longitude', label: '경도', example: '126.9780' },
            { value: 'percentage', label: '퍼센트', example: '85.5' },
            { value: 'custom', label: '커스텀 범위', example: '0.0-100.0' }
        ],
        boolean: [
            { value: 'boolean', label: '불린', example: 'true/false' },
            { value: 'weighted', label: '가중치 불린', example: '70% true' }
        ]
    };

    // Initialize on page load
    document.addEventListener('DOMContentLoaded', initializeFakeit);

    function initializeFakeit() {
        setupEventListeners();
    }

    function setupEventListeners() {
        // Schema file upload
        const fileInput = document.getElementById('schema-file-input');
        const dropZone = document.getElementById('schema-drop-zone');
        const removeBtn = document.getElementById('remove-schema-btn');

        fileInput.addEventListener('change', handleSchemaFileSelect);
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
                handleSchemaFileSelect({ target: { files } });
            }
        });

        // Remove schema button
        if (removeBtn) {
            removeBtn.addEventListener('click', clearSchema);
        }

        // Action buttons
        document.getElementById('clear-all-btn').addEventListener('click', clearAll);
        document.getElementById('generate-data-btn').addEventListener('click', generateData);
        document.getElementById('export-json-btn').addEventListener('click', () => exportData('json'));
        document.getElementById('export-csv-btn').addEventListener('click', () => exportData('csv'));
    }

    function handleSchemaFileSelect(event) {
        const file = event.target.files[0];
        if (!file) return;

        // Validate file
        const maxSize = 5 * 1024 * 1024; // 5MB
        if (file.size > maxSize) {
            showSchemaError('파일 크기가 5MB를 초과합니다.');
            return;
        }

        if (!file.name.toLowerCase().endsWith('.json')) {
            showSchemaError('JSON 파일만 지원됩니다.');
            return;
        }

        showSelectedSchemaFile(file);
        parseSchemaFile(file);
    }

    function showSelectedSchemaFile(file) {
        document.getElementById('selected-schema-name').textContent = file.name;
        document.getElementById('selected-schema-size').textContent = formatFileSize(file.size);
        document.getElementById('selected-schema-info').classList.remove('hidden');
        document.getElementById('schema-error').classList.add('hidden');
    }

    function parseSchemaFile(file) {
        const reader = new FileReader();
        reader.onload = function(e) {
            try {
                const schema = JSON.parse(e.target.result);
                validateSchema(schema);
                fakeitData.schema = schema;
                generateFieldConfigs(schema);
                
                // Enable buttons
                document.getElementById('clear-all-btn').disabled = false;
                document.getElementById('generate-data-btn').disabled = false;
                
            } catch (error) {
                showSchemaError('JSON 파싱 오류: ' + error.message);
            }
        };
        
        reader.onerror = function() {
            showSchemaError('파일을 읽을 수 없습니다.');
        };
        
        reader.readAsText(file);
    }

    function validateSchema(schema) {
        if (typeof schema !== 'object' || schema === null) {
            throw new Error('유효하지 않은 JSON 구조입니다.');
        }

        if (schema.type !== 'object') {
            throw new Error('스키마의 타입이 object가 아닙니다.');
        }

        if (!schema.properties || typeof schema.properties !== 'object') {
            throw new Error('properties 필드가 없거나 유효하지 않습니다.');
        }

        if (Object.keys(schema.properties).length === 0) {
            throw new Error('최소한 하나의 속성이 필요합니다.');
        }
    }

    function getSmartTypeRecommendation(fieldName, fieldType) {
        // Normalize field name for matching
        const normalizedName = fieldName.toLowerCase().replace(/[-_\s]/g, '_');
        
        // First try exact match
        if (smartTypeMapping[normalizedName]) {
            return smartTypeMapping[normalizedName];
        }
        
        // Try partial matching for compound names
        for (const [pattern, typeValue] of Object.entries(smartTypeMapping)) {
            if (normalizedName.includes(pattern) || pattern.includes(normalizedName)) {
                return typeValue;
            }
        }
        
        // Default based on field type
        const typeDefaults = {
            'string': 'name',
            'integer': 'number',
            'number': 'float',
            'boolean': 'boolean'
        };
        
        return typeDefaults[fieldType] || 'name';
    }

    function generateFieldConfigs(schema) {
        const container = document.getElementById('field-configs');
        const properties = schema.properties || {};
        
        container.innerHTML = '';
        fakeitData.fieldConfigs = {};

        // Show field config container
        document.getElementById('no-schema-message').classList.add('hidden');
        document.getElementById('field-config-container').classList.remove('hidden');

        Object.entries(properties).forEach(([fieldName, fieldSchema]) => {
            const fieldType = fieldSchema.type || 'string';
            const availableTypes = fakeDataTypes[fieldType] || fakeDataTypes.string;
            
            // Get smart type recommendation
            const recommendedType = getSmartTypeRecommendation(fieldName, fieldType);
            
            // Find the recommended type in available types, or use first as default
            let defaultTypeIndex = 0;
            for (let i = 0; i < availableTypes.length; i++) {
                if (availableTypes[i].value === recommendedType) {
                    defaultTypeIndex = i;
                    break;
                }
            }
            
            const fieldConfig = createFieldConfigElement(fieldName, fieldSchema, availableTypes, defaultTypeIndex);
            container.appendChild(fieldConfig);
            
            // Set default config with smart recommendation
            fakeitData.fieldConfigs[fieldName] = {
                type: availableTypes[defaultTypeIndex].value,
                params: {}
            };
        });
    }

    function createFieldConfigElement(fieldName, fieldSchema, availableTypes, defaultTypeIndex = 0) {
        const div = document.createElement('div');
        div.className = 'border border-gray-200 rounded-lg p-4';
        
        const fieldType = fieldSchema.type || 'string';
        const isRequired = fakeitData.schema.required?.includes(fieldName) || false;
        const isSmartRecommended = defaultTypeIndex > 0; // If not first option, it was smart-recommended
        
        div.innerHTML = `
            <div class="flex items-center justify-between mb-3">
                <div>
                    <h4 class="text-sm font-medium text-gray-900">
                        ${escapeHtml(fieldName)}
                        ${isRequired ? '<span class="text-red-500 text-xs ml-1">*</span>' : ''}
                        ${isSmartRecommended ? '<span class="text-green-600 text-xs ml-2" title="자동 추천됨">✨</span>' : ''}
                    </h4>
                    <p class="text-xs text-gray-500">타입: ${fieldType}</p>
                    ${fieldSchema.description ? `<p class="text-xs text-gray-600 mt-1">${escapeHtml(fieldSchema.description)}</p>` : ''}
                </div>
                <div class="flex items-center space-x-2">
                    <label class="inline-flex items-center">
                        <input type="checkbox" class="field-enabled" data-field="${fieldName}" checked 
                               class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded">
                        <span class="ml-2 text-xs text-gray-700">활성화</span>
                    </label>
                </div>
            </div>
            
            <div class="field-config-content">
                <div class="mb-3">
                    <label class="block text-xs font-medium text-gray-700 mb-1">
                        데이터 타입
                        ${isSmartRecommended ? '<span class="text-green-600 text-xs ml-1">(자동 선택됨)</span>' : ''}
                    </label>
                    <select class="fake-type-select block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 text-sm" 
                            data-field="${fieldName}">
                        ${availableTypes.map((type, index) => `
                            <option value="${type.value}" title="${type.example}" ${index === defaultTypeIndex ? 'selected' : ''}>
                                ${type.label} (예: ${type.example})
                            </option>
                        `).join('')}
                    </select>
                </div>
                
                <div class="field-params" data-field="${fieldName}">
                    ${generateParamInputs(fieldName, availableTypes[defaultTypeIndex], fieldSchema)}
                </div>
            </div>
        `;

        // Add event listeners
        const enableCheckbox = div.querySelector('.field-enabled');
        const typeSelect = div.querySelector('.fake-type-select');
        const configContent = div.querySelector('.field-config-content');

        enableCheckbox.addEventListener('change', (e) => {
            configContent.style.opacity = e.target.checked ? '1' : '0.5';
            configContent.style.pointerEvents = e.target.checked ? 'auto' : 'none';
        });

        typeSelect.addEventListener('change', (e) => {
            const selectedType = availableTypes.find(t => t.value === e.target.value);
            const paramsContainer = div.querySelector('.field-params');
            paramsContainer.innerHTML = generateParamInputs(fieldName, selectedType, fieldSchema);
            
            // Update config
            fakeitData.fieldConfigs[fieldName].type = e.target.value;
            fakeitData.fieldConfigs[fieldName].params = {};
        });

        return div;
    }

    function generateParamInputs(fieldName, typeConfig, fieldSchema) {
        if (typeConfig.value === 'custom') {
            if (fieldSchema.type === 'string') {
                return `
                    <div>
                        <label class="block text-xs font-medium text-gray-700 mb-1">커스텀 값들 (쉼표로 구분)</label>
                        <textarea class="param-input block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 text-sm" 
                                  data-field="${fieldName}" data-param="values" rows="2" 
                                  placeholder="예: Apple, Banana, Cherry"></textarea>
                    </div>
                `;
            } else if (fieldSchema.type === 'integer') {
                return `
                    <div class="grid grid-cols-2 gap-2">
                        <div>
                            <label class="block text-xs font-medium text-gray-700 mb-1">최소값</label>
                            <input type="number" class="param-input block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 text-sm" 
                                   data-field="${fieldName}" data-param="min" value="1" />
                        </div>
                        <div>
                            <label class="block text-xs font-medium text-gray-700 mb-1">최대값</label>
                            <input type="number" class="param-input block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 text-sm" 
                                   data-field="${fieldName}" data-param="max" value="100" />
                        </div>
                    </div>
                `;
            } else if (fieldSchema.type === 'number') {
                return `
                    <div class="grid grid-cols-2 gap-2">
                        <div>
                            <label class="block text-xs font-medium text-gray-700 mb-1">최소값</label>
                            <input type="number" step="0.01" class="param-input block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 text-sm" 
                                   data-field="${fieldName}" data-param="min" value="0.0" />
                        </div>
                        <div>
                            <label class="block text-xs font-medium text-gray-700 mb-1">최대값</label>
                            <input type="number" step="0.01" class="param-input block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 text-sm" 
                                   data-field="${fieldName}" data-param="max" value="100.0" />
                        </div>
                    </div>
                `;
            }
        } else if (typeConfig.value === 'weighted' && fieldSchema.type === 'boolean') {
            return `
                <div>
                    <label class="block text-xs font-medium text-gray-700 mb-1">true 확률 (%)</label>
                    <input type="range" min="0" max="100" value="50" 
                           class="param-input w-full" data-field="${fieldName}" data-param="trueProbability" />
                    <div class="text-xs text-gray-500 mt-1">50%</div>
                </div>
            `;
        } else if (typeConfig.value === 'lorem' && fieldSchema.type === 'string') {
            return `
                <div>
                    <label class="block text-xs font-medium text-gray-700 mb-1">단어 수</label>
                    <input type="number" min="1" max="100" value="5" 
                           class="param-input block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 text-sm" 
                           data-field="${fieldName}" data-param="wordCount" />
                </div>
            `;
        }
        
        return '<p class="text-xs text-gray-500">이 타입은 추가 설정이 필요하지 않습니다.</p>';
    }

    async function generateData() {
        const count = parseInt(document.getElementById('data-count-input').value) || 10;
        
        if (count < 1 || count > 10000) {
            alert('데이터 개수는 1개에서 10,000개 사이여야 합니다.');
            return;
        }

        showLoadingModal();

        try {
            // Collect field configurations
            const fieldConfigs = {};
            document.querySelectorAll('.field-enabled:checked').forEach(checkbox => {
                const fieldName = checkbox.dataset.field;
                const typeSelect = document.querySelector(`.fake-type-select[data-field="${fieldName}"]`);
                const fakeType = typeSelect.value;
                
                const params = {};
                document.querySelectorAll(`.param-input[data-field="${fieldName}"]`).forEach(input => {
                    const paramName = input.dataset.param;
                    params[paramName] = input.value;
                });

                fieldConfigs[fieldName] = {
                    type: fakeType,
                    params: params
                };
            });

            // Generate data via API
            const response = await fetch('/api/fakeit/generate', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    schema: fakeitData.schema,
                    fieldConfigs: fieldConfigs,
                    count: count
                })
            });

            if (!response.ok) {
                const error = await response.text();
                throw new Error(error || 'Failed to generate data');
            }

            const generatedData = await response.json();
            fakeitData.generatedData = generatedData.data;
            
            displayGeneratedData(fakeitData.generatedData);
            
            // Enable export buttons
            document.getElementById('export-json-btn').disabled = false;
            document.getElementById('export-csv-btn').disabled = false;

        } catch (error) {
            console.error('Data generation error:', error);
            alert('데이터 생성 실패: ' + error.message);
        } finally {
            hideLoadingModal();
        }
    }

    function displayGeneratedData(data) {
        const container = document.getElementById('data-preview');
        
        if (!data || data.length === 0) {
            container.innerHTML = `
                <div class="text-center text-gray-500 py-12">
                    <p class="text-sm">생성된 데이터가 없습니다.</p>
                </div>
            `;
            return;
        }

        // Use schema property order instead of Object.keys
        let headers;
        if (fakeitData.schema && fakeitData.schema.properties) {
            // Get headers in the order they appear in the schema
            headers = Object.keys(fakeitData.schema.properties);
            // Add any additional fields that might exist in data but not in schema
            const dataKeys = Object.keys(data[0]);
            dataKeys.forEach(key => {
                if (!headers.includes(key)) {
                    headers.push(key);
                }
            });
        } else {
            // Fallback to data keys if schema is not available
            headers = Object.keys(data[0]);
        }
        
        const previewData = data.slice(0, 20); // Show first 20 rows
        
        const tableHTML = `
            <div class="overflow-x-auto">
                <table class="min-w-full divide-y divide-gray-200">
                    <thead class="bg-gray-50">
                        <tr>
                            ${headers.map(header => `
                                <th class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    ${escapeHtml(header)}
                                </th>
                            `).join('')}
                        </tr>
                    </thead>
                    <tbody class="bg-white divide-y divide-gray-200">
                        ${previewData.map(row => `
                            <tr>
                                ${headers.map(header => `
                                    <td class="px-3 py-2 text-sm text-gray-900 max-w-xs truncate">
                                        ${escapeHtml(String(row[header] || ''))}
                                    </td>
                                `).join('')}
                            </tr>
                        `).join('')}
                    </tbody>
                </table>
            </div>
            ${data.length > 20 ? `
                <div class="mt-4 text-center text-sm text-gray-500">
                    총 ${data.length}개 중 처음 20개 행을 표시 중입니다.
                </div>
            ` : ''}
        `;

        container.innerHTML = tableHTML;
    }

    function exportData(format) {
        if (!fakeitData.generatedData) {
            alert('먼저 데이터를 생성해주세요.');
            return;
        }

        let content, fileName, mimeType;

        if (format === 'json') {
            content = JSON.stringify(fakeitData.generatedData, null, 2);
            fileName = 'fake_data.json';
            mimeType = 'application/json';
        } else if (format === 'csv') {
            content = convertToCSV(fakeitData.generatedData);
            fileName = 'fake_data.csv';
            mimeType = 'text/csv';
        }

        downloadFile(content, fileName, mimeType);
    }

    function convertToCSV(data) {
        if (!data || data.length === 0) return '';

        // Use schema property order for CSV headers too
        let headers;
        if (fakeitData.schema && fakeitData.schema.properties) {
            headers = Object.keys(fakeitData.schema.properties);
            // Add any additional fields
            const dataKeys = Object.keys(data[0]);
            dataKeys.forEach(key => {
                if (!headers.includes(key)) {
                    headers.push(key);
                }
            });
        } else {
            headers = Object.keys(data[0]);
        }
        
        const csvRows = [];

        // Add headers
        csvRows.push(headers.map(header => `"${header.replace(/"/g, '""')}"`).join(','));

        // Add data rows
        data.forEach(row => {
            const values = headers.map(header => {
                const value = row[header] || '';
                return `"${String(value).replace(/"/g, '""')}"`;
            });
            csvRows.push(values.join(','));
        });

        return csvRows.join('\n');
    }

    function downloadFile(content, fileName, mimeType) {
        const blob = new Blob([content], { type: mimeType });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = fileName;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    }

    function clearSchema() {
        // Reset file input
        document.getElementById('schema-file-input').value = '';
        document.getElementById('selected-schema-info').classList.add('hidden');
        document.getElementById('schema-error').classList.add('hidden');
        document.getElementById('no-schema-message').classList.remove('hidden');
        document.getElementById('field-config-container').classList.add('hidden');
        
        // Reset state
        fakeitData.schema = null;
        fakeitData.fieldConfigs = {};
        
        // Disable buttons
        document.getElementById('clear-all-btn').disabled = true;
        document.getElementById('generate-data-btn').disabled = true;
    }

    function clearAll() {
        clearSchema();
        fakeitData.generatedData = null;
        
        // Clear preview
        document.getElementById('data-preview').innerHTML = `
            <div class="text-center text-gray-500 py-12">
                <svg class="mx-auto h-12 w-12 text-gray-300 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M3 14h18m-9-4v8m-7 0V4a1 1 0 011-1h14a1 1 0 011 1v16a1 1 0 01-1 1H5a1 1 0 01-1-1z"/>
                </svg>
                <p class="text-sm">데이터를 생성하면 여기에 미리보기가 표시됩니다</p>
            </div>
        `;
        
        // Disable export buttons
        document.getElementById('export-json-btn').disabled = true;
        document.getElementById('export-csv-btn').disabled = true;
    }

    function showSchemaError(message) {
        document.getElementById('schema-error-message').textContent = message;
        document.getElementById('schema-error').classList.remove('hidden');
        document.getElementById('selected-schema-info').classList.add('hidden');
    }

    function showLoadingModal() {
        document.getElementById('loading-modal').classList.remove('hidden');
    }

    function hideLoadingModal() {
        document.getElementById('loading-modal').classList.add('hidden');
    }

    function formatFileSize(bytes) {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }

    function escapeHtml(str) {
        const div = document.createElement('div');
        div.textContent = str;
        return div.innerHTML;
    }

})();