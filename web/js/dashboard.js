// Dashboard functionality
let assetPieChart = null;
let netWorthChart = null;

// Initialize dashboard
async function initDashboard() {
    try {
        await loadSummary();
        await loadRecentAssets();
        await loadRecentDebts();
        await loadCharts();
    } catch (error) {
        handleError(error);
    }
}

// Load summary data
async function loadSummary() {
    try {
        const data = await api.getSummary();
        
        document.getElementById('totalAssets').textContent = formatCurrency(data.total_assets || 0);
        document.getElementById('totalDebts').textContent = formatCurrency(data.total_debts || 0);
        document.getElementById('netWorth').textContent = formatCurrency(data.net_worth || 0);
        
        const profitLoss = data.total_profit_loss || 0;
        const profitLossElement = document.getElementById('profitLoss');
        profitLossElement.textContent = formatCurrency(profitLoss);
        profitLossElement.className = `amount ${profitLoss >= 0 ? 'positive' : 'negative'}`;
    } catch (error) {
        console.error('Failed to load summary:', error);
        // Set default values
        document.getElementById('totalAssets').textContent = formatCurrency(0);
        document.getElementById('totalDebts').textContent = formatCurrency(0);
        document.getElementById('netWorth').textContent = formatCurrency(0);
        document.getElementById('profitLoss').textContent = formatCurrency(0);
    }
}

// Load recent assets
async function loadRecentAssets() {
    const tbody = document.getElementById('recentAssetsTable');
    
    try {
        const assets = await api.getAssets();
        
        if (!assets || assets.length === 0) {
            tbody.innerHTML = '<tr><td colspan="6" class="empty">No assets found. Add your first asset!</td></tr>';
            return;
        }

        // Show only first 5 assets
        const recentAssets = assets.slice(0, 5);
        
        tbody.innerHTML = recentAssets.map(asset => {
            const totalValue = calculateTotalValue(asset.current_value, asset.quantity);
            const profitLoss = calculateProfitLoss(asset.buy_price, asset.current_value, asset.quantity);
            
            return `
                <tr>
                    <td><strong>${asset.name}</strong></td>
                    <td><span class="${getBadgeClass(asset.type)}">${formatAssetType(asset.type)}</span></td>
                    <td>${asset.quantity}</td>
                    <td>${formatCurrency(asset.current_value, asset.currency)}</td>
                    <td><strong>${formatCurrency(totalValue, asset.currency)}</strong></td>
                    <td class="${profitLoss >= 0 ? 'positive' : 'negative'}">
                        <strong>${formatCurrency(profitLoss, asset.currency)}</strong>
                    </td>
                </tr>
            `;
        }).join('');
    } catch (error) {
        console.error('Failed to load assets:', error);
        tbody.innerHTML = '<tr><td colspan="6" class="empty">Failed to load assets</td></tr>';
    }
}

// Load recent debts
async function loadRecentDebts() {
    const tbody = document.getElementById('recentDebtsTable');
    
    try {
        const debts = await api.getDebts();
        
        if (!debts || debts.length === 0) {
            tbody.innerHTML = '<tr><td colspan="5" class="empty">No debts found.</td></tr>';
            return;
        }

        // Show only first 5 debts
        const recentDebts = debts.slice(0, 5);
        
        tbody.innerHTML = recentDebts.map(debt => `
            <tr>
                <td><strong>${debt.name}</strong></td>
                <td><span class="${getBadgeClass(debt.type)}">${formatAssetType(debt.type)}</span></td>
                <td>${formatCurrency(debt.principal, debt.currency)}</td>
                <td><strong>${formatCurrency(debt.current_value, debt.currency)}</strong></td>
                <td>${debt.interest_rate}%</td>
            </tr>
        `).join('');
    } catch (error) {
        console.error('Failed to load debts:', error);
        tbody.innerHTML = '<tr><td colspan="5" class="empty">Failed to load debts</td></tr>';
    }
}

// Load charts
async function loadCharts() {
    try {
        const assets = await api.getAssets();
        
        if (assets && assets.length > 0) {
            createAssetPieChart(assets);
            createNetWorthChart(assets);
        }
    } catch (error) {
        console.error('Failed to load charts:', error);
    }
}

// Create asset distribution pie chart
function createAssetPieChart(assets) {
    const ctx = document.getElementById('assetPieChart');
    if (!ctx) return;

    // Aggregate by type
    const assetsByType = assets.reduce((acc, asset) => {
        const totalValue = calculateTotalValue(asset.current_value, asset.quantity);
        acc[asset.type] = (acc[asset.type] || 0) + totalValue;
        return acc;
    }, {});

    const labels = Object.keys(assetsByType).map(formatAssetType);
    const data = Object.values(assetsByType);
    const colors = [
        '#4f46e5', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6',
        '#ec4899', '#14b8a6', '#f97316', '#06b6d4'
    ];

    if (assetPieChart) {
        assetPieChart.destroy();
    }

    assetPieChart = new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: labels,
            datasets: [{
                data: data,
                backgroundColor: colors.slice(0, labels.length),
                borderWidth: 2,
                borderColor: '#ffffff'
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            plugins: {
                legend: {
                    position: 'bottom',
                },
                tooltip: {
                    callbacks: {
                        label: function(context) {
                            const label = context.label || '';
                            const value = formatCurrency(context.parsed);
                            const total = context.dataset.data.reduce((a, b) => a + b, 0);
                            const percentage = ((context.parsed / total) * 100).toFixed(1);
                            return `${label}: ${value} (${percentage}%)`;
                        }
                    }
                }
            }
        }
    });
}

// Create net worth bar chart
function createNetWorthChart(assets) {
    const ctx = document.getElementById('netWorthChart');
    if (!ctx) return;

    const totalAssets = assets.reduce((sum, asset) => {
        return sum + calculateTotalValue(asset.current_value, asset.quantity);
    }, 0);

    const totalInvested = assets.reduce((sum, asset) => {
        return sum + (asset.buy_price * asset.quantity);
    }, 0);

    const profitLoss = totalAssets - totalInvested;

    if (netWorthChart) {
        netWorthChart.destroy();
    }

    netWorthChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: ['Invested', 'Current Value', 'Profit/Loss'],
            datasets: [{
                label: 'Amount (USD)',
                data: [totalInvested, totalAssets, profitLoss],
                backgroundColor: [
                    '#64748b',
                    '#4f46e5',
                    profitLoss >= 0 ? '#10b981' : '#ef4444'
                ],
                borderRadius: 8,
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            plugins: {
                legend: {
                    display: false
                },
                tooltip: {
                    callbacks: {
                        label: function(context) {
                            return formatCurrency(context.parsed.y);
                        }
                    }
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: {
                        callback: function(value) {
                            return '$' + value.toLocaleString();
                        }
                    }
                }
            }
        }
    });
}

// Initialize on page load
document.addEventListener('DOMContentLoaded', initDashboard);

// Refresh every 30 seconds
setInterval(loadSummary, 30000);
