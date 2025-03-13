const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080/api/v1';

async function fetchAPI(url, options = {}) {
    const headers = {
        'Content-Type': 'application/json',
        ...options.headers,
    };

    const token = localStorage.getItem('token');
    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE_URL}${url}`, {
        ...options,
        headers,
    });

    const data = await response.json();

    if (!response.ok) {
        const error = new Error(data.error || 'Something went wrong');
        error.status = response.status;
        error.data = data;
        throw error;
    }

    return data;
}

const authAPI = {
    register: (userData) => {
        return fetchAPI('/auth/register', {
            method: 'POST',
            body: JSON.stringify(userData),
        });
    },

    login: (credentials) => {
        return fetchAPI('/auth/login', {
            method: 'POST',
            body: JSON.stringify(credentials),
        });
    },

    getProfile: () => {
        return fetchAPI('/users/userinfo');
    },
};

const userAPI = {
    getProfile: () => {
        return fetchAPI('/users/userinfo');
    },
};

export { authAPI, userAPI };