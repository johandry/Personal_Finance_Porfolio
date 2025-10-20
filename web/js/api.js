// API Client
class FinanceAPI {
    constructor(baseURL) {
        this.baseURL = baseURL;
    }

    async request(endpoint, options = {}) {
        const url = `${this.baseURL}${endpoint}`;
        const config = {
            headers: {
                'Content-Type': 'application/json',
                ...options.headers,
            },
            ...options,
        };

        try {
            const response = await fetch(url, config);
            
            if (!response.ok) {
                const error = await response.json().catch(() => ({}));
                throw new Error(error.error || `HTTP ${response.status}: ${response.statusText}`);
            }

            // Handle empty responses
            const text = await response.text();
            return text ? JSON.parse(text) : {};
        } catch (error) {
            throw error;
        }
    }

    // Assets
    async getAssets() {
        return this.request('/assets');
    }

    async getAsset(id) {
        return this.request(`/assets/${id}`);
    }

    async createAsset(data) {
        return this.request('/assets', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateAsset(id, data) {
        return this.request(`/assets/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteAsset(id) {
        return this.request(`/assets/${id}`, {
            method: 'DELETE',
        });
    }

    async getAssetHistory(id) {
        return this.request(`/assets/${id}/history`);
    }

    // Debts
    async getDebts() {
        return this.request('/debts');
    }

    async getDebt(id) {
        return this.request(`/debts/${id}`);
    }

    async createDebt(data) {
        return this.request('/debts', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateDebt(id, data) {
        return this.request(`/debts/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteDebt(id) {
        return this.request(`/debts/${id}`, {
            method: 'DELETE',
        });
    }

    // Summary
    async getNetWorth() {
        return this.request('/networth');
    }

    async getSummary() {
        return this.request('/summary');
    }

    // Health check
    async healthCheck() {
        return this.request('/health');
    }
}

// Initialize API client
const api = new FinanceAPI(API_CONFIG.baseURL);
