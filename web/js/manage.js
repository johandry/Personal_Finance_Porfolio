// Manage page functionality
let currentEditingAssetId = null;
let currentEditingDebtId = null;

// Initialize manage page
function initManage() {
    setupTabs();
    setupModals();
    setupExportImport();
    loadAssets();
    loadDebts();
}

// Setup tab navigation
function setupTabs() {
    const tabButtons = document.querySelectorAll('.tab-btn');
    const tabContents = document.querySelectorAll('.tab-content');

    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            const tabName = button.getAttribute('data-tab');

            // Remove active class from all
            tabButtons.forEach(btn => btn.classList.remove('active'));
            tabContents.forEach(content => content.classList.remove('active'));

            // Add active class to current
            button.classList.add('active');
            document.getElementById(`${tabName}-tab`).classList.add('active');
        });
    });
}

// Setup modals
function setupModals() {
    // Asset modal
    const assetModal = document.getElementById('assetModal');
    const addAssetBtn = document.getElementById('addAssetBtn');
    const closeAssetModal = document.getElementById('closeAssetModal');
    const cancelAssetBtn = document.getElementById('cancelAssetBtn');
    const assetForm = document.getElementById('assetForm');

    addAssetBtn.addEventListener('click', () => {
        openAssetModal();
    });

    closeAssetModal.addEventListener('click', () => {
        closeModal(assetModal);
    });

    cancelAssetBtn.addEventListener('click', () => {
        closeModal(assetModal);
    });

    assetForm.addEventListener('submit', handleAssetSubmit);

    // Listen for changes in asset type and source to hide/show current value
    const assetType = document.getElementById('assetType');
    const assetSource = document.getElementById('assetSource');
    const currentValueField = document.getElementById('assetCurrentValue').closest('.form-group');

    function toggleCurrentValueField() {
        const type = assetType.value;
        const source = assetSource.value;
        
        // Hide current value for stocks with market_api source
        if (type === 'stock' && source === 'market_api') {
            currentValueField.style.display = 'none';
            document.getElementById('assetCurrentValue').removeAttribute('required');
        } else {
            currentValueField.style.display = 'block';
        }
    }

    assetType.addEventListener('change', toggleCurrentValueField);
    assetSource.addEventListener('change', toggleCurrentValueField);

    // Debt modal
    const debtModal = document.getElementById('debtModal');
    const addDebtBtn = document.getElementById('addDebtBtn');
    const closeDebtModal = document.getElementById('closeDebtModal');
    const cancelDebtBtn = document.getElementById('cancelDebtBtn');
    const debtForm = document.getElementById('debtForm');

    addDebtBtn.addEventListener('click', () => {
        openDebtModal();
    });

    closeDebtModal.addEventListener('click', () => {
        closeModal(debtModal);
    });

    cancelDebtBtn.addEventListener('click', () => {
        closeModal(debtModal);
    });

    debtForm.addEventListener('submit', handleDebtSubmit);

    // Close modal when clicking outside
    window.addEventListener('click', (event) => {
        if (event.target === assetModal) {
            closeModal(assetModal);
        }
        if (event.target === debtModal) {
            closeModal(debtModal);
        }
    });
}

// Asset Modal Functions
function openAssetModal(asset = null) {
    const modal = document.getElementById('assetModal');
    const title = document.getElementById('assetModalTitle');
    const form = document.getElementById('assetForm');

    if (asset) {
        // Edit mode
        title.textContent = 'Edit Asset';
        currentEditingAssetId = asset.id;
        
        document.getElementById('assetId').value = asset.id;
        document.getElementById('assetName').value = asset.name;
        document.getElementById('assetType').value = asset.type;
        document.getElementById('assetBuyPrice').value = asset.buy_price;
        document.getElementById('assetCurrentValue').value = asset.current_value;
        document.getElementById('assetQuantity').value = asset.quantity;
        document.getElementById('assetCurrency').value = asset.currency || 'USD';
        document.getElementById('assetPurchaseDate').value = formatDateForInput(asset.purchase_date);
        document.getElementById('assetSource').value = asset.source || 'manual';
    } else {
        // Create mode
        title.textContent = 'Add New Asset';
        currentEditingAssetId = null;
        form.reset();
        
        // Set default date to today
        document.getElementById('assetPurchaseDate').value = new Date().toISOString().split('T')[0];
    }

    // Trigger field visibility update
    const typeField = document.getElementById('assetType');
    const sourceField = document.getElementById('assetSource');
    if (typeField && sourceField) {
        typeField.dispatchEvent(new Event('change'));
    }

    modal.classList.add('active');
}

function openDebtModal(debt = null) {
    const modal = document.getElementById('debtModal');
    const title = document.getElementById('debtModalTitle');
    const form = document.getElementById('debtForm');

    if (debt) {
        // Edit mode
        title.textContent = 'Edit Debt';
        currentEditingDebtId = debt.id;
        
        document.getElementById('debtId').value = debt.id;
        document.getElementById('debtName').value = debt.name;
        document.getElementById('debtType').value = debt.type;
        document.getElementById('debtPrincipal').value = debt.principal;
        document.getElementById('debtCurrentValue').value = debt.current_value;
        document.getElementById('debtInterestRate').value = debt.interest_rate;
        document.getElementById('debtCurrency').value = debt.currency || 'USD';
        document.getElementById('debtStartDate').value = formatDateForInput(debt.start_date);
    } else {
        // Create mode
        title.textContent = 'Add New Debt';
        currentEditingDebtId = null;
        form.reset();
        
        // Set default date to today
        document.getElementById('debtStartDate').value = new Date().toISOString().split('T')[0];
    }

    modal.classList.add('active');
}

function closeModal(modal) {
    modal.classList.remove('active');
}

// Handle asset form submission
async function handleAssetSubmit(event) {
    event.preventDefault();

    const assetType = document.getElementById('assetType').value;
    const assetSource = document.getElementById('assetSource').value;

    const data = {
        name: document.getElementById('assetName').value,
        type: assetType,
        buy_price: parseFloat(document.getElementById('assetBuyPrice').value),
        quantity: parseFloat(document.getElementById('assetQuantity').value),
        currency: document.getElementById('assetCurrency').value,
        purchase_date: document.getElementById('assetPurchaseDate').value,
        source: assetSource,
    };

    // Only include current_value if not a stock with market_api source
    const currentValue = document.getElementById('assetCurrentValue').value;
    if (currentValue && !(assetType === 'stock' && assetSource === 'market_api')) {
        data.current_value = parseFloat(currentValue);
    }

    try {
        if (currentEditingAssetId) {
            // Update existing asset
            await api.updateAsset(currentEditingAssetId, {
                name: data.name,
                current_value: data.current_value,
                quantity: data.quantity,
                source: data.source,
            });
            showToast('Asset updated successfully!', 'success');
        } else {
            // Create new asset
            await api.createAsset(data);
            showToast('Asset created successfully!', 'success');
        }

        closeModal(document.getElementById('assetModal'));
        loadAssets();
    } catch (error) {
        handleError(error);
    }
}

// Handle debt form submission
async function handleDebtSubmit(event) {
    event.preventDefault();

    const data = {
        name: document.getElementById('debtName').value,
        type: document.getElementById('debtType').value,
        principal: parseFloat(document.getElementById('debtPrincipal').value),
        currency: document.getElementById('debtCurrency').value,
        interest_rate: parseFloat(document.getElementById('debtInterestRate').value || 0),
        start_date: document.getElementById('debtStartDate').value,
    };

    const currentValue = document.getElementById('debtCurrentValue').value;
    if (currentValue) {
        data.current_value = parseFloat(currentValue);
    }

    try {
        if (currentEditingDebtId) {
            // Update existing debt
            await api.updateDebt(currentEditingDebtId, {
                name: data.name,
                current_value: data.current_value,
                interest_rate: data.interest_rate,
            });
            showToast('Debt updated successfully!', 'success');
        } else {
            // Create new debt
            await api.createDebt(data);
            showToast('Debt created successfully!', 'success');
        }

        closeModal(document.getElementById('debtModal'));
        loadDebts();
    } catch (error) {
        handleError(error);
    }
}

// Load assets
async function loadAssets() {
    const tbody = document.getElementById('assetsTable');
    tbody.innerHTML = '<tr><td colspan="9" class="loading">Loading...</td></tr>';

    try {
        const assets = await api.getAssets();

        if (!assets || assets.length === 0) {
            tbody.innerHTML = '<tr><td colspan="9" class="empty">No assets found. Click "Add New Asset" to get started!</td></tr>';
            return;
        }

        tbody.innerHTML = assets.map(asset => {
            const totalValue = calculateTotalValue(asset.current_value, asset.quantity);
            const profitLoss = calculateProfitLoss(asset.buy_price, asset.current_value, asset.quantity);

            return `
                <tr>
                    <td><strong>${asset.name}</strong></td>
                    <td><span class="${getBadgeClass(asset.type)}">${formatAssetType(asset.type)}</span></td>
                    <td>${formatCurrency(asset.buy_price, asset.currency)}</td>
                    <td>${formatCurrency(asset.current_value, asset.currency)}</td>
                    <td>${asset.quantity}</td>
                    <td><strong>${formatCurrency(totalValue, asset.currency)}</strong></td>
                    <td class="${profitLoss >= 0 ? 'positive' : 'negative'}">
                        <strong>${formatCurrency(profitLoss, asset.currency)}</strong>
                    </td>
                    <td>${formatDate(asset.purchase_date)}</td>
                    <td>
                        <button class="btn btn-edit" onclick="editAsset('${asset.id}')">Edit</button>
                        <button class="btn btn-danger" onclick="deleteAsset('${asset.id}', '${asset.name}')">Delete</button>
                    </td>
                </tr>
            `;
        }).join('');
    } catch (error) {
        handleError(error);
        tbody.innerHTML = '<tr><td colspan="9" class="empty">Failed to load assets</td></tr>';
    }
}

// Load debts
async function loadDebts() {
    const tbody = document.getElementById('debtsTable');
    tbody.innerHTML = '<tr><td colspan="7" class="loading">Loading...</td></tr>';

    try {
        const debts = await api.getDebts();

        if (!debts || debts.length === 0) {
            tbody.innerHTML = '<tr><td colspan="7" class="empty">No debts found.</td></tr>';
            return;
        }

        tbody.innerHTML = debts.map(debt => `
            <tr>
                <td><strong>${debt.name}</strong></td>
                <td><span class="${getBadgeClass(debt.type)}">${formatAssetType(debt.type)}</span></td>
                <td>${formatCurrency(debt.principal, debt.currency)}</td>
                <td><strong>${formatCurrency(debt.current_value, debt.currency)}</strong></td>
                <td>${debt.interest_rate}%</td>
                <td>${formatDate(debt.start_date)}</td>
                <td>
                    <button class="btn btn-edit" onclick="editDebt('${debt.id}')">Edit</button>
                    <button class="btn btn-danger" onclick="deleteDebt('${debt.id}', '${debt.name}')">Delete</button>
                </td>
            </tr>
        `).join('');
    } catch (error) {
        handleError(error);
        tbody.innerHTML = '<tr><td colspan="7" class="empty">Failed to load debts</td></tr>';
    }
}

// Edit asset
async function editAsset(id) {
    try {
        const asset = await api.getAsset(id);
        openAssetModal(asset);
    } catch (error) {
        handleError(error);
    }
}

// Delete asset
async function deleteAsset(id, name) {
    if (!confirm(`Are you sure you want to delete "${name}"?`)) {
        return;
    }

    try {
        await api.deleteAsset(id);
        showToast('Asset deleted successfully!', 'success');
        loadAssets();
    } catch (error) {
        handleError(error);
    }
}

// Edit debt
async function editDebt(id) {
    try {
        const debt = await api.getDebt(id);
        openDebtModal(debt);
    } catch (error) {
        handleError(error);
    }
}

// Delete debt
async function deleteDebt(id, name) {
    if (!confirm(`Are you sure you want to delete "${name}"?`)) {
        return;
    }

    try {
        await api.deleteDebt(id);
        showToast('Debt deleted successfully!', 'success');
        loadDebts();
    } catch (error) {
        handleError(error);
    }
}

// Setup export/import functionality
function setupExportImport() {
    // Assets Export Buttons
    document.getElementById('exportAssetsJSON').addEventListener('click', () => {
        window.location.href = `${API_CONFIG.baseURL}/export/assets/json`;
    });

    document.getElementById('exportAssetsCSV').addEventListener('click', () => {
        window.location.href = `${API_CONFIG.baseURL}/export/assets/csv`;
    });

    // Debts Export Buttons
    document.getElementById('exportDebtsJSON').addEventListener('click', () => {
        window.location.href = `${API_CONFIG.baseURL}/export/debts/json`;
    });

    document.getElementById('exportDebtsCSV').addEventListener('click', () => {
        window.location.href = `${API_CONFIG.baseURL}/export/debts/csv`;
    });

    // Assets Import Button
    document.getElementById('importAssetsBtn').addEventListener('click', () => {
        document.getElementById('importAssetsFile').click();
    });

    document.getElementById('importAssetsFile').addEventListener('change', async (e) => {
        const file = e.target.files[0];
        if (!file) return;

        const fileExtension = file.name.split('.').pop().toLowerCase();
        if (fileExtension !== 'json' && fileExtension !== 'csv') {
            showToast('Please select a JSON or CSV file', 'error');
            return;
        }

        await importAssets(file, fileExtension);
        e.target.value = ''; // Reset file input
    });

    // Debts Import Button
    document.getElementById('importDebtsBtn').addEventListener('click', () => {
        document.getElementById('importDebtsFile').click();
    });

    document.getElementById('importDebtsFile').addEventListener('change', async (e) => {
        const file = e.target.files[0];
        if (!file) return;

        const fileExtension = file.name.split('.').pop().toLowerCase();
        if (fileExtension !== 'json' && fileExtension !== 'csv') {
            showToast('Please select a JSON or CSV file', 'error');
            return;
        }

        await importDebts(file, fileExtension);
        e.target.value = ''; // Reset file input
    });
}

// Import assets from file
async function importAssets(file, format) {
    try {
        const formData = new FormData();
        formData.append('file', file);

        const response = await fetch(`${API_CONFIG.baseURL}/import/assets/${format}`, {
            method: 'POST',
            body: file,
            headers: {
                'Content-Type': format === 'json' ? 'application/json' : 'text/csv',
            }
        });

        if (!response.ok) {
            throw new Error('Import failed');
        }

        const result = await response.json();
        
        let message = `Imported ${result.imported} assets`;
        if (result.skipped > 0) {
            message += `, skipped ${result.skipped} duplicates`;
        }
        if (result.errors && result.errors.length > 0) {
            message += `, ${result.errors.length} errors`;
            console.error('Import errors:', result.errors);
        }

        showToast(message, result.errors && result.errors.length > 0 ? 'warning' : 'success');
        loadAssets(); // Reload assets table
    } catch (error) {
        handleError(error);
    }
}

// Import debts from file
async function importDebts(file, format) {
    try {
        const formData = new FormData();
        formData.append('file', file);

        const response = await fetch(`${API_CONFIG.baseURL}/import/debts/${format}`, {
            method: 'POST',
            body: file,
            headers: {
                'Content-Type': format === 'json' ? 'application/json' : 'text/csv',
            }
        });

        if (!response.ok) {
            throw new Error('Import failed');
        }

        const result = await response.json();
        
        let message = `Imported ${result.imported} debts`;
        if (result.skipped > 0) {
            message += `, skipped ${result.skipped} duplicates`;
        }
        if (result.errors && result.errors.length > 0) {
            message += `, ${result.errors.length} errors`;
            console.error('Import errors:', result.errors);
        }

        showToast(message, result.errors && result.errors.length > 0 ? 'warning' : 'success');
        loadDebts(); // Reload debts table
    } catch (error) {
        handleError(error);
    }
}

// Initialize on page load
document.addEventListener('DOMContentLoaded', initManage);
