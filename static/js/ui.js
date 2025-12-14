const UI = {
    notificationTimeout: null,

    showNotification(message, type = 'info') {
        const notification = document.getElementById('notification');
        notification.textContent = message;
        notification.className = `notification ${type}`;
        notification.classList.remove('hidden');

        // Clear any existing timeout
        if (this.notificationTimeout) {
            clearTimeout(this.notificationTimeout);
        }

        // Set new timeout
        this.notificationTimeout = setTimeout(() => {
            notification.classList.add('hidden');
            this.notificationTimeout = null;
        }, 5000);
    },

    showLoading(element) {
        if (element) element.classList.remove('hidden');
    },

    hideLoading(element) {
        if (element) element.classList.add('hidden');
    },

    clearFormErrors(formPrefix) {
        const errorElements = document.querySelectorAll(`[id^="${formPrefix}-"][id$="-error"]`);
        errorElements.forEach(el => el.textContent = '');
    },

    displayValidationErrors(errors, formPrefix) {
        this.clearFormErrors(formPrefix);

        if (Array.isArray(errors)) {
            errors.forEach(error => {
                const fieldName = error.field.toLowerCase();
                const errorElement = document.getElementById(`${formPrefix}-${fieldName}-error`);
                if (errorElement) {
                    errorElement.textContent = error.message;
                }
            });
        }
    },

    showAuthContainer() {
        document.getElementById('auth-container').classList.remove('hidden');
        document.getElementById('app-container').classList.add('hidden');
    },

    showAppContainer() {
        document.getElementById('auth-container').classList.add('hidden');
        document.getElementById('app-container').classList.remove('hidden');
    },

    showConfirmModal(message, onConfirm) {
        const modal = document.getElementById('confirm-modal');
        const messageEl = document.getElementById('confirm-message');

        messageEl.textContent = message;
        modal.classList.remove('hidden');

        const handleConfirm = async () => {
            modal.classList.add('hidden');
            await onConfirm();
            cleanup();
        };

        const handleCancel = () => {
            modal.classList.add('hidden');
            cleanup();
        };

        const cleanup = () => {
            document.getElementById('confirm-yes').removeEventListener('click', handleConfirm);
            document.getElementById('confirm-no').removeEventListener('click', handleCancel);
        };

        document.getElementById('confirm-yes').addEventListener('click', handleConfirm);
        document.getElementById('confirm-no').addEventListener('click', handleCancel);
    },

    async showProfileModal() {
        try {
            const profile = await Auth.getProfile();
            document.getElementById('profile-username').textContent = profile.username;
            document.getElementById('profile-email').textContent = profile.email;
            document.getElementById('profile-role').textContent = profile.role;
            document.getElementById('profile-modal').classList.remove('hidden');
        } catch (error) {
            this.showNotification('Failed to load profile', 'error');
        }
    },

    hideProfileModal() {
        document.getElementById('profile-modal').classList.add('hidden');
    }
};
