// API Configuration
const API_CONFIG = {
    baseURL: 'http://localhost:8080/api/v1',
    timeout: 10000,
};

// Format currency
function formatCurrency(amount, currency = 'USD') {
    // Validate currency code
    if (!currency || currency.trim() === '') {
        currency = 'USD';
    }
    
    // Ensure currency is a string and uppercase
    currency = String(currency).toUpperCase().trim();
    
    try {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: currency,
            minimumFractionDigits: 2,
            maximumFractionDigits: 2,
        }).format(amount);
    } catch (e) {
        // Fallback if currency code is invalid
        console.warn(`Invalid currency code: ${currency}, using USD as fallback`);
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD',
            minimumFractionDigits: 2,
            maximumFractionDigits: 2,
        }).format(amount);
    }
}

// Format date
function formatDate(dateString) {
    if (!dateString) return 'N/A';
    
    const date = new Date(dateString);
    
    // Check if date is valid
    if (isNaN(date.getTime())) return 'Invalid Date';
    
    // Try using Intl.DateTimeFormat if available
    if (typeof Intl !== 'undefined' && Intl.DateTimeFormat) {
        try {
            return new Intl.DateTimeFormat('en-US', {
                year: 'numeric',
                month: 'short',
                day: 'numeric',
            }).format(date);
        } catch (e) {
            // Fallback if Intl fails
        }
    }
    
    // Fallback: Manual formatting
    const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 
                    'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    const month = months[date.getMonth()];
    const day = date.getDate();
    const year = date.getFullYear();
    
    return `${month} ${day}, ${year}`;
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
