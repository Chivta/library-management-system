const Auth = {
    getToken() {
        if (!AppState.authToken) {
            AppState.authToken = localStorage.getItem('authToken');
        }
        return AppState.authToken;
    },

    setToken(token) {
        AppState.authToken = token;
        localStorage.setItem('authToken', token);
    },

    clearToken() {
        AppState.authToken = null;
        localStorage.removeItem('authToken');
    },

    isAuthenticated() {
        return !!this.getToken();
    },

    async register(username, email, password) {
        return await apiRequest('/auth/register', {
            method: 'POST',
            body: JSON.stringify({ username, email, password }),
            skipAuth: true
        });
    },

    async login(username, password) {
        return await apiRequest('/auth/login', {
            method: 'POST',
            body: JSON.stringify({ username, password }),
            skipAuth: true
        });
    },

    async getProfile() {
        return await apiRequest('/auth/profile');
    },

    logout() {
        this.clearToken();
        AppState.currentUser = null;
        UI.showAuthContainer();
        UI.showNotification('Logged out successfully', 'success');
    }
};

function switchAuthTab(tabName) {
    document.querySelectorAll('.auth-tab').forEach(btn => {
        btn.classList.remove('active');
        if (btn.dataset.authTab === tabName) {
            btn.classList.add('active');
        }
    });

    document.querySelectorAll('.auth-form-container').forEach(container => {
        container.classList.remove('active');
    });
    document.getElementById(`${tabName}-form-container`).classList.add('active');
}

async function handleLoginSubmit(e) {
    e.preventDefault();

    const submitBtn = document.getElementById('login-submit-btn');
    submitBtn.disabled = true;
    submitBtn.textContent = 'Logging in...';

    UI.clearFormErrors('login');

    const username = document.getElementById('login-username').value.trim();
    const password = document.getElementById('login-password').value;

    try {
        const response = await Auth.login(username, password);
        Auth.setToken(response.token);
        AppState.currentUser = {
            id: response.user_id,
            username: response.username,
            email: response.email,
            role: response.role
        };

        document.getElementById('username-display').textContent = response.username;
        UI.showAppContainer();
        UI.showNotification(`Welcome back, ${response.username}!`, 'success');
        await Books.load();
    } catch (error) {
        if (error.validationErrors) {
            UI.displayValidationErrors(error.validationErrors, 'login');
        }
        UI.showNotification(error.message, 'error');
    } finally {
        submitBtn.disabled = false;
        submitBtn.textContent = 'Login';
    }
}

async function handleRegisterSubmit(e) {
    e.preventDefault();

    const submitBtn = document.getElementById('register-submit-btn');
    submitBtn.disabled = true;
    submitBtn.textContent = 'Creating account...';

    UI.clearFormErrors('register');

    const username = document.getElementById('register-username').value.trim();
    const email = document.getElementById('register-email').value.trim();
    const password = document.getElementById('register-password').value;
    const passwordConfirm = document.getElementById('register-password-confirm').value;

    if (password !== passwordConfirm) {
        document.getElementById('register-password-confirm-error').textContent = 'Passwords do not match';
        submitBtn.disabled = false;
        submitBtn.textContent = 'Register';
        return;
    }

    try {
        const response = await Auth.register(username, email, password);
        Auth.setToken(response.token);
        AppState.currentUser = {
            id: response.user_id,
            username: response.username,
            email: response.email,
            role: response.role
        };

        document.getElementById('username-display').textContent = response.username;
        UI.showAppContainer();
        UI.showNotification(`Welcome, ${response.username}! Your account has been created.`, 'success');
        await Books.load();
    } catch (error) {
        if (error.validationErrors) {
            UI.displayValidationErrors(error.validationErrors, 'register');
        }
        UI.showNotification(error.message, 'error');
    } finally {
        submitBtn.disabled = false;
        submitBtn.textContent = 'Register';
    }
}
