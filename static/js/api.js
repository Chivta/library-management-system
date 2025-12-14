async function apiRequest(endpoint, options = {}) {
    try {
        const headers = {
            'Content-Type': 'application/json',
            ...options.headers
        };

        if (!options.skipAuth && Auth.getToken()) {
            headers['Authorization'] = `Bearer ${Auth.getToken()}`;
        }

        const response = await fetch(`${API_BASE}${endpoint}`, {
            headers,
            ...options
        });

        if (response.status === 204) {
            return { success: true };
        }

        const data = await response.json();

        if (!response.ok) {
            if (response.status === 401) {
                Auth.clearToken();
                UI.showAuthContainer();
                throw new Error('Session expired. Please login again.');
            } else if (response.status === 403) {
                throw new Error(data.error || 'Access denied');
            } else if (response.status === 404) {
                throw new Error(data.error || 'Resource not found');
            } else if (response.status === 400) {
                if (data.errors) {
                    throw { validationErrors: data.errors, message: 'Validation failed' };
                }
                throw new Error(data.error || 'Invalid request');
            } else if (response.status === 409) {
                throw new Error(data.error || 'Conflict');
            } else if (response.status === 500) {
                throw new Error(data.error || 'Server error occurred');
            } else {
                throw new Error(data.error || 'An error occurred');
            }
        }

        return data;
    } catch (error) {
        if (error.validationErrors) {
            throw error;
        }
        if (error.message === 'Failed to fetch') {
            throw new Error('Unable to connect to the server. Please ensure the server is running.');
        }
        throw error;
    }
}

const BooksAPI = {
    async getAll() {
        return await apiRequest('/books/');
    },

    async getById(id) {
        return await apiRequest(`/books/${id}`);
    },

    async create(bookData) {
        return await apiRequest('/books/', {
            method: 'POST',
            body: JSON.stringify(bookData)
        });
    },

    async update(id, bookData) {
        return await apiRequest(`/books/${id}`, {
            method: 'PUT',
            body: JSON.stringify(bookData)
        });
    },

    async delete(id) {
        return await apiRequest(`/books/${id}`, {
            method: 'DELETE'
        });
    },

    async deleteAll() {
        return await apiRequest('/books/', {
            method: 'DELETE'
        });
    }
};

const ReadersAPI = {
    async getAll() {
        return await apiRequest('/readers/');
    },

    async getById(id) {
        return await apiRequest(`/readers/${id}`);
    },

    async create(readerData) {
        return await apiRequest('/readers/', {
            method: 'POST',
            body: JSON.stringify(readerData)
        });
    },

    async update(id, readerData) {
        return await apiRequest(`/readers/${id}`, {
            method: 'PUT',
            body: JSON.stringify(readerData)
        });
    },

    async delete(id) {
        return await apiRequest(`/readers/${id}`, {
            method: 'DELETE'
        });
    },

    async deleteAll() {
        return await apiRequest('/readers/', {
            method: 'DELETE'
        });
    }
};
