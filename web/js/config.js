// API Configuration
const API_CONFIG = {
    baseURL: 'http://localhost:8080/api/v1',
    timeout: 10000,
};

// Format currency
function formatCurrency(amount, currency = 'USD') {
    return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: currency,
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
    }).format(amount);
}

// Format date
function formatDate(dateString) {
    const date = new Date(dateString);
    return new Intl.DateFormat('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
    }).format(date);
}

// Format date for input field
function formatDateForInput(dateString) {
    const date = new Date(dateString);
    return date.toISOString().split('T')[0];
}

// Calculate profit/loss
function calculateProfitLoss(buyPrice, currentValue, quantity) {
    return (currentValue - buyPrice) * quantity;
}

// Calculate total value
function calculateTotalValue(currentValue, quantity) {
    return currentValue * quantity;
}

// Format asset type
function formatAssetType(type) {
    return type.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase());
}

// Get badge class for type
function getBadgeClass(type) {
    return `badge ${type.toLowerCase().replace('_', '-')}`;
}

// Show toast notification
function showToast(message, type = 'success') {
    const toast = document.getElementById('toast');
    if (!toast) return;

    toast.textContent = message;
    toast.className = `toast ${type} show`;

    setTimeout(() => {
        toast.className = 'toast';
    }, 3000);
}

// Handle API errors
function handleError(error) {
    console.error('API Error:', error);
    const message = error.message || 'An unexpected error occurred';
    showToast(message, 'error');
}
